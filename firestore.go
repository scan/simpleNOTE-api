package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type firestoreDB struct {
	client     *firestore.Client
	logger     *zap.Logger
	collection string
}

func newFirestoreDB(client *firestore.Client, logger *zap.Logger) (*firestoreDB, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the Firestore
	// service.
	err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("firestoredb: could not connect: %v", err)
	}
	return &firestoreDB{
		client:     client,
		collection: "notes",
	}, nil
}

func (db *firestoreDB) Close() error {
	return db.client.Close()
}
