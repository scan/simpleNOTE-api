package main

import (
	"context"

	"go.uber.org/zap"
)

type Resolver struct {
	logger *zap.Logger
	db     *datastoreDB
}

func (r *accountResolver) Notes(ctx context.Context, obj *Account, first int, skip int, after *string) (*NoteConnection, error) {
	notes, err := r.db.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	edges := make([]*NoteEdge, len(notes))
	for i, note := range notes {
		edges[i] = &NoteEdge{
			Node:   note,
			Cursor: "",
		}
	}

	return &NoteConnection{
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage:     false,
			HasPreviousPage: false,
			EndCursor:       "",
			StartCursor:     "",
		},
	}, nil
}

func (r *mutationResolver) CreateNote(ctx context.Context, input NoteInput) (*Note, error) {
	note, err := r.db.AddNote(ctx, input.Title, input.Contents)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *mutationResolver) EditNote(ctx context.Context, id string, note NoteInput) (*Note, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveNote(ctx context.Context, id string) (bool, error) {
	err := r.db.DeleteNote(ctx, id)
	return err == nil, err
}

func (r *queryResolver) Me(ctx context.Context) (*Account, error) {
	return &Account{Email: ""}, nil
}

// Account returns AccountResolver implementation.
func (r *Resolver) Account() AccountResolver { return &accountResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type accountResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
