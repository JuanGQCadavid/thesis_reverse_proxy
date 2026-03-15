package ports

import (
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type Cache interface {
	Get(path string) (*domain.CacheData, error)
	Save(path string, data *domain.CacheData) error
}
