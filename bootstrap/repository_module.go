package bootstrap

import (
	"fmt"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories"
	repocontracts "github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type RepositoryModule struct {
	InstallationRepository repocontracts.InstallationRepository
	OrderRepository        repocontracts.OrderRepository
	WebhookEventRepository repocontracts.WebhookEventRepository
}

func NewRepositoryModule(orm repocontracts.ORM) (*RepositoryModule, error) {
	if orm == nil {
		return nil, fmt.Errorf("new repository module: %w", ErrNilDependency)
	}

	installationRepository, err := repositories.NewInstallationRepository(orm)
	if err != nil {
		return nil, err
	}

	orderRepository, err := repositories.NewOrderRepository(orm)
	if err != nil {
		return nil, err
	}

	webhookEventRepository, err := repositories.NewWebhookEventRepository(orm)
	if err != nil {
		return nil, err
	}

	return &RepositoryModule{
		InstallationRepository: installationRepository,
		OrderRepository:        orderRepository,
		WebhookEventRepository: webhookEventRepository,
	}, nil
}
