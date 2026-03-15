package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	gatewaycontracts "github.com/geovanne-gallinati/AppStoreAppDemo/app/gateway/appmax/contracts"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

const tokenExpiryBuffer = 60 * time.Second

type Clock func() time.Time

type TokenManager interface {
	AppToken(ctx context.Context) (string, error)
	MerchantToken(ctx context.Context, inst *models.Installation) (string, error)
}

type tokenManager struct {
	gateway         gatewaycontracts.TokenGateway
	appClientID     string
	appClientSecret string
	clock           Clock
	cache           MerchantTokenCache
	logger          Logger

	mu             sync.Mutex
	cachedAppToken string
	appTokenExpiry time.Time
}

var _ TokenManager = (*tokenManager)(nil)

func NewTokenManagerWithGateway(gateway gatewaycontracts.TokenGateway, appClientID, appClientSecret string) (TokenManager, error) {
	return newTokenManager(gateway, appClientID, appClientSecret, time.Now, nil, nil)
}

func NewTokenManagerWithGatewayDeps(
	gateway gatewaycontracts.TokenGateway,
	appClientID, appClientSecret string,
	cache MerchantTokenCache,
	logger Logger,
) (TokenManager, error) {
	return NewTokenManagerWithGatewayDepsAndClock(gateway, appClientID, appClientSecret, cache, logger, time.Now)
}

func NewTokenManagerWithGatewayDepsAndClock(
	gateway gatewaycontracts.TokenGateway,
	appClientID, appClientSecret string,
	cache MerchantTokenCache,
	logger Logger,
	clock Clock,
) (TokenManager, error) {
	return newTokenManager(gateway, appClientID, appClientSecret, clock, cache, logger)
}

func newTokenManager(
	gateway gatewaycontracts.TokenGateway,
	appClientID, appClientSecret string,
	clock Clock,
	cache MerchantTokenCache,
	logger Logger,
) (*tokenManager, error) {
	if gateway == nil {
		return nil, fmt.Errorf("new token manager: %w", ErrNilDependency)
	}
	if strings.TrimSpace(appClientID) == "" || strings.TrimSpace(appClientSecret) == "" {
		return nil, fmt.Errorf("new token manager: %w", ErrInvalidConfig)
	}
	if clock == nil {
		return nil, fmt.Errorf("new token manager: %w", ErrNilDependency)
	}

	return &tokenManager{
		gateway:         gateway,
		appClientID:     appClientID,
		appClientSecret: appClientSecret,
		clock:           clock,
		cache:           coalesceMerchantTokenCache(cache),
		logger:          coalesceLogger(logger),
	}, nil
}

func (t *tokenManager) AppToken(ctx context.Context) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cachedAppToken != "" && t.clock().Before(t.appTokenExpiry) {
		return t.cachedAppToken, nil
	}

	resp, err := t.gateway.GetToken(ctx, t.appClientID, t.appClientSecret)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenFetch, err)
	}

	expiry := time.Duration(resp.ExpiresIn)*time.Second - tokenExpiryBuffer
	t.cachedAppToken = resp.AccessToken
	t.appTokenExpiry = t.clock().Add(expiry)

	return t.cachedAppToken, nil
}

func (t *tokenManager) MerchantToken(ctx context.Context, inst *models.Installation) (string, error) {
	if inst == nil {
		return "", fmt.Errorf("%w: installation is nil", ErrNilDependency)
	}
	if strings.TrimSpace(inst.MerchantClientID) == "" || strings.TrimSpace(inst.MerchantClientSecret) == "" {
		return "", fmt.Errorf("%w: merchant credentials missing", ErrInvalidConfig)
	}

	redisKey := fmt.Sprintf("merchant_token:%d", inst.ID)

	cached := t.cache.GetString(ctx, redisKey)
	if cached != "" {
		return cached, nil
	}

	resp, err := t.gateway.GetToken(ctx, inst.MerchantClientID, inst.MerchantClientSecret)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTokenFetch, err)
	}

	ttl := time.Duration(resp.ExpiresIn)*time.Second - tokenExpiryBuffer
	if putErr := t.cache.Put(ctx, redisKey, resp.AccessToken, ttl); putErr != nil {
		t.logger.Warningf("token_manager: failed to cache merchant token for installation %d: %v", inst.ID, putErr)
	}

	return resp.AccessToken, nil
}
