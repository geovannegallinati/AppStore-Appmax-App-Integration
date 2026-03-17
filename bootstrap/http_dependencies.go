package bootstrap

import (
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
	repocontracts "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/repositories/contracts"
)

type HTTPDependencies struct {
	HealthController       *controllers.HealthController
	InstallController      *controllers.InstallController
	MerchantAuthController *controllers.MerchantAuthController
	CheckoutController     *controllers.CheckoutController
	WebhookController      *controllers.WebhookController
	InstallationRepository repocontracts.InstallationRepository
}

func NewHTTPDependencies() (*HTTPDependencies, error) {
	cfg, err := LoadAppmaxConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewHTTPDependenciesWithORM(cfg, facades.Orm())
}

func NewHTTPDependenciesWithORM(cfg AppmaxConfig, orm repocontracts.ORM) (*HTTPDependencies, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if orm == nil {
		return nil, ErrNilDependency
	}

	gateways, err := NewGatewayModule(cfg)
	if err != nil {
		return nil, err
	}

	repositories, err := NewRepositoryModule(orm)
	if err != nil {
		return nil, err
	}

	services, err := NewServiceModule(cfg, gateways, repositories)
	if err != nil {
		return nil, err
	}

	ctrlModule, err := NewControllerModule(cfg, services)
	if err != nil {
		return nil, err
	}

	return &HTTPDependencies{
		HealthController:       controllers.NewHealthController(cfg.AppPublicURL),
		InstallController:      ctrlModule.InstallController,
		MerchantAuthController: ctrlModule.MerchantAuthController,
		CheckoutController:     ctrlModule.CheckoutController,
		WebhookController:      ctrlModule.WebhookController,
		InstallationRepository: repositories.InstallationRepository,
	}, nil
}
