package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
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

func (db *firestoreDB) DeleteNote(ctx context.Context, id string) error {
	if _, err := db.client.Collection(db.collection).Doc(id).Delete(ctx); err != nil {
		return fmt.Errorf("firestore: delete: %v", err)
	}
	return nil
}

func (db *firestoreDB) ListNotes(ctx context.Context) ([]*Note, error) {
	notes := make([]*Note, 0)
	iter := db.client.Collection(db.collection).Query.OrderBy("CreatedAt", firestore.Desc).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("firestoredb: could not list books: %v", err)
		}
		note := &Note{}
		doc.DataTo(note)
		db.logger.Info("found note", zap.String("id", note.ID), zap.String("title", note.Title))
		notes = append(notes, note)
	}

	return notes, nil
}
