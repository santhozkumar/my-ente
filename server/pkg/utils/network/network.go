package network

import "github.com/gin-gonic/gin"

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("CF-Connecting-IP")
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}
