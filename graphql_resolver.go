package main

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
)

type Resolver struct{}

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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
