package main

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"

	"google.golang.org/api/iterator"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Note struct {
	Key        *datastore.Key `json:"-" datastore:"__key__"`
	Title      string         `json:"title"`
	Contents   *string        `json:"contents" datastore:",noindex"`
	CreatedAt  time.Time      `json:"createdAt"`
	ModifiedAt *time.Time     `json:"modifiedAt"`
	DeletedAt  *time.Time     `json:"deletedAt"`
}

func (Note) IsNode() {}

func (note Note) ID() string {
	return note.Key.Name
}

type datastoreDB struct {
	client     *datastore.Client
	logger     *zap.Logger
	collection string
}

func newDB(client *datastore.Client, logger *zap.Logger) (*datastoreDB, error) {
	return &datastoreDB{
		client:     client,
		collection: "notes",
	}, nil
}

func (db *datastoreDB) Close() error {
	return db.client.Close()
}

func (db *datastoreDB) AddNote(ctx context.Context, title string, contents *string) (note Note, err error) {
	key := datastore.NameKey("Note", uuid.New().String(), nil)
	note = Note{
		Title:      title,
		Contents:   contents,
		CreatedAt:  time.Now(),
		ModifiedAt: nil,
		DeletedAt:  nil,
	}

	key, err = db.client.Put(ctx, key, &note)
	if err != nil {
		return Note{}, errors.Wrap(err, "adding note failed")
	}

	note.Key = key

	return note, nil
}

func (db *datastoreDB) DeleteNote(ctx context.Context, id string) error {
	key := datastore.NameKey("Node", id, nil)

	return db.client.Delete(ctx, key)
}

func (db *datastoreDB) ListNotes(ctx context.Context) ([]*Note, error) {
	notes := make([]*Note, 0)

	query := datastore.NewQuery("Note").Order("-CreatedAt")
	iter := db.client.Run(ctx, query)

	for {
		var note Note
		_, err := iter.Next(&note)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not list notes")
		}

		db.logger.Info("found note", zap.String("id", note.Key.String()), zap.String("title", note.Title))
		notes = append(notes, &note)
	}

	return notes, nil
}
