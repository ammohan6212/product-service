package gcsclient

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
)

var (
	Client     *storage.Client
	Ctx        context.Context
	BucketName = "ecommerce-details"
)

func ConnectGCS() error {
	Ctx = context.Background()
	var err error
	Client, err = storage.NewClient(Ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %v", err)
	}
	return nil
}

func CloseGCS() {
	if Client != nil {
		Client.Close()
	}
}
