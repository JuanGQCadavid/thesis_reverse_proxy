package ports

import (
	"github.com/JuanGQCadavid/thesis_reverse_proxy/poxy/internal/core/domain"
)

type Cache interface {
	Get(key domain.HttpStatusLineMultipart) (*domain.CacheData, error)
	Save(key domain.HttpStatusLineMultipart, data *domain.CacheData) error
}
