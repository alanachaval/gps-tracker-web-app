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

// GetFrames return way for user, starting in lastTrack
func (db *MySQL) GetFrames(userID int64, lastTrack int64) ([]Frame, error) {

	frames := []Frame{}

	rows, err := db.DB.Query("SELECT * FROM gpsTrack WHERE userID = ? and trackNumber > ? ORDER BY trackNumber ASC", userID, lastTrack)

	if err != nil {

		return []Frame{}, errors.Wrap(err, "Got an error in SELECT query.")
	}

	defer rows.Close()

	for rows.Next() {

		f := Frame{}
		err := rows.Scan(&f.Id, &f.UserID, &f.TrackNumber, &f.Time, &f.Longitude, &f.Latitude, &f.Status, &f.LatitudeHemisphere, &f.LongitudeHemisphere,
			&f.EarthVelocity, &f.Track, &f.Date, &f.MagneticVariation, &f.DirectionVariation, &f.SystemPosition, &f.Checksum)

		if err != nil {

			return []Frame{}, errors.Wrap(err, "Got an error in SELECT query.")
		}

		frames = append(frames, f)
	}

	return frames, nil
}

// GetAllWay return all way for user
func (db *MySQL) GetAllWay(userID int64) ([]Frame, error) {
	return db.GetFrames(userID, 0)
}

// AddFrame Recieve frame in string format
func (db *MySQL) AddFrame(newFrame Frame, userID int64) error {

	_, err := db.DB.Exec(
		"INSERT INTO gpsTrack (userId, trackNumber, time, longitude, latitude, status, latitudeHemisphere, longitudeHemisphere, earthVelocity, track, date, magneticVariation, directionVariation, systemPosition) SELECT * FROM (SELECT ? as a,? as b,? as c,? as d,? as e,? as f,? as g,? as h,? as i,? as j,? as k,? as l,? as m,? as n) AS tmp WHERE NOT EXISTS (SELECT userId FROM gpsTrack WHERE userId = ? AND trackNumber = ?) LIMIT 1",
		userID, newFrame.TrackNumber, newFrame.Time, newFrame.Longitude, newFrame.Latitude, newFrame.Status, newFrame.LatitudeHemisphere, newFrame.LongitudeHemisphere, newFrame.EarthVelocity,
		newFrame.Track, newFrame.Date, newFrame.MagneticVariation, newFrame.DirectionVariation, newFrame.SystemPosition, userID, newFrame.TrackNumber)

	if err != nil {
		return errors.Wrap(err, "Could not insert frame into database")
	}

	return nil
}

// GetUserID retrieves userID
func (db *MySQL) GetUserID(user string) (int64, error) {
	result, err := db.DB.Query("SELECT id from gpsframes.user WHERE user = ?", user)
	if err != nil {
		return 0, errors.Wrap(err, "Got an error in SELECT query.")
	}
	if !result.Next() {
		return 0, errors.Wrap(err, "User not exists.")
	}
	var userID int64
	err = result.Scan(&userID)
	if err != nil {
		return 0, errors.Wrap(err, "Got an error in SELECT query.")
	}
	return userID, nil
}
