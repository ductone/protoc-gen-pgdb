package v1

import (
	"github.com/doug-martin/goqu/v9/exp"
)

type DBReflectMessage interface {
	DBReflect(dialect Dialect) Message
}

type Message interface {
	Descriptor() Descriptor

	Record(opts ...RecordOptionsFunc) (exp.Record, error)
	SearchData(opts ...RecordOptionsFunc) []*SearchContent
	Dialect() Dialect
}

// PrimaryKeyer is implemented by the reflected Message of every top-level
// (non-nested-only) proto — i.e. every message that maps to its own Postgres
// row and therefore has a primary key. Nested-only messages are stored inline
// and have no key, so they do NOT implement it; type-assert before calling.
type PrimaryKeyer interface {
	// PKSK returns the row's primary key as stored in Postgres: the pb$pk and
	// pb$sk column values joined with "|" (matching the pb$pksk generated
	// column, `pb$pk || '|' || pb$sk`). It is derived only from the key fields —
	// no proto marshal, FTS, or search-data computation — so it is much cheaper
	// than pulling the value out of Record().
	PKSK() string
}

type ColumnExpression interface {
	Identifier() exp.IdentifierExpression
}
