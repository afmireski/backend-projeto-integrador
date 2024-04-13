package services

import (
	"github.com/afmireski/garchop-api/internal/entities"
	"github.com/afmireski/garchop-api/internal/ports"

	myTypes "github.com/afmireski/garchop-api/internal/types"
	customErrors "github.com/afmireski/garchop-api/internal/errors"
	"github.com/afmireski/garchop-api/internal/validators"
)

type TiersService struct {
	repository ports.TiersRepositoryPort
}

func NewTiersService(repository ports.TiersRepositoryPort) *TiersService {
	return &TiersService{
		repository: repository,
	}
}

func (s *TiersService) FindAll() ([]entities.Tier, *customErrors.InternalError) {

	where := myTypes.Where{}
	data, err := s.repository.FindAll(where); if err != nil {
		return nil, customErrors.NewInternalError("a failure occurred when try to find the tiers", 500, []string{err.Error()})
	}

	response := entities.BuildTiersFromModels(data)
	return response, nil
}

func (s *TiersService) FindById(id string) (*entities.Tier, *customErrors.InternalError) {
	if !validators.IsValidUuid(id) {
		return nil, customErrors.NewInternalError("invalid id", 400, []string{"the id must be a valid uuid"})
	}

	where := myTypes.Where{
		"deleted_at": map[string]string{"is": "null"},
	}

	data, err := s.repository.FindById(id, where); if err != nil {
		return nil, customErrors.NewInternalError("a failure occurred when try to find the tiers", 500, []string{err.Error()})
	} else if data == nil {
		return nil, customErrors.NewInternalError("tier not found", 404, []string{})
	}
	response := entities.BuildTierFromModel(*data)
	
	return &response, nil
}


