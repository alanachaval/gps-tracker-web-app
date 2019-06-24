package src

type Frame struct {
	Id                  int64
	UserID              int64
	TrackNumber         int64
	Time                string
	Longitude           string
	Latitude            string
	Status              string
	LatitudeHemisphere  string
	LongitudeHemisphere string
	EarthVelocity       string
	Track               string
	Date                string
	MagneticVariation   string
	DirectionVariation  string
	SystemPosition      string
	Checksum            string
}

type FramesDTO struct {
	User   string
	Frames []Frame
}
