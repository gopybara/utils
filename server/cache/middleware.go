package cache

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chloyka/ginannot"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gopybara/utils/http_utils"
	redisProvider "github.com/gopybara/utils/providers/redis"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CacheMiddlewareInterface struct {
	DropCache ginannot.Route      `gin:"DELETE /cache"`
	Cache     ginannot.Middleware `middleware:"name=cache"`
}

type CacheMiddleware struct {
	CacheMiddlewareInterface
	redisProvider *redisProvider.RedisProvider
	log           *zap.Logger
	cacheConfig   *CacheConfig
	localStore    sync.Map
}

func NewCacheMiddleware(
	redisProvider *providers.RedisProvider,
	log *zap.Logger,
	cacheConfig ...*CacheConfig,
) ginannot.Handler {
	if len(cacheConfig) == 0 {
		ttl := time.Hour * 48
		cacheConfig = []*CacheConfig{
			{
				TTL: ttl,
			},
		}
	}

	middleware := &CacheMiddleware{
		redisProvider: redisProvider,
		log:           log,
		cacheConfig:   cacheConfig[0],
		localStore:    sync.Map{},
	}
	go middleware.revalidateLocalCache()

	return middleware
}

func (m *CacheMiddleware) Cache(ctx *gin.Context) {
	cacheKey := m.getCacheKey(ctx)

	// Firstly check local cache store
	if entry := m.getFromLocal(cacheKey); entry != nil {
		// Add headers from cache
		m.addHeadersFromCache(ctx, entry)

		ctx.String(entry.StatusCode, string(entry.Body))
		ctx.Abort()
		return
	}

	// Secondly check redis cache store
	if entry := m.getFromRedis(cacheKey); entry != nil {
		// Add headers from cache
		m.addHeadersFromCache(ctx, entry)

		// Store in local cache store
		entry.Headers[HeaderXCacheValidUntilKey] = []string{time.Now().Add(m.redisProvider.Client().TTL(context.Background(), cacheKey).Val()).String()}
		m.localStore.Store(cacheKey, entry)
		ctx.String(entry.StatusCode, string(entry.Body))
		ctx.Abort()
		return
	}

	w := &writer{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = w
	ctx.Next()

	if ctx.Writer.Status() != 200 {
		ctx.Abort()
		return
	}

	k, exists := ctx.Get(StatusSkipCacheKey)
	if exists && k == StatusSkipCache {
		ctx.Abort()
		return
	}

	ctx.Writer.Header().Set(HeaderXCache, HeaderXCacheMiss)
	cache := &httpCacheItem{
		StatusCode: ctx.Writer.Status(),
		Headers:    ctx.Writer.Header(),
		Body:       w.body.Bytes(),
	}

	// Set cache valid until
	m.SetCache(cacheKey, cache)
}

func (m *CacheMiddleware) getFromLocal(key string) *httpCacheItem {
	if data, ok := m.localStore.Load(key); ok {
		entry := data.(*httpCacheItem)
		return entry
	}

	return nil
}

func (m *CacheMiddleware) getFromRedis(key string) *httpCacheItem {
	if data, err := m.redisProvider.Client().Get(context.Background(), key).Bytes(); err == nil {
		entry := &httpCacheItem{}
		if err := json.Unmarshal(data, entry); err != nil {
			m.log.Error("cannot unmarshal cache entry", zap.Error(err))
			return nil
		}

		return entry
	}

	return nil
}

func (m *CacheMiddleware) getCacheKey(ctx *gin.Context) string {
	cacheKey := ctx.Request.URL.Path

	if m.cacheConfig.KeyFunc != nil {
		cacheKey = m.cacheConfig.KeyFunc(ctx)
	}

	return cacheKey
}

func (m *CacheMiddleware) addHeadersFromCache(ctx *gin.Context, entry *httpCacheItem) {
	for k, h := range entry.Headers {
		for _, v := range h {
			ctx.Writer.Header().Add(k, v)
		}
	}

	ctx.Writer.Header().Set(HeaderXCache, HeaderXCacheHit)
}

type InvalidateRequest struct {
	Key string `json:"key"`
}

func (m *CacheMiddleware) Invalidate(ctx *gin.Context) {
	var request InvalidateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		http_utils.ThrowHTTPError(ctx, err)
		return
	}

	m.localStore.Range(func(key, value interface{}) bool {
		if strings.Contains(key.(string), request.Key) {
			m.localStore.Delete(key.(string))
			m.redisProvider.Client().Del(context.Background(), fmt.Sprintf("*%s*", key.(string)))
		}

		return true
	})
}

func (m *CacheMiddleware) SetCache(key string, entry *httpCacheItem) {
	entry.Headers[HeaderXCacheValidUntilKey] = []string{strconv.Itoa(int(time.Now().Add(m.cacheConfig.TTL).Unix()))}
	m.localStore.Store(key, entry)
	redisData, err := json.Marshal(entry)
	if err != nil {
		m.log.Error("cannot marshal cache entry", zap.Error(err))
		return
	}
	_ = m.redisProvider.Client().Set(context.Background(), key, redisData, m.cacheConfig.TTL)
}

func (m *CacheMiddleware) revalidateLocalCache() {
	for {
		m.localStore.Range(func(key, value interface{}) bool {
			entry := value.(*httpCacheItem)
			if entry.Headers[HeaderXCacheValidUntilKey] == nil {
				return true
			}
			cacheTime, err := strconv.Atoi(entry.Headers[HeaderXCacheValidUntilKey][0])
			if err != nil {
				m.localStore.Delete(key)
			}
			if time.Now().Unix() > int64(cacheTime) {
				m.localStore.Delete(key)
			}

			return true
		})
		time.Sleep(time.Minute)
	}
}

func (m *CacheMiddleware) parseTime(value string) time.Time {
	t, err := time.Parse("2023-06-24 16:26:01.054309587 +0000 UTC m=+173568.738855828", value)
	if err != nil {
		m.log.Error("cannot parse time", zap.Error(err))
		return time.Time{}
	}

	return t
}

func (m *CacheMiddleware) DropCache(ctx *gin.Context) {
	go m.Invalidate(ctx)
	ctx.AbortWithStatus(204)
}
