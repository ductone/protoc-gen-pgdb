package v1

type Dialect uint8

const (
	DialectUnspecified Dialect = iota
	DialectV13
	DialectV17
)

const DefaultDialect = DialectV13

func (d Dialect) String() string {
	switch d {
	case DialectV13:
		return "DIALECT_V13"
	case DialectV17:
		return "DIALECT_V17"
	case DialectUnspecified:
		fallthrough
	default:
		return "DIALECT_UNSPECIFIED"
	}
}

func DialectOrDefault(d Dialect) Dialect {
	switch d {
	case DialectV17:
		return DialectV17
	default:
		return DefaultDialect
	}
}
