package support

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/geovanne-gallinati/AppStoreAppDemo/config"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/foundation"
	frameworklog "github.com/goravel/framework/log"
	frameworkvalidation "github.com/goravel/framework/validation"
	frameworkview "github.com/goravel/framework/view"
	goravel_gin "github.com/goravel/gin"
	goravel_postgres "github.com/goravel/postgres"
	goravel_redis "github.com/goravel/redis"
)

func RunWithFrameworkBootstrap(m *testing.M) int {
	applyTestConfigOverrides()

	foundation.Setup().
		WithProviders(func() []contractsfoundation.ServiceProvider {
			return []contractsfoundation.ServiceProvider{
				&frameworklog.ServiceProvider{},
				&frameworkvalidation.ServiceProvider{},
				&frameworkview.ServiceProvider{},
				&goravel_gin.ServiceProvider{},
				&goravel_postgres.ServiceProvider{},
				&goravel_redis.ServiceProvider{},
			}
		}).
		Create()

	return m.Run()
}

func applyTestConfigOverrides() {
	config := facades.Config()
	config.Add("cache", map[string]any{
		"default": "memory",
		"prefix":  "appmax_checkout_",
		"stores": map[string]any{
			"memory": map[string]any{
				"driver": "memory",
			},
			"redis": map[string]any{
				"driver":     "redis",
				"connection": "default",
			},
		},
	})

	config.Add("logging", map[string]any{
		"default": "stack",
		"channels": map[string]any{
			"stack": map[string]any{
				"driver":   "stack",
				"channels": []string{"single"},
			},
			"single": map[string]any{
				"driver": "single",
				"path":   filepath.Join(os.TempDir(), fmt.Sprintf("appstoreappdemo-unit-%d.log", os.Getpid())),
				"level":  "debug",
				"days":   1,
				"print":  false,
			},
		},
	})
}
