package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type Note struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Contents   *string    `json:"contents"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedAt *time.Time `json:"modifiedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}

func (Note) IsNode() {}

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

func (db *firestoreDB) AddNote(ctx context.Context, title string, contents *string) (note Note, err error) {
	ref := db.client.Collection(db.collection).NewDoc()

	note = Note{
		ID:         ref.ID,
		Title:      title,
		Contents:   contents,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
		DeletedAt:  nil,
	}

	if _, err := ref.Create(ctx, note); err != nil {
		return Note{}, fmt.Errorf("create: %v", err)
	}
	return note, nil
}
