package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasbrake/actracker/model"
	"github.com/jasbrake/actracker/state"
)

// Start creates the API web handler
func Start(c *model.Config) {
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
		errChan <- R.Run(fmt.Sprintf(":%d", c.Port))
	}()

	err := <-errChan
	if err != nil {
		panic(err)
	}
}
