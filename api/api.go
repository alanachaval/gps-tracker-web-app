package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/alanachaval/gps-tracker-web-app/src"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Api struct {
	database *src.MySQL
	key      string
}

func newApi(db *src.MySQL, key string) *Api {

	return &Api{
		database: db,
		key:      key,
	}
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	storage, err := src.NewMySQL(os.Getenv("user"), os.Getenv("password"), os.Getenv("host"), os.Getenv("dbName"))

	if err != nil {
		errors.Wrap(err, "Could not establish connection with the database")
	}

	key := os.Getenv("privKey")
	api := newApi(storage, key)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/frames", api.getFrames)
	r.POST("/frame", api.postFrames)

	return r
}

func Start() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func (a *Api) getFrames(c *gin.Context) {
	response, err := a.database.GetFrames(1, 0)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *Api) postFrames(c *gin.Context) {
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	// PENDIENTE LEER JSON
	frames := strings.Split(reqBody, "\n")
	err := a.AddFramesToDB("GPSTrackerUser", frames)
	if err == nil {
		c.JSON(200, gin.H{})
	} else {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
	}
}

// AddFramesToDB insert the frames for the user
func (a *Api) AddFramesToDB(user string, frames []string) error {

	userID, err := a.database.GetUserID("GPSTrackerUser")
	if err != nil {
		return errors.Wrap(err, "Cant retrieve user")
	}
	for _, f := range frames {
		err := a.database.AddFrame(f, userID)
		if err != nil {
			return errors.Wrap(err, "Cant insert frames")
		}
	}
	return nil
}
