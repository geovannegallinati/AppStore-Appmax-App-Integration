package services_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	gatewayappmax "github.com/geovanne-gallinati/AppStoreAppDemo/app/gateway/appmax"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTokenGateway struct {
	getTokenFunc func(ctx context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error)
}

func (m *mockTokenGateway) GetToken(ctx context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error) {
	if m.getTokenFunc == nil {
		return gatewayappmax.TokenResponse{}, nil
	}

	return m.getTokenFunc(ctx, clientID, clientSecret)
}

type memoryTokenCache struct {
	mu       sync.Mutex
	values   map[string]string
	putErr   error
	lastTTL  time.Duration
	lastKey  string
	lastData string
}

func newMemoryTokenCache() *memoryTokenCache {
	return &memoryTokenCache{values: map[string]string{}}
}

func (c *memoryTokenCache) GetString(_ context.Context, key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.values[key]
}

func (c *memoryTokenCache) Put(_ context.Context, key, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastKey = key
	c.lastData = value
	c.lastTTL = ttl
	if c.putErr != nil {
		return c.putErr
	}
	c.values[key] = value
	return nil
}

type captureLogger struct {
	warnings []string
	errors   []string
}

func (l *captureLogger) Warningf(format string, args ...any) {
	l.warnings = append(l.warnings, format)
}

func (l *captureLogger) Errorf(format string, args ...any) {
	l.errors = append(l.errors, format)
}

func TestAppToken_FetchesOnFirstCall(t *testing.T) {
	calls := 0
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error) {
			calls++
			assert.Equal(t, "id", clientID)
			assert.Equal(t, "secret", clientSecret)
			return gatewayappmax.TokenResponse{AccessToken: "app-token-abc", ExpiresIn: 3600}, nil
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "id", "secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)
	tok, err := mgr.AppToken(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "app-token-abc", tok)
	assert.Equal(t, 1, calls)
}

func TestAppToken_UsesCache(t *testing.T) {
	calls := 0
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, _, _ string) (gatewayappmax.TokenResponse, error) {
			calls++
			return gatewayappmax.TokenResponse{AccessToken: "app-token-abc", ExpiresIn: 3600}, nil
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "id", "secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)

	_, err = mgr.AppToken(context.Background())
	require.NoError(t, err)
	_, err = mgr.AppToken(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 1, calls, "second call should hit cache, not re-fetch")
}

func TestAppToken_RefetchesWhenExpired(t *testing.T) {
	calls := 0
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, _, _ string) (gatewayappmax.TokenResponse, error) {
			calls++
			return gatewayappmax.TokenResponse{AccessToken: "app-token-abc", ExpiresIn: 61}, nil
		},
	}

	now := time.Unix(1000, 0)
	clock := func() time.Time { return now }

	mgr, err := services.NewTokenManagerWithGatewayDepsAndClock(gateway, "id", "secret", newMemoryTokenCache(), &captureLogger{}, clock)
	require.NoError(t, err)

	_, err = mgr.AppToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, calls)

	now = now.Add(2 * time.Second)

	_, err = mgr.AppToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, calls, "should re-fetch after expiry (61s - 60s buffer = 1s TTL)")
}

func TestAppToken_ConcurrentCallsResultInOneFetch(t *testing.T) {
	calls := 0
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, _, _ string) (gatewayappmax.TokenResponse, error) {
			calls++
			time.Sleep(10 * time.Millisecond)
			return gatewayappmax.TokenResponse{AccessToken: "app-token-abc", ExpiresIn: 3600}, nil
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "id", "secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			tok, err := mgr.AppToken(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, "app-token-abc", tok)
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, calls, "concurrent calls should coalesce into a single fetch")
}

func TestAppToken_ReturnsErrorOnGatewayFailure(t *testing.T) {
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, _, _ string) (gatewayappmax.TokenResponse, error) {
			return gatewayappmax.TokenResponse{}, errors.New("unauthorized")
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "bad-id", "bad-secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)
	_, err = mgr.AppToken(context.Background())

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrTokenFetch)
}

