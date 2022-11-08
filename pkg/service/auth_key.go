package service

import (
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type APIKeyStrategy interface {
	Create() (string, error)
	List(limit, offset int, status string) ([]model.AuthStrategy, error)
	Approve(ID string) error
}

type aPIKeyStrategy struct {
	repo repository.AuthStrategy
}

func NewAPIKeyStrategy(repo repository.AuthStrategy) APIKeyStrategy {
	return aPIKeyStrategy{repo}
}
func (g aPIKeyStrategy) Create() (string, error) {
	uuiKey := "str." + uuidWithoutHyphens()
	hashed := common.ToSha256(uuiKey)
	err := g.repo.CreateAPIKey("", repository.AuthTypeAPIKey, hashed, true)
	return uuiKey, err
}

func (g aPIKeyStrategy) List(limit, offset int, status string) ([]model.AuthStrategy, error) {
	log.Info().Int("limit", limit).Int("offset", offset).Str("status", status).Msg("List")
	if status != "" {
		return g.ListByStatus(limit, offset, status)
	}
	return g.repo.List(limit, offset)
}

func (g aPIKeyStrategy) ListByStatus(limit, offset int, status string) ([]model.AuthStrategy, error) {
	if limit == 0 {
		limit = 100
	}
	return g.repo.ListByStatus(limit, offset, status)
}

// Approve updates the APIKey status and creates an entry on redis
func (g aPIKeyStrategy) Approve(ID string) error {
	m, err := g.repo.UpdateStatus(ID, "active")
	if err != nil {
		return err
	}
	return g.repo.CreateAPIKey(m.ID, repository.AuthTypeAPIKey, m.Data, false)
}
