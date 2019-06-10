package src

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MySQL struct {
	DB *sql.DB
}

// Init MySql connection
func NewMySQL(dbUser, dbPassword, dbHost, dbName string) (*MySQL, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		dbUser, dbPassword, dbHost, dbName,
	)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		return &MySQL{}, err
	}

	return &MySQL{DB: db}, nil
}

func (db *MySQL) GetFrames(id int) ([]Frame, error) {

	frames := []Frame{}

	rows, err := db.DB.Query("SELECT * FROM gpsTrack WHERE id >= ?", id)

	if err != nil {

		return []Frame{}, errors.Wrap(err, "Got an error in SELECT query.")
	}
	for rows.Next() {

		f := Frame{}
		err := rows.Scan(&f.Id, &f.Time, &f.Longitude, &f.Latitude, &f.Status, &f.LatitudeHemisphere, &f.LongitudeHemisphere,
			&f.EarthVelocity, &f.Track, &f.Date, &f.MagneticVariation, &f.DirectionVariation, &f.SystemPosition, &f.Checksum)

		if err != nil {

			return []Frame{}, errors.Wrap(err, "Got an error in SELECT query.")
		}

		frames = append(frames, f)
	}

	return frames, nil
}
