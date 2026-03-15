package domain

import "time"

type CacheData struct {
	TTL  time.Time
	Data *HttpPackage
}
