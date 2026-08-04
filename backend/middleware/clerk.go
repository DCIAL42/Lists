package middleware

import (
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

func WrapClerkMiddleware(clerkMiddleWare func(http.Handler) http.Handler) gin.HandlerFunc {
	clerk.SetKey("sk_test_XsiGhkWZeft81IydOVTHATS9UQZcrNDXckEPlym9M6")

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

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())

		if !ok {
			fmt.Println("notok")
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("claims", claims)

		c.Next()
	}
}
