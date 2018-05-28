package store

import (
	"bytes"
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

func NewGCS() ObjectStore {
	ctx := context.Background()

	// projectID := "beautifulthings-204814"

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bucketName := "bt-main"

	bucket := client.Bucket(bucketName)

	// if err := bucket.Create(ctx, projectID, nil); err != nil {
	// 	log.Fatalf("Failed to create bucket: %v", err)
	// }

	return &gcsStore{
		c: client,
		b: bucket,
	}
}

type gcsStore struct {
	c *storage.Client
	b *storage.BucketHandle
}

func (g *gcsStore) Get(url string) ([]byte, error) {
	obj := g.b.Object(url)
	ctx := context.Background()

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf.Bytes(), nil
}

func (g *gcsStore) Set(url string, val []byte) error {
	obj := g.b.Object(url)
	ctx := context.Background()

	w := obj.NewWriter(ctx)
	if _, err := w.Write(val); err != nil {
		return errors.WithStack(err)
	}

	if err := w.Close(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
