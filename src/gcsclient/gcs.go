package gcsclient

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"log"
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

	// ✅ Log success
	log.Printf("✅ Connected to Google Cloud Storage bucket: %s\n", BucketName)
	return nil
}

func CloseGCS() {
	if Client != nil {
		Client.Close()
		log.Println("🛑 GCS client connection closed.")
	}
}
