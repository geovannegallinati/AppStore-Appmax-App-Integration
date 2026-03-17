package bootstrap

import (
	"fmt"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
)

type ControllerModule struct {
	InstallController      *controllers.InstallController
	MerchantAuthController *controllers.MerchantAuthController
	CheckoutController     *controllers.CheckoutController
	WebhookController      *controllers.WebhookController
}

func NewControllerModule(cfg AppmaxConfig, services *ServiceModule) (*ControllerModule, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if services == nil {
		return nil, fmt.Errorf("new controller module: %w", ErrNilDependency)
	}

	installController, err := controllers.NewInstallController(services.AppmaxService, services.InstallService, cfg.AdminURL, cfg.AppPublicURL, cfg.AppIDUUID, cfg.AppIDNumeric)
	if err != nil {
		return nil, err
	}

	merchantAuthController, err := controllers.NewMerchantAuthController(services.TokenManager)
	if err != nil {
		return nil, err
	}

	checkoutController, err := controllers.NewCheckoutController(services.CheckoutService)
	if err != nil {
		return nil, err
	}

	webhookController, err := controllers.NewWebhookController(services.WebhookService)
	if err != nil {
		return nil, err
	}

	return &ControllerModule{
		InstallController:      installController,
		MerchantAuthController: merchantAuthController,
		CheckoutController:     checkoutController,
		WebhookController:      webhookController,
	}, nil
}
