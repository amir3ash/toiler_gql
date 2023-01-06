package scallers

import (
	"database/sql"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalNullTime(t sql.NullTime) graphql.Marshaler {
	if !t.Valid {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.Time.Format(time.RFC3339Nano)))
	})
}

func UnmarshalNullTime(v interface{}) (sql.NullTime, error) {
	if tmpStr, ok := v.(string); ok {
		if time, err := time.Parse(time.RFC3339Nano, tmpStr); err != nil {
			return sql.NullTime{
				Time:  time,
				Valid: true,
			}, nil
		}
	}
	return sql.NullTime{}, errors.New("time should be RFC3339Nano formatted string")
}
