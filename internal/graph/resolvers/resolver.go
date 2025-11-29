package resolvers

import (
	"go-graphql/internal/graph"
)

type Resolver struct {
	*graph.Resolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r.Resolver}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r.Resolver}
}

type queryResolver struct {
	*graph.Resolver
}

type mutationResolver struct {
	*graph.Resolver
}
