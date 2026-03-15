package repositories

import (
	"context"
	"fmt"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type installationRepository struct {
	orm contracts.ORM
}

var _ contracts.InstallationRepository = (*installationRepository)(nil)

func NewInstallationRepository(orm contracts.ORM) (contracts.InstallationRepository, error) {
	if orm == nil {
		return nil, fmt.Errorf("new installation repository: %w", ErrNilORM)
	}

	return &installationRepository{orm: orm}, nil
}

func (r *installationRepository) FindByExternalKey(_ context.Context, key string) (*models.Installation, error) {
	var inst models.Installation
	err := r.orm.Query().Where("external_key = ?", key).First(&inst)
	if err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if inst.ID == 0 {
		return nil, nil
	}
	return &inst, nil
}

func (r *installationRepository) Create(_ context.Context, inst *models.Installation) error {
	return r.orm.Query().Create(inst)
}

func (r *installationRepository) Save(_ context.Context, inst *models.Installation) error {
	return r.orm.Query().Save(inst)
}
