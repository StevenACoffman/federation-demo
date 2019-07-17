// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package products

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql/introspection"
)

func (ec *executionContext) __resolve__service(ctx context.Context) (introspection.Service, error) {
	if ec.DisableIntrospection {
		return introspection.Service{}, errors.New("federated introspection disabled")
	}
	return introspection.Service{
		SDL: `type Product @key(fields: "upc price") {
	upc: String!
	name: String
	price: Int
	weight: Int
}
type Query @extends {
	topProducts(first: Int = 5): [Product]
}
`,
	}, nil
}

func (ec *executionContext) __resolve_entities(ctx context.Context, representations []map[string]interface{}) ([]_Entity, error) {
	list := []_Entity{}
	for _, rep := range representations {
		typeName, ok := rep["__typename"].(string)
		if !ok {
			return nil, errors.New("__typename must be an existing string")
		}
		switch typeName {

		case "Product":
			id, ok := rep["upc"].(string)
			if !ok {
				return nil, errors.New("opsies")
			}
			resp, err := ec.resolvers.Product().ResolveEntity(ctx, id)
			if err != nil {
				return nil, err
			}
			list = append(list, resp)

		default:
			return nil, errors.New("unknown type: " + typeName)
		}
	}
	return list, nil
}
