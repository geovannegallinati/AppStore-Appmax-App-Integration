//go:build appmax_live

package appmax

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"

	gatewayappmax "github.com/geovanne-gallinati/AppStoreAppDemo/app/gateway/appmax"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

var _ services.TokenManager = (*testTokenManager)(nil)

type testTokenManager struct {
	appToken      string
	merchantToken string
}

func (m *testTokenManager) AppToken(_ context.Context) (string, error) {
	return m.appToken, nil
}

func (m *testTokenManager) MerchantToken(_ context.Context, _ *models.Installation) (string, error) {
	return m.merchantToken, nil
}

var fixtureMu sync.Mutex

func testCtx(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func env(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("env var %s is not set — add it to .env", key)
	}
	return v
}

func callbackURL() string {
	base := strings.TrimRight(strings.TrimSpace(env("NGROK_URL")), "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "https://" + base
	}
	return base + "/integrations/appmax/callback/install"
}


func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func fetchOAuthToken(authURL, clientID, clientSecret string) (string, int, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequest(http.MethodPost, authURL+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", 0, fmt.Errorf("oauth2/token unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", 0, err
	}
	return out.AccessToken, out.ExpiresIn, nil
}

func runOAuthFlow(authURL, apiURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	appIDUUID := strings.TrimSpace(env("APP_ID_UUID"))
	attemptExternalKey := buildInstallAttemptExternalKey()
	cbURL := callbackURL()
	log.Printf("Using callback URL: %s", cbURL)
	log.Printf("Installation attempt external_key: %s", attemptExternalKey)

	tmpTokenMgr := &testTokenManager{appToken: creds.AppToken}
	gateway, gatewayErr := gatewayappmax.NewClient(authURL, apiURL)
	if gatewayErr != nil {
		return fmt.Errorf("new appmax client: %w", gatewayErr)
	}
	tmpSvc, svcErr := services.NewAppmaxServiceWithGateway(tmpTokenMgr, gateway)
	if svcErr != nil {
		return fmt.Errorf("new appmax service: %w", svcErr)
	}

	candidates := make([]string, 0, 1)
	if appIDUUID != "" {
		candidates = append(candidates, appIDUUID)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("no app id candidate available for authorize (APP_ID_UUID)")
	}

	var (
		hash              string
		authErr           error
		selectedAuthorize string
	)
	for _, candidate := range candidates {
		hash, authErr = tmpSvc.Authorize(ctx, candidate, attemptExternalKey, cbURL)
		if authErr == nil {
			selectedAuthorize = candidate
			break
		}
		log.Printf("Authorize failed with app_id=%s: %v", candidate, authErr)
	}
	if authErr != nil {
		return fmt.Errorf("authorize installation: %w", authErr)
	}
	log.Printf("Authorize app_id selected: %s", selectedAuthorize)
	log.Printf("Authorization hash acquired: %s", hash)

	redirectURL := fmt.Sprintf("%s/appstore/integration/%s", sandboxBCURL, hash)
	generateToken, browserErr := authorizeInstallationInBrowser(redirectURL)
	if browserErr != nil {
		return fmt.Errorf("browser authorization flow: %w", browserErr)
	}
	if strings.TrimSpace(generateToken) == "" {
		generateToken = hash
	}
	log.Printf("Token selected for app/client/generate: %s", generateToken)

	preWait := 10 * time.Second
	if s := strings.TrimSpace(os.Getenv("PRE_GENERATE_WAIT")); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			preWait = d
		}
	}
	log.Printf("Waiting %s before calling /app/client/generate...", preWait)
	select {
	case <-ctx.Done():
		return fmt.Errorf("pre-generate wait canceled: %w", ctx.Err())
	case <-time.After(preWait):
	}

	freshTok, freshExpiresIn, refreshErr := fetchOAuthToken(authURL, env("APPMAX_CLIENT_ID"), env("APPMAX_CLIENT_SECRET"))
	if refreshErr == nil {
		tmpTokenMgr.appToken = freshTok
		creds.AppToken = freshTok
		creds.AppTokenExpiry = time.Now().Add(time.Duration(freshExpiresIn-60) * time.Second)
		if saveErr := saveCreds(creds); saveErr != nil {
			log.Printf("warn: save refreshed app token: %v", saveErr)
		}
		log.Printf("App token refreshed before credential generation (expires in %ds)", freshExpiresIn)
	} else {
		log.Printf("warn: could not refresh app token before generate: %v", refreshErr)
	}

	mClientID, mSecret, genErr := tmpSvc.GenerateMerchantCreds(ctx, generateToken)
	if genErr != nil {
		return fmt.Errorf("generate merchant creds: %w", genErr)
	}
	log.Printf("Merchant credentials obtained: client_id=%s", mClientID)
	log.Printf("AppMax will POST credentials to %s to complete the handshake", cbURL)

	creds.MerchantClientID = mClientID
	creds.MerchantClientSecret = mSecret
	creds.MerchantToken = ""
	creds.MerchantTokenExpiry = time.Time{}
	if saveErr := saveCreds(creds); saveErr != nil {
		log.Printf("warn: save creds: %v", saveErr)
	}
	return nil
}

