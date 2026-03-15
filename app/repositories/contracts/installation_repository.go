package contracts

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type InstallationRepository interface {
	FindByExternalKey(ctx context.Context, key string) (*models.Installation, error)
	Create(ctx context.Context, inst *models.Installation) error
	Save(ctx context.Context, inst *models.Installation) error
}
