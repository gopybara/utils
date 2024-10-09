package cache

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"time"
)

type CacheConfig struct {
	TTL     time.Duration
	KeyFunc func(ctx *gin.Context) string
}

type httpCacheItem struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

const (
	HeaderXCache              = "X-Cache"
	HeaderXCacheHit           = "HIT"
	HeaderXCacheSkip          = "SKIP"
	HeaderXCacheMiss          = "MISS"
	StatusSkipCacheKey        = "skip-cache"
	StatusSkipCache           = 499
	HeaderXCacheValidUntilKey = "X-Cache-Valid-Until"
)

type writer struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *writer) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
