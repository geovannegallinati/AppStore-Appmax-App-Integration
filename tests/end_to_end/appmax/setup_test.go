//go:build appmax_live

package appmax

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	gatewayappmax "github.com/geovanne-gallinati/AppStoreAppDemo/app/gateway/appmax"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

const (
	cardSuccess = "4000000000000010"
	cardFail    = "4000000000000028"
	docCPF      = "25226493029"

	sandboxAuthURL     = "https://auth.sandboxappmax.com.br"
	sandboxAPIURL      = "https://api.sandboxappmax.com.br"
	sandboxBCURL       = "https://breakingcode.sandboxappmax.com.br"
)

var (
	appmaxSvc  services.AppmaxService
	testInst   *models.Installation
	creds      *SandboxCreds
	customerID int
	upsellHash string
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")

	var loadErr error
	creds, loadErr = loadCreds()
	if loadErr != nil {
		log.Fatalf("load creds: %v", loadErr)
	}

	if !creds.IsAppTokenValid() {
		log.Println("Getting app token...")
		tok, expiresIn, tokErr := fetchOAuthToken(sandboxAuthURL, env("APPMAX_CLIENT_ID"), env("APPMAX_CLIENT_SECRET"))
		if tokErr != nil {
			log.Fatalf("GetAppToken: %v", tokErr)
		}
		creds.AppToken = tok
		creds.AppTokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)
		if saveErr := saveCreds(creds); saveErr != nil {
			log.Printf("warn: save creds: %v", saveErr)
		}
		log.Printf("App token acquired (expires in %ds)", expiresIn)
	} else {
		log.Println("Using cached app token")
	}

	if !creds.IsMerchantCredentialsReady() {
		log.Println("Merchant credentials not found — running OAuth flow...")
		if oauthErr := runOAuthFlow(sandboxAuthURL, sandboxAPIURL); oauthErr != nil {
			log.Fatalf("OAuth flow: %v", oauthErr)
		}
	} else {
		log.Println("Using cached merchant credentials")
	}

	if !creds.IsMerchantTokenValid() {
		log.Println("Getting merchant token...")
		tok, expiresIn, tokErr := fetchOAuthToken(sandboxAuthURL, creds.MerchantClientID, creds.MerchantClientSecret)
		if tokErr != nil {
			log.Fatalf("GetMerchantToken: %v", tokErr)
		}
		creds.MerchantToken = tok
		creds.MerchantTokenExpiry = time.Now().Add(time.Duration(expiresIn-60) * time.Second)
		if saveErr := saveCreds(creds); saveErr != nil {
			log.Printf("warn: save creds: %v", saveErr)
		}
		log.Printf("Merchant token acquired (expires in %ds)", expiresIn)
	} else {
		log.Println("Using cached merchant token")
	}

	tokenMgr := &testTokenManager{
		appToken:      creds.AppToken,
		merchantToken: creds.MerchantToken,
	}

	gateway, gatewayErr := gatewayappmax.NewClient(sandboxAuthURL, sandboxAPIURL)
	if gatewayErr != nil {
		log.Fatalf("new appmax client: %v", gatewayErr)
	}

	var svcErr error
	appmaxSvc, svcErr = services.NewAppmaxServiceWithGateway(tokenMgr, gateway)
	if svcErr != nil {
		log.Fatalf("new appmax service: %v", svcErr)
	}
	testInst = &models.Installation{
		ID:                   1,
		MerchantClientID:     creds.MerchantClientID,
		MerchantClientSecret: creds.MerchantClientSecret,
	}

	os.Exit(m.Run())
}
