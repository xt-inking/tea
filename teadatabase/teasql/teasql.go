package teasql

import (
	"sync"
)

var sqlBuildersMu sync.RWMutex

var sqlBuilders = make(map[string]sqlBuilder)

func RegisterSqlBuilder(name string, sqlBuilder sqlBuilder) {
	sqlBuildersMu.Lock()
	defer sqlBuildersMu.Unlock()
	if sqlBuilder == nil {
		panic("teasql: Register sqlBuilder is nil")
	}
	if _, dup := sqlBuilders[name]; dup {
		panic("teasql: Register called twice for sqlBuilder " + name)
	}
	sqlBuilders[name] = sqlBuilder
}

type sqlBuilder interface {
	identifierQuoter
	placeholderGeneratorGetter
}

type identifierQuoter interface {
	QuoteIdentifier(identifier string) string
}

type placeholderGeneratorGetter interface {
	GetPlaceholderGenerator() PlaceholderGenerator
}

type PlaceholderGenerator interface {
	GeneratePlaceholder() string
	Clone() PlaceholderGenerator
}