func TestMerchantToken_FetchesAndCaches(t *testing.T) {
	calls := 0
	cache := newMemoryTokenCache()
	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error) {
			calls++
			assert.Equal(t, "merchant-id", clientID)
			assert.Equal(t, "merchant-secret", clientSecret)
			return gatewayappmax.TokenResponse{AccessToken: "merchant-token", ExpiresIn: 3600}, nil
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "id", "secret", cache, &captureLogger{})
	require.NoError(t, err)

	inst := &models.Installation{ID: 10, MerchantClientID: "merchant-id", MerchantClientSecret: "merchant-secret"}
	token1, err := mgr.MerchantToken(context.Background(), inst)
	require.NoError(t, err)
	token2, err := mgr.MerchantToken(context.Background(), inst)
	require.NoError(t, err)

	assert.Equal(t, "merchant-token", token1)
	assert.Equal(t, "merchant-token", token2)
	assert.Equal(t, 1, calls)
	assert.Equal(t, "merchant_token:10", cache.lastKey)
	assert.Equal(t, "merchant-token", cache.lastData)
	assert.True(t, cache.lastTTL > 0)
}

func TestMerchantToken_WarnsWhenCachePutFails(t *testing.T) {
	cache := newMemoryTokenCache()
	cache.putErr = errors.New("redis down")
	logger := &captureLogger{}

	gateway := &mockTokenGateway{
		getTokenFunc: func(_ context.Context, _, _ string) (gatewayappmax.TokenResponse, error) {
			return gatewayappmax.TokenResponse{AccessToken: "merchant-token", ExpiresIn: 3600}, nil
		},
	}

	mgr, err := services.NewTokenManagerWithGatewayDeps(gateway, "id", "secret", cache, logger)
	require.NoError(t, err)

	inst := &models.Installation{ID: 15, MerchantClientID: "merchant-id", MerchantClientSecret: "merchant-secret"}
	token, err := mgr.MerchantToken(context.Background(), inst)
	require.NoError(t, err)
	assert.Equal(t, "merchant-token", token)
	assert.NotEmpty(t, logger.warnings)
}

func TestMerchantToken_RejectsNilInstallation(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGatewayDeps(&mockTokenGateway{}, "id", "secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)

	_, err = mgr.MerchantToken(context.Background(), nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

func TestMerchantToken_RejectsMissingCredentials(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGatewayDeps(&mockTokenGateway{}, "id", "secret", newMemoryTokenCache(), &captureLogger{})
	require.NoError(t, err)

	_, err = mgr.MerchantToken(context.Background(), &models.Installation{ID: 22})

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrInvalidConfig)
}

func TestTokenManagerConstructor_RejectsInvalidConfig(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGatewayDeps(&mockTokenGateway{}, "", "secret", newMemoryTokenCache(), &captureLogger{})

	require.Error(t, err)
	assert.Nil(t, mgr)
	assert.ErrorIs(t, err, services.ErrInvalidConfig)
}

func TestTokenManagerConstructor_RejectsNilGateway(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGatewayDeps(nil, "id", "secret", newMemoryTokenCache(), &captureLogger{})

	require.Error(t, err)
	assert.Nil(t, mgr)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

func TestTokenManagerConstructor_WithGatewayWrapperSuccess(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGateway(&mockTokenGateway{}, "id", "secret")

	require.NoError(t, err)
	assert.NotNil(t, mgr)
}

func TestTokenManagerConstructor_RejectsNilClock(t *testing.T) {
	mgr, err := services.NewTokenManagerWithGatewayDepsAndClock(
		&mockTokenGateway{},
		"id",
		"secret",
		newMemoryTokenCache(),
		&captureLogger{},
		nil,
	)

	require.Error(t, err)
	assert.Nil(t, mgr)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}
