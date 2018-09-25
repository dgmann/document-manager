package api

type Rotator interface {
	Rotate(img []byte, degrees float64) (*Image, error)
}
