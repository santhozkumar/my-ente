package middleware

import "github.com/patrickmn/go-cache"


type AuthMiddleware struct {
    Cache  cache.Cache
    userAuthRepo 
}
