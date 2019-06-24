package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/alanachaval/gps-tracker-web-app/src"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	r.POST("/frame", api.postFrames)

	return r
}

func Start() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080

	key := os.Getenv("privKey")
	cert := os.Getenv("cert")
	//r.Run(":443")

	err := r.RunTLS(":443", cert, key)
	if err != nil {
		fmt.Println("Could not start WebServer")
	}
}

func (a *Api) getFrames(c *gin.Context) {

	lastFrame := c.Query("lastTrack")
	user := c.Query("user")

	var err error
	var response []src.Frame

	if user == "" {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		userID, err := a.database.GetUserID(user)
		if err != nil {
			if lastFrame == "" {
				response, err = a.database.GetFrames(userID, 0)
			} else {
				lastFrameInt, err := strconv.ParseInt(lastFrame, 10, 64)
				if err != nil {
					response, err = a.database.GetFrames(userID, lastFrameInt)
				}
			}
		}
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *Api) postFrames(c *gin.Context) {
	bodyRequest := src.FramesDTO{}

	err := c.Bind(&bodyRequest)
	err = a.AddFramesToDB(bodyRequest.User, bodyRequest.Frames)
	if err == nil {
		c.JSON(200, gin.H{"msg": "Ok"})
	} else {
		c.JSON(400, gin.H{
			"error_msg": err.Error(),
		})
	}
}

// AddFramesToDB insert the frames for the user
func (a *Api) AddFramesToDB(user string, frames []src.Frame) error {

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
