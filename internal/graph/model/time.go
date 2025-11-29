package model

import (
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalTime converts time.Time to a GraphQL string (RFC3339).
func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalString(t.Format(time.RFC3339))
}

// UnmarshalTime parses a GraphQL string into time.Time.
func UnmarshalTime(v interface{}) (time.Time, error) {
	str, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("time must be a string")
	}
	return time.Parse(time.RFC3339, str)
}
