package filesystem

import (
	"context"

	"github.com/dgmann/document-manager/internal/backend/storage"
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type ImageService struct {
	*DiskStorage
	tracer oteltrace.Tracer
}

func NewImageService(directory string) (*ImageService, error) {
	repository, err := NewDiskStorage(directory)
	if err != nil {
		return nil, err
	}
	tracer := otel.GetTracerProvider().Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)
	return &ImageService{DiskStorage: repository, tracer: tracer}, nil
}

func (f *ImageService) GetByRecordId(ctx context.Context, id string) (map[string]*storage.Image, error) {
	images := make(map[string]*storage.Image, 0)
	p := storage.NewKey(id)
	err := f.ForEach(ctx, p, func(resource storage.KeyedResource, err error) error {
		fileName := resource.Key()[len(resource.Key())-1]
		images[fileName] = storage.NewImage(resource.Data(), resource.Format())
		return err
	})
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (f *ImageService) Copy(ctx context.Context, fromId string, toId string) error {
	return f.CopyFolder(ctx, storage.NewKey(fromId), storage.NewKey(toId))
}
