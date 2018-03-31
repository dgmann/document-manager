package shared

type ImageResult struct {
	PageNumber uint `json:"pageNumber"`
	Image []byte `json:"image"`
}
