package api

import (
	"errors"
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

	R.GET("/player", func(c *gin.Context) {
		name := c.Query("name")
		if len(name) <= 0 {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("missing 'name' query parameter"))
			return
		}

		gamePlayers, err := model.GetPlayerGames(name)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gamePlayers)
	})

	R.GET("/player_autocomplete", func(c *gin.Context) {
		name := c.Query("q")
		if len(name) <= 0 {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("missing 'q' query parameter"))
			return
		}

		names, err := model.GetPlayerNames(name)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, names)
	})

	go func() {
		errChan <- R.Run(fmt.Sprintf(":%d", c.Port))
	}()

	err := <-errChan
	if err != nil {
		panic(err)
	}
}
