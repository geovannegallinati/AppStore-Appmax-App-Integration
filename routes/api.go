package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/middleware"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/bootstrap"
)

func Api() {
	healthController := controllers.NewHealthController(bootstrap.LoadAppPublicURLFromEnv())

	facades.Route().Static("public", "./public")
	facades.Route().Get("/", healthController.RootFrontend)
	facades.Route().Get("/health", healthController.Check)

	deps, err := bootstrap.NewHTTPDependencies()
	if err != nil {
		facades.Log().Warningf("routes api: bootstrap http dependencies failed: %v; registered only core health routes", err)
		return
	}

	facades.Route().Get("/integrations/appmax/callback/install", deps.InstallController.CallbackGuide)
	facades.Route().Post("/integrations/appmax/callback/install", deps.InstallController.Callback)

	facades.Route().Get("/install/start", deps.InstallController.Start)
	facades.Route().Get("/webhooks/appmax", deps.WebhookController.Guide)
	facades.Route().Post("/webhooks/appmax", deps.WebhookController.Handle)

	facades.Route().Prefix("/installations/{key}").Middleware(middleware.MerchantContext(deps.InstallationRepository)).Group(func(r route.Router) {
		r.Get("/merchant/token", deps.MerchantAuthController.SyncToken)
	})

	facades.Route().Prefix("/checkout/{key}").Middleware(middleware.MerchantContext(deps.InstallationRepository)).Group(func(r route.Router) {
		r.Post("/order", deps.CheckoutController.CreateOrder)
		r.Get("/status/{order_id}", deps.CheckoutController.Status)
		r.Get("/installments", deps.CheckoutController.Installments)
		r.Post("/pay/credit-card", deps.CheckoutController.PayCreditCard)
		r.Post("/pay/pix", deps.CheckoutController.PayPix)
		r.Post("/pay/boleto", deps.CheckoutController.PayBoleto)
		r.Post("/refund", deps.CheckoutController.Refund)
		r.Post("/tokenize", deps.CheckoutController.Tokenize)
		r.Post("/tracking", deps.CheckoutController.AddTracking)
		r.Post("/upsell", deps.CheckoutController.Upsell)
	})
}
