package pdfprocessor

import (
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"fmt"
	log "github.com/sirupsen/logrus"
	"bytes"
	"image/jpeg"
	"io"
)

type ImageResult struct {
	PageNumber uint `json:"pageNumber"`
	Image []byte `json:"image"`
}

func ToImages(file io.ReadSeeker) []ImageResult {
	pdfReader, err := pdf.NewPdfReader(file)
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Panic(err)
	}

	var images []ImageResult
	for i := 0; i < numPages; i++ {
		log.Debugf("Page %d:", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			log.Panic(err)
		}

		// List images on the page.
		rgbImages, err := extractImagesOnPage(page)
		if err != nil {
			log.Panic(err)
		}
		_ = rgbImages

		for idx, img := range rgbImages {
			gimg, err := img.ToGoImage()
			if err != nil {
				log.Panic(err)
			}

			buf := new(bytes.Buffer)
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(buf, gimg, &opt)
			if err != nil {
				log.Panic(err)
			}
			images = append(images, ImageResult{PageNumber:uint(idx), Image:buf.Bytes()})
		}
	}
	return images
}

func extractImagesOnPage(page *pdf.PdfPage) ([]*pdf.Image, error) {
	contents, err := page.GetAllContentStreams()
	if err != nil {
		return nil, err
	}

	return extractImagesInContentStream(contents, page.Resources)
}

func extractImagesInContentStream(contents string, resources *pdf.PdfPageResources) ([]*pdf.Image, error) {
	rgbImages := []*pdf.Image{}
	cstreamParser := pdfcontent.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return nil, err
	}

	processedXObjects := map[string]bool{}

	// Range through all the content stream operations.
	for _, op := range *operations {
		if op.Operand == "BI" && len(op.Params) == 1 {
			// BI: Inline image.

			iimg, ok := op.Params[0].(*pdfcontent.ContentStreamInlineImage)
			if !ok {
				continue
			}

			img, err := iimg.ToImage(resources)
			if err != nil {
				return nil, err
			}

			cs, err := iimg.GetColorSpace(resources)
			if err != nil {
				return nil, err
			}
			if cs == nil {
				// Default if not specified?
				cs = pdf.NewPdfColorspaceDeviceGray()
			}
			fmt.Printf("Cs: %T\n", cs)

			rgbImg, err := cs.ImageToRGB(*img)
			if err != nil {
				return nil, err
			}

			rgbImages = append(rgbImages, &rgbImg)
		} else if op.Operand == "Do" && len(op.Params) == 1 {
			// Do: XObject.
			name := op.Params[0].(*pdfcore.PdfObjectName)

			// Only process each one once.
			_, has := processedXObjects[string(*name)]
			if has {
				continue
			}
			processedXObjects[string(*name)] = true

			_, xtype := resources.GetXObjectByName(*name)
			if xtype == pdf.XObjectTypeImage {
				fmt.Printf(" XObject Image: %s\n", *name)

				ximg, err := resources.GetXObjectImageByName(*name)
				if err != nil {
					return nil, err
				}

				img, err := ximg.ToImage()
				if err != nil {
					return nil, err
				}

				rgbImg, err := ximg.ColorSpace.ImageToRGB(*img)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, &rgbImg)
			} else if xtype == pdf.XObjectTypeForm {
				// Go through the XObject Form content stream.
				xform, err := resources.GetXObjectFormByName(*name)
				if err != nil {
					return nil, err
				}

				formContent, err := xform.GetContentStream()
				if err != nil {
					return nil, err
				}

				// Process the content stream in the Form object too:
				formResources := xform.Resources
				if formResources == nil {
					formResources = resources
				}

				// Process the content stream in the Form object too:
				formRgbImages, err := extractImagesInContentStream(string(formContent), formResources)
				if err != nil {
					return nil, err
				}
				rgbImages = append(rgbImages, formRgbImages...)
			}
		}
	}

	return rgbImages, nil
}
