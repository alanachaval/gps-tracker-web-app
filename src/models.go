package src

type Frame struct {
	Id                  int
	Time                string
	Longitude           float32
	Latitude            float32
	Status              string
	LatitudeHemisphere  string
	LongitudeHemisphere string
	EarthVelocity       float32
	Track               float32
	Date                string
	MagneticVariation   float32
	DirectionVariation  string
	SystemPosition      string
	Checksum            string
}
