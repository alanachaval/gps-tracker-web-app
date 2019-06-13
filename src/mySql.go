package src

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func (db *MySQL) GetAllWay() ([]Frame, error) {

	frames := []Frame{}

	rows, err := db.DB.Query("SELECT * FROM gpsTrack")

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

// Recieve frame in string format
func (db *MySQL) AddFrame(frame string) (Frame, error) {

	s := strings.Split(frame, ",")

	longitude, err := strconv.ParseFloat(s[2], 32)
	if err == nil {
		fmt.Println("Wrong!!")
	}
	latitude, err := strconv.ParseFloat(s[3], 32)
	if err == nil {
		fmt.Println("Wrong!!")
	}

	earthVelocity, err := strconv.ParseFloat(s[7], 32)
	if err == nil {
		fmt.Println("Wrong!!")
	}

	track, err := strconv.ParseFloat(s[8], 32)
	if err == nil {
		fmt.Println("Wrong!!")
	}

	magnetic, err := strconv.ParseFloat(s[10], 32)
	if err == nil {
		fmt.Println("Wrong!!")
	}

	newFrame := Frame{
		Time:                s[1],
		Longitude:           longitude,
		Latitude:            latitude,
		Status:              s[4],
		LatitudeHemisphere:  s[5],
		LongitudeHemisphere: s[6],
		EarthVelocity:       earthVelocity,
		Track:               track,
		Date:                s[9],
		MagneticVariation:   magnetic,
		DirectionVariation:  s[11],
		SystemPosition:      s[12],
	}

	result, err := db.DB.Exec(
		"INSERT INTO frames (time, longitude, latitude, status, latitudeHemisphere, longitudeHemisphere, earthVelocity, track, date, magneticVariation, directionVariation, systemPosition) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		newFrame.Time, newFrame.Longitude, newFrame.Latitude, newFrame.Status, newFrame.LatitudeHemisphere, newFrame.LongitudeHemisphere, newFrame.EarthVelocity,
		newFrame.Track, newFrame.Date, newFrame.MagneticVariation, newFrame.DirectionVariation, newFrame.SystemPosition)

	if err != nil {
		return Frame{}, errors.Wrap(err, "Could not insert frame into database")
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return Frame{}, errors.Wrap(err, "Could not get last inserted project ID")
	}

	newFrame.Id = lastID

	return newFrame, nil

}
