package middleware

// import (
// 	"github.com/gin-gonic/gin"
// )

// // CorsConfig holds CORS configuration
// type CorsConfig struct {
// 	AllowOrigins     []string
// 	AllowMethods     []string
// 	AllowHeaders     []string
// 	ExposeHeaders    []string
// 	AllowCredentials bool
// }

// // DefaultCorsConfig returns the default CORS configuration
// func DefaultCorsConfig() CorsConfig {
// 	return CorsConfig{
// 		AllowOrigins: []string{"*"},
// 		AllowMethods: []string{
// 			"GET",
// 			"POST",
// 			"PUT",
// 			"PATCH",
// 			"DELETE",
// 			"OPTIONS",
// 		},
// 		AllowHeaders: []string{
// 			"Origin",
// 			"Content-Type",
// 			"Content-Length",
// 			"Accept-Encoding",
// 			"X-CSRF-Token",
// 			"Authorization",
// 			"Accept",
// 			"Cache-Control",
// 			"X-Requested-With",
// 		},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 	}
// }

// // Cors returns a middleware handler that handles CORS
// func Cors() gin.HandlerFunc {
// 	config := DefaultCorsConfig()
// 	return CorsWithConfig(config)
// }

// // CorsWithConfig returns a middleware handler with custom config
// func CorsWithConfig(config CorsConfig) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Set headers
// 		origin := c.Request.Header.Get("Origin")
// 		if origin != "" {
// 			for _, allowOrigin := range config.AllowOrigins {
// 				if allowOrigin == "*" || allowOrigin == origin {
// 					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
// 					break
// 				}
// 			}
// 		}

// 		// Set other headers
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods",
// 			joinStrings(config.AllowMethods))
// 		c.Writer.Header().Set("Access-Control-Allow-Headers",
// 			joinStrings(config.AllowHeaders))
// 		c.Writer.Header().Set("Access-Control-Expose-Headers",
// 			joinStrings(config.ExposeHeaders))

// 		// Handle preflight requests
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// // Helper function to join strings with comma
// func joinStrings(strings []string) string {
// 	if len(strings) == 0 {
// 		return ""
// 	}
// 	result := strings[0]
// 	for _, s := range strings[1:] {
// 		result += ", " + s
// 	}
// 	return result
// }
