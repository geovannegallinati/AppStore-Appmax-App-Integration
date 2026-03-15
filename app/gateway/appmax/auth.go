package appmax

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetToken(ctx context.Context, clientID, clientSecret string) (TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.authBaseURL+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return TokenResponse{}, fmt.Errorf("new token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("send token request: %w", err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, http.StatusOK); err != nil {
		return TokenResponse{}, err
	}

	var out TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return TokenResponse{}, fmt.Errorf("decode token response: %w", err)
	}

	return out, nil
}

func (c *Client) Authorize(ctx context.Context, appToken, appID, externalKey, callbackURL string) (string, error) {
	payload := AuthorizeRequest{
		AppID:       appID,
		ExternalKey: externalKey,
		URLCallback: callbackURL,
	}

	out, err := doAndDecode[AuthorizeResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/app/authorize", payload, appToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return "", fmt.Errorf("authorize: %w", err)
	}

	return out.Data.Token, nil
}

func (c *Client) GenerateMerchantCreds(ctx context.Context, appToken, hash string) (string, string, error) {
	payload := GenerateCredsRequest{Token: hash}

	out, err := doAndDecode[GenerateCredsResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/app/client/generate", payload, appToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return "", "", fmt.Errorf("generate merchant creds: %w", err)
	}

	return out.Data.Client.ClientID, out.Data.Client.ClientSecret, nil
}