func authorizeInstallationInBrowser(redirectURL string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	cdpCtx, cdpCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cdpCancel()

	chromedpCtx, chromedpCancel := context.WithTimeout(cdpCtx, 3*time.Minute)
	defer chromedpCancel()

	bcLogin := env("APPMAX_LOGIN")
	bcPassword := env("APPMAX_PASSWORD")
	storeName := strings.TrimSpace(os.Getenv("APPMAX_STORE_NAME"))
	if storeName == "" {
		storeName = fmt.Sprintf("E2E Test Store %d", time.Now().Unix())
	}

	var finalURL string
	var screenshotBuf []byte
	browserErr := chromedp.Run(chromedpCtx,
		chromedp.Navigate(redirectURL),
		chromedp.Sleep(5*time.Second),
		chromedp.Screenshot(`body`, &screenshotBuf, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`input[type="email"], input[name="email"], input[type="text"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[type="email"], input[name="email"], input[type="text"]`, bcLogin, chromedp.ByQuery),
		chromedp.SendKeys(`input[type="password"]`, bcPassword, chromedp.ByQuery),
		chromedp.Click(`button[type="submit"], input[type="submit"]`, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
		chromedp.Navigate(redirectURL),
		chromedp.Sleep(5*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result string
			script := fmt.Sprintf(`(function(name){
				var input = document.querySelector('input[placeholder="Digite aqui"]');
				if(!input) return 'store name input not found';
				var setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype,'value').set;
				setter.call(input, name);
				input.dispatchEvent(new Event('input',{bubbles:true}));
				input.dispatchEvent(new Event('change',{bubbles:true}));
				return 'store name set: '+name;
			})(%q)`, storeName)
			if err := chromedp.Evaluate(script, &result).Do(ctx); err != nil {
				log.Printf("fill store name: %v", err)
				return nil
			}
			log.Printf("Fill store name: %s", result)
			return nil
		}),
		chromedp.Sleep(time.Second),
		chromedp.Click(`[data-test="am_dropdown-v3"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result string
			err := chromedp.Evaluate(`(function(){
				var li = document.querySelector('li.am_dropdown-v3--list-item');
				if(li){ li.click(); return 'li click: '+li.textContent.trim(); }
				return 'not found';
			})()`, &result).Do(ctx)
			if err != nil {
				log.Printf("select company: %v", err)
				return nil
			}
			log.Printf("Select company: %s", result)
			return nil
		}),
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result string
			err := chromedp.Evaluate(`(function(){
				var label = document.querySelector('[data-test="am-v3-checkbox-v3"] label, label.am-v3-checkbox-v3, .config-app_checkbox label');
				if(label){ label.click(); return 'clicked checkbox label'; }
				return 'checkbox not found';
			})()`, &result).Do(ctx)
			if err != nil {
				log.Printf("check checkbox: %v", err)
				return nil
			}
			log.Printf("Check checkbox: %s", result)
			return nil
		}),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var result string
			err := chromedp.Evaluate(`(function(){
				var btns = document.querySelectorAll('button');
				for(var i=0;i<btns.length;i++){
					var t=btns[i].textContent.trim().toLowerCase();
					if((t.includes('salvar')||t.includes('autorizar')||t.includes('confirmar'))&&!btns[i].disabled){
						btns[i].click();
						return 'clicked: '+btns[i].textContent.trim();
					}
				}
				for(var i=0;i<btns.length;i++){
					if(btns[i].textContent.trim().toLowerCase().includes('salvar')){
						btns[i].disabled=false;
						btns[i].click();
						return 'force-clicked Salvar';
					}
				}
				return 'no enabled save button found';
			})()`, &result).Do(ctx)
			if err != nil {
				log.Printf("click save: %v", err)
				return nil
			}
			log.Printf("Click save result: %s", result)
			return nil
		}),
		chromedp.Sleep(8*time.Second),
		chromedp.Location(&finalURL),
	)

	if len(screenshotBuf) > 0 {
		if writeErr := os.WriteFile("/tmp/bc_install_step1.png", screenshotBuf, 0o644); writeErr == nil {
			log.Printf("Screenshot saved: /tmp/bc_install_step1.png")
		}
	}

	if browserErr != nil {
		return "", browserErr
	}
	log.Printf("Final browser URL after authorization: %s", finalURL)

	parsedURL, parseErr := url.Parse(finalURL)
	if parseErr != nil {
		return "", nil
	}
	return strings.TrimSpace(parsedURL.Query().Get("token")), nil
}


func buildInstallAttemptExternalKey() string {
	suffixBytes := make([]byte, 8)
	if _, err := rand.Read(suffixBytes); err != nil {
		return fmt.Sprintf("e2e-install-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("e2e-install-%d-%s", time.Now().Unix(), hex.EncodeToString(suffixBytes))
}

func ensureCustomerID(t *testing.T) int {
	t.Helper()

	fixtureMu.Lock()
	cachedID := customerID
	fixtureMu.Unlock()
	if cachedID > 0 {
		return cachedID
	}

	id, err := appmaxSvc.CreateOrUpdateCustomer(testCtx(t), testInst, services.CustomerInput{
		FirstName:      "E2E",
		LastName:       "Teste",
		Email:          fmt.Sprintf("e2e+%d@teste.com", time.Now().UnixMilli()),
		Phone:          "51983655100",
		DocumentNumber: docCPF,
		Address: &services.Address{
			Postcode: "91520270",
			Street:   "Rua Francisco Carneiro da Rocha",
			Number:   "582",
			City:     "Porto Alegre",
			State:    "RS",
			District: "Centro",
		},
		IP: "177.92.0.1",
		Products: []services.Product{
			{SKU: "PROD-001", Name: "Produto E2E", Quantity: 1, UnitValue: 5000, Type: "digital"},
		},
	})
	require.NoError(t, err)
	require.Greater(t, id, 0)

	fixtureMu.Lock()
	if customerID == 0 {
		customerID = id
	}
	cachedID = customerID
	fixtureMu.Unlock()

	return cachedID
}

func createOrderForCustomer(t *testing.T, customerID int, sku, name string, unitValue int, productType string) int {
	t.Helper()

	if productType == "" {
		productType = "digital"
	}

	orderID, err := appmaxSvc.CreateOrder(testCtx(t), testInst, services.OrderInput{
		CustomerID:    customerID,
		DiscountValue: 0,
		ShippingValue: 0,
		Products: []services.Product{
			{SKU: sku, Name: name, Quantity: 1, UnitValue: unitValue, Type: productType},
		},
	})
	require.NoError(t, err)
	require.Greater(t, orderID, 0)
	return orderID
}

func createApprovedCreditCardPayment(t *testing.T) (int, int, string) {
	t.Helper()

	cid := ensureCustomerID(t)
	oid := createOrderForCustomer(t, cid, "PROD-APPROVED", "Produto Aprovado E2E", 5000, "digital")

	result, err := appmaxSvc.CreditCard(testCtx(t), testInst, services.CreditCardInput{
		OrderID:              oid,
		CustomerID:           cid,
		Number:               cardSuccess,
		CVV:                  "123",
		ExpirationMonth:      "12",
		ExpirationYear:       "28",
		HolderName:           "TESTE E2E",
		HolderDocumentNumber: docCPF,
		Installments:         1,
		SoftDescriptor:       "E2ETEST",
	})
	require.NoError(t, err)
	require.Greater(t, result.PaymentID, 0)

	if result.UpsellHash != "" {
		fixtureMu.Lock()
		if upsellHash == "" {
			upsellHash = result.UpsellHash
		}
		fixtureMu.Unlock()
	}

	return oid, cid, result.UpsellHash
}

func ensureUpsellHash(t *testing.T) string {
	t.Helper()

	fixtureMu.Lock()
	cached := upsellHash
	fixtureMu.Unlock()
	if cached != "" {
		return cached
	}

	_, _, generatedHash := createApprovedCreditCardPayment(t)
	if generatedHash == "" {
		t.Skip("no upsell_hash returned in approved payment flow")
	}

	fixtureMu.Lock()
	if upsellHash == "" {
		upsellHash = generatedHash
	}
	cached = upsellHash
	fixtureMu.Unlock()

	return cached
}
