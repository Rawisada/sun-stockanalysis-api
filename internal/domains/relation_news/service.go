package relation_news

import (
	"errors"
	"strings"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type RelationNewsService interface {
	CreateRelations(symbol string, relationSymbols []string) error
}

type RelationNewsServiceImpl struct {
	repo repository.RelationNewsRepository
}

func NewRelationNewsService(repo repository.RelationNewsRepository) RelationNewsService {
	return &RelationNewsServiceImpl{repo: repo}
}

func (s *RelationNewsServiceImpl) CreateRelations(symbol string, relationSymbols []string) error {
	base := strings.TrimSpace(symbol)
	if base == "" {
		return errors.New("symbol is required")
	}

	unique := map[string]struct{}{}
	unique[base] = struct{}{}
	for _, rel := range relationSymbols {
		rel = strings.TrimSpace(rel)
		if rel == "" {
			continue
		}
		unique[rel] = struct{}{}
	}

	items := make([]models.RelationNews, 0, len(unique))
	for rel := range unique {
		items = append(items, models.RelationNews{
			Symbol:         base,
			RelationSymbol: rel,
			IsActive:       true,
		})
	}

	return s.repo.CreateMany(items)
}
