package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasbrake/actracker/state"
)

// Start creates the API web handler
func Start() {
	errChan := make(chan error)
	R := gin.Default()

	if gin.Mode() == gin.DebugMode {
		R.Use(cors.Default())
	}

	R.GET("/servers", func(c *gin.Context) {
		servers := state.GetActiveServers()
		c.JSON(http.StatusOK, servers)
	})

	go func() {
		errChan <- R.Run(":3000")
	}()

	err := <-errChan
	if err != nil {
		panic(err)
	}
}
