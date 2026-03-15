package bootstrap

import (
	"fmt"

	adaptergoravel "github.com/geovanne-gallinati/AppStoreAppDemo/app/adapters/goravel"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

type ServiceModule struct {
	TokenManager    services.TokenManager
	AppmaxService   services.AppmaxService
	InstallService  services.InstallService
	CheckoutService services.CheckoutService
	WebhookService  services.WebhookService
}

func NewServiceModule(cfg AppmaxConfig, gateways *GatewayModule, repositories *RepositoryModule) (*ServiceModule, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if gateways == nil || repositories == nil {
		return nil, fmt.Errorf("new service module: %w", ErrNilDependency)
	}

	logger := adaptergoravel.NewLogger()
	cache := adaptergoravel.NewCache()

	tokenManager, err := services.NewTokenManagerWithGatewayDeps(
		gateways.AppmaxGateway,
		cfg.AppClientID,
		cfg.AppClientSecret,
		cache,
		logger,
	)
	if err != nil {
		return nil, err
	}

	appmaxService, err := services.NewAppmaxServiceWithGateway(tokenManager, gateways.AppmaxGateway)
	if err != nil {
		return nil, err
	}

	installService, err := services.NewInstallService(repositories.InstallationRepository)
	if err != nil {
		return nil, err
	}

	checkoutService, err := services.NewCheckoutService(appmaxService, repositories.OrderRepository, logger)
	if err != nil {
		return nil, err
	}

	webhookService, err := services.NewWebhookService(repositories.WebhookEventRepository, repositories.OrderRepository, logger)
	if err != nil {
		return nil, err
	}

	return &ServiceModule{
		TokenManager:    tokenManager,
		AppmaxService:   appmaxService,
		InstallService:  installService,
		CheckoutService: checkoutService,
		WebhookService:  webhookService,
	}, nil
}
