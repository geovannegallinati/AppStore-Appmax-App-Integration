package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

const InstallationContextKey = "installation"

func MerchantContext(installRepo contracts.InstallationRepository) http.Middleware {
	return func(ctx http.Context) {
		key := ctx.Request().Input("key")
		if key == "" {
			key = ctx.Request().Route("key")
		}

		inst, err := installRepo.FindByExternalKey(ctx.Context(), key)
		if err != nil || inst == nil {
			ctx.Request().AbortWithStatusJson(404, http.Json{
				"message": "installation not found",
			})
			return
		}

		ctx.WithValue(InstallationContextKey, inst)
		ctx.Request().Next()
	}
}
