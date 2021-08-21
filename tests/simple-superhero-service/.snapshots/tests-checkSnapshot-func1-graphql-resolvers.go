package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/alexisvisco/graphql-proxy-grpc/tests/simple-superhero-service/gen/go/protos/superheroes"
)

// generated with custom plugin :)

type Resolver struct {
	SuperHeroService superheroes.SuperHeroServiceClient
}

func (r *mutationResolver) SuperHeroServiceNewSuperHero(ctx context.Context, in *superheroes.CallSuperHeroRequest) (*superheroes.CallSuperHeroResponse, error) {
	response, err := r.SuperHeroService.NewSuperHero(ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) SuperHeroServiceCallSuperHero(ctx context.Context, in *superheroes.CallSuperHeroRequest) (*superheroes.CallSuperHeroResponse, error) {
	response, err := r.SuperHeroService.CallSuperHero(ctx, in)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

