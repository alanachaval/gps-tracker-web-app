package main

import (
	"net/http"
	"os"

	"github.com/alanachaval/gps-tracker-web-app/src"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Api struct {
	database *src.MySQL
}

func newApi(db *src.MySQL) *Api {

	return &Api{
		database: db,
	}
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	storage, err := src.NewMySQL(os.Getenv("user"), os.Getenv("password"), os.Getenv("host"), os.Getenv("dbName"))

	if err != nil {
		errors.Wrap(err, "Could not establish connection with the database")
	}

	api := newApi(storage)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/frames", api.getFrames)

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func (a *Api) getFrames(c *gin.Context) {
	response, err := a.database.GetFrames(0)

	if err != nil {
		c.String(http.StatusInternalServerError, "Could not get projects.")
		return
	}

	c.JSON(http.StatusOK, response)
}
