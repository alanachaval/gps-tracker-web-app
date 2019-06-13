package api

import (
	"net/http"
	"os"

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
	response, err := a.database.GetFrames(0)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *Api) postFrames(c *gin.Context) {
	var frame string
	err := c.Bind(&frame)
	if err == nil {
		c.JSON(200, gin.H{

			//FALTA CODIGO QUE AGARRE EL BODY, LO PROCESE Y DEVUELVA LOS DATOS SIN ENCRIPTAR
		})
	} else {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
	}
}

func (a *Api) DecryptFrames(key string, encryptedFrames []string) ([]string, error) {

	cipherKey := []byte(a.key)

	frames := []string{}
	for _, f := range encryptedFrames {
		decoded, err := src.Decrypt(cipherKey, f)
		if err != nil {
			return nil, errors.Wrap(err, "Could not decrypt the frames")
		}
		frames = append(frames, decoded)
	}

	return frames, nil
}

func (a *Api) AddDecryptFramesToDB(key string, encryptedFrames []string) ([]src.Frame, error) {

	planeFrames, err := a.DecryptFrames(a.key, encryptedFrames)

	if err != nil {
		return nil, errors.Wrap(err, "Could not decrypt the frames")
	}

	finalFrames := []src.Frame{}
	var aFrame src.Frame
	for _, f := range planeFrames {
		aFrame, err = a.database.AddFrame(f)
		if err != nil {
			return nil, errors.Wrap(err, "Could not decrypt the frames")
		}

		finalFrames = append(finalFrames, aFrame)
	}

	return finalFrames, nil
}
