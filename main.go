package main

import (
	"github.com/goravel/framework/foundation"

	goravel_cache      "github.com/goravel/framework/cache"
	goravel_database   "github.com/goravel/framework/database"
	goravel_log        "github.com/goravel/framework/log"
	goravel_route      "github.com/goravel/framework/route"
	goravel_validation "github.com/goravel/framework/validation"
	goravel_view       "github.com/goravel/framework/view"
	goravel_gin        "github.com/goravel/gin"
	goravel_postgres   "github.com/goravel/postgres"
	goravel_redis      "github.com/goravel/redis"

	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	contractsschema     "github.com/goravel/framework/contracts/database/schema"

	_ "github.com/geovanne-gallinati/AppStoreAppDemo/config"
	"github.com/geovanne-gallinati/AppStoreAppDemo/database/migrations"
	"github.com/geovanne-gallinati/AppStoreAppDemo/routes"
)

func main() {
	foundation.Setup().
		WithProviders(func() []contractsfoundation.ServiceProvider {
			return []contractsfoundation.ServiceProvider{
				&goravel_log.ServiceProvider{},
				&goravel_cache.ServiceProvider{},
				&goravel_database.ServiceProvider{},
				&goravel_route.ServiceProvider{},
				&goravel_validation.ServiceProvider{},
				&goravel_view.ServiceProvider{},
				&goravel_gin.ServiceProvider{},
				&goravel_postgres.ServiceProvider{},
				&goravel_redis.ServiceProvider{},
			}
		}).
		WithRouting(func() {
			routes.Api()
		}).
		WithMigrations(func() []contractsschema.Migration {
			return migrations.All()
		}).
		Create().
		Start()
}
