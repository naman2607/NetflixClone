//Route definitions and routing logic for mapping URLs to handlers.

package authRoute

import (
	"github.com/gin-gonic/gin"
	authHandler "github.com/naman2607/netflixClone/handlers"
)

func InitAuthRoute(r *gin.Engine) {
	r.POST("/signup", func(c *gin.Context) {
		authHandler.Signup(c)
		c.Request.Body.Close()
	})

	r.POST("/login", func(c *gin.Context) {
		authHandler.Signin(c)
		c.Request.Body.Close()
	})
}
