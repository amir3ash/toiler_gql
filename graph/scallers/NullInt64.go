package scallers

import (
	"database/sql"

	"github.com/99designs/gqlgen/graphql"
)


func MarshalNullInt64(ns sql.NullInt64) graphql.Marshaler {
	if !ns.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalInt64(ns.Int64)
}

// UnmarshalNullString is a custom unmarshaller.
func UnmarshalNullInt64(v interface{}) (sql.NullInt64, error) {
	if v == nil {
		return sql.NullInt64{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalInt64(v)
	return sql.NullInt64{Int64: s, Valid: true}, err
}