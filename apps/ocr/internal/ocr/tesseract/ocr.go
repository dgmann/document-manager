//go:build !test

package tesseract

import (
	"fmt"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/dgmann/gosseract"
	"github.com/sirupsen/logrus"
	"ocr/internal/ocr"
)

type Client struct {
	client *gosseract.Client
}

func NewClient() *Client {
	return &Client{
		client: gosseract.NewClient(),
	}
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) CheckOrientation(pages []ocr.PageWithContent) ([]api.PageUpdate, error) {
	needsUpdate := false
	updatedPages := make([]api.PageUpdate, len(pages))
	for i, p := range pages {
		updatedPages[i] = api.PageUpdate{
			Id: p.Id,
		}

		// If there is no image to process, nothing to do
		if len(p.Image) == 0 {
			continue
		}

		if err := c.client.SetImageFromBytes(p.Image); err != nil {
			return nil, fmt.Errorf("error loading image: %w", err)
		}
		if err := c.client.SetLanguage("osd"); err != nil {
			return nil, fmt.Errorf("error setting language: %w", err)
		}
		if err := c.client.SetPageSegMode(gosseract.PSM_OSD_ONLY); err != nil {
			return nil, fmt.Errorf("error setting page segmentation mode: %w", err)
		}
		osdResult, err := c.client.DetectOrientationScript()
		if err != nil {
			return nil, err
		}
		if osdResult.OrientationDegree == 0 {
			continue
		}
		if osdResult.OrientationConfidence < 10.0 {
			logrus.Infof("orientation confidence too low. Confidence: %f.\n", osdResult.OrientationConfidence)
			continue
		}
		logrus.Infof("page %s is rotated %d degree with confidence %f. Correcting orientation.\n", p.Id, osdResult.OrientationDegree, osdResult.OrientationConfidence)
		needsUpdate = true
		updatedPages[i].Rotate = float64(osdResult.OrientationDegree * -1)
	}
	if !needsUpdate {
		return nil, nil
	}
	return updatedPages, nil
}

func (c *Client) ExtractText(pages []ocr.PageWithContent) ([]api.PageUpdate, error) {
	updatedPages := make([]api.PageUpdate, len(pages))
	for i, p := range pages {
		updatedPages[i] = api.PageUpdate{Id: p.Id}

		// If there is no image to process, we just append the page without modification
		if len(p.Image) == 0 {
			continue
		}
		if err := c.client.SetImageFromBytes(p.Image); err != nil {
			return nil, fmt.Errorf("error loading image: %w", err)
		}
		if err := c.client.SetLanguage("deu", "eng"); err != nil {
			return nil, fmt.Errorf("error setting language: %w", err)
		}
		if err := c.client.SetPageSegMode(gosseract.PSM_AUTO); err != nil {
			return nil, fmt.Errorf("error setting page segmentation mode: %w", err)
		}
		text, err := c.client.Text()
		if err != nil {
			return nil, err
		}

		updatedPages[i].Content = &text
	}
	return updatedPages, nil
}
