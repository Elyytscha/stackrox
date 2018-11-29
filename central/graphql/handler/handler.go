package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/stackrox/rox/central/graphql/resolvers"
	"github.com/stackrox/rox/central/graphql/schema"
)

// Handler returns an HTTP handler for the graphql api endpoint
func Handler() http.Handler {
	s := schema.Schema()
	ourSchema, err := graphql.ParseSchema(s, resolvers.New())
	if err != nil {
		fmt.Fprintf(os.Stderr, "s: %q", s)
		panic(err)
	}
	return &relay.Handler{Schema: ourSchema}
}
