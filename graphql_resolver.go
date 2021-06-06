package main

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"go.uber.org/zap"
)

type Resolver struct {
	logger *zap.Logger
	db     *firestoreDB
}

func (r *accountResolver) Notes(ctx context.Context, obj *Account, first int, skip int, after *string) (*NoteConnection, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateNote(ctx context.Context, note NoteInput) (*Note, error) {
	panic("not implemented")
}

func (r *mutationResolver) EditNote(ctx context.Context, id string, note NoteInput) (*Note, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveNote(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *queryResolver) Me(ctx context.Context) (*Account, error) {
	panic("not implemented")
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
