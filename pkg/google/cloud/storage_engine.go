package cloud

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/anhdht/go-common/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"io"
	"strings"
)

type (
	StorageEngine interface {
		DownloadFile(ctx context.Context, bucket string, path string) ([]byte, error)
		UploadFile(ctx context.Context, bucket string, path string, data []byte) error
		DeleteFile(ctx context.Context, bucket string, path string) error
		ListFile(ctx context.Context, bucket string, path string) ([]string, error)
	}

	storageEngine struct {
		client *storage.Client
	}
)

func NewStorageEngine() StorageEngine {
	ctx := context.Background()

	logger := log.Logger(ctx)

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Fatal("encounter error while creating new storage client", zap.Error(err))
	}

	return &storageEngine{
		client: client,
	}
}

func (ptr *storageEngine) DownloadFile(ctx context.Context, bucket string, path string) ([]byte, error) {
	logger := log.Logger(ctx)

	b := ptr.client.Bucket(bucket)
	reader, err := b.Object(path).NewReader(ctx)
	if err != nil {
		logger.Error("encounter error while downloading data from: "+path, zap.Error(err))
		return nil, err
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		logger.Error("encounter error reading data from: "+path, zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

func (ptr *storageEngine) UploadFile(ctx context.Context, bucket string, path string, data []byte) error {
	logger := log.Logger(ctx)

	b := ptr.client.Bucket(bucket)
	writer := b.Object(path).NewWriter(ctx)
	l := len(data)
	count := 0
	for count < l {
		wrote, err := writer.Write(data[count:])
		if err != nil {
			logger.Error("encounter error while writing file: "+path, zap.Error(err))
			return err
		}
		count += wrote
	}

	if err := writer.Close(); err != nil {
		logger.Error("encounter error while closing file: "+path, zap.Error(err))
		return err
	}

	return nil
}

func (ptr *storageEngine) DeleteFile(ctx context.Context, bucket string, path string) error {
	logger := log.Logger(ctx)

	b := ptr.client.Bucket(bucket)
	err := b.Object(path).Delete(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("encounter error while deleting object in path: %s in bucket: %s",
			path, bucket), zap.Error(err))
		return err
	}
	return nil
}

func (ptr *storageEngine) ListFile(ctx context.Context, bucket string, path string) ([]string, error) {
	logger := log.Logger(ctx)

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	path = strings.TrimPrefix(path, "/")
	logger.Sugar().Infof("travelling to gs://%s/%s", bucket, path)

	b := ptr.client.Bucket(bucket)
	it := b.Objects(ctx, &storage.Query{
		Prefix:    path,
		Delimiter: "/",
	})

	var files []string

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return files, err
		}

		if strings.HasPrefix(attrs.ContentType, "image/") {
			f := strings.TrimPrefix(attrs.Name, path)
			files = append(files, f)
		}
	}

	return files, nil
}
