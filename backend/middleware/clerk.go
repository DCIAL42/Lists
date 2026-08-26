package middleware

import (
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

func WrapClerkMiddleware(clerkMiddleWare func(http.Handler) http.Handler) gin.HandlerFunc {
	clerkSecret, ok := os.LookupEnv("CLERK_SECRET_KEY")

	if !ok {
		panic("clerk secret key not found")
	}

	clerk.SetKey(clerkSecret)

	return func(c *gin.Context) {
		var called bool
		var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			c.Request = r
		})

		clerkMiddleWare(next).ServeHTTP(c.Writer, c.Request)

		if !called {
			c.Abort()
			return
		}

		c.Next()
	}
}

func MaybeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())

		if !ok {
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("claims", claims)

		c.Next()
	}
}

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())

		if !ok {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("claims", claims)

		c.Next()
	}
}
