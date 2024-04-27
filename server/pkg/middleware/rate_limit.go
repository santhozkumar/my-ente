package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/santhozkumar/my-ente/pkg/utils/network"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimitMiddleware struct {
	limit10ReqPerMin  *limiter.Limiter
	limit200ReqPerSec *limiter.Limiter
}

func NewRateLimitMiddlware() *RateLimitMiddleware {

	return &RateLimitMiddleware{
		limit10ReqPerMin:  rateLimiter("10-M"),
		limit200ReqPerSec: rateLimiter("200-S"),
	}
}


func rateLimiter(formatted string) *limiter.Limiter {
	store := memory.NewStore()
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		panic(err)
	}

	instance := limiter.New(store, rate)
	return instance
}

func (r *RateLimitMiddleware) getRateLimit(path string) *limiter.Limiter {
    if path == "/users"{
        return r.limit10ReqPerMin
    }

    return r.limit200ReqPerSec
}


func (r *RateLimitMiddleware) APIRateLimitMiddleWare(urlsanitizer func(_ *gin.Context) string) func (*gin.Context){
    return func (c *gin.Context){
        urlPath := urlsanitizer(c)
        rateLimiter := r.getRateLimit(urlPath)
        if rateLimiter != nil {
            key := fmt.Sprintf("%s-%s",network.GetClientIP(c), urlPath)
            limitContent, err := rateLimiter.Increment(c, key, 1)
            if err != nil {
                log.Print("Failed to check rate limit", err)
                c.Next()
            }
            if limitContent.Reached {
                c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate Limit Reached, try later"})
                log.Print("rate limit reached")
                return 
            }
        }
        c.Next()
    }
}
