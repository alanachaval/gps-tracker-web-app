package src

type Frame struct {
	Id                  int64
	UserID              int64
	TrackNumber         int64
	Time                string
	Longitude           float64
	Latitude            float64
	Status              string
	LatitudeHemisphere  string
	LongitudeHemisphere string
	EarthVelocity       float64
	Track               float64
	Date                string
	MagneticVariation   float64
	DirectionVariation  string
	SystemPosition      string
	Checksum            string
}
