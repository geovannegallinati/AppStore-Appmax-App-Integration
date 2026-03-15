package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type UpsertInstallationInput struct {
	AppID                string
	ExternalKey          string
	MerchantClientID     string
	MerchantClientSecret string
}

type InstallService interface {
	Upsert(ctx context.Context, input UpsertInstallationInput) (inst *models.Installation, created bool, err error)
}

type installService struct {
	installRepo contracts.InstallationRepository
}

var _ InstallService = (*installService)(nil)

func NewInstallService(installRepo contracts.InstallationRepository) (InstallService, error) {
	if installRepo == nil {
		return nil, fmt.Errorf("new install service: %w", ErrNilDependency)
	}

	return &installService{installRepo: installRepo}, nil
}

func (s *installService) Upsert(ctx context.Context, input UpsertInstallationInput) (*models.Installation, bool, error) {
	existing, err := s.installRepo.FindByExternalKey(ctx, input.ExternalKey)
	if err != nil {
		return nil, false, fmt.Errorf("install upsert lookup: %w", err)
	}

	if existing != nil {
		existing.MerchantClientID = input.MerchantClientID
		existing.MerchantClientSecret = input.MerchantClientSecret
		existing.InstalledAt = time.Now()
		if err := s.installRepo.Save(ctx, existing); err != nil {
			return nil, false, fmt.Errorf("install upsert save: %w", err)
		}
		return existing, false, nil
	}

	inst := &models.Installation{
		AppID:                input.AppID,
		ExternalKey:          input.ExternalKey,
		MerchantClientID:     input.MerchantClientID,
		MerchantClientSecret: input.MerchantClientSecret,
		ExternalID:           uuid.New().String(),
		InstalledAt:          time.Now(),
	}
	if err := s.installRepo.Create(ctx, inst); err != nil {
		return nil, false, fmt.Errorf("install upsert create: %w", err)
	}
	return inst, true, nil
}
