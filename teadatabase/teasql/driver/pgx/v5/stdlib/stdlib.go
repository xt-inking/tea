package stdlib

import (
	"strconv"
	"sync/atomic"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/tea-frame-go/tea/teadatabase/teasql"
)

func init() {
	teasql.RegisterSqlBuilder("pgx", sqlBuilder{})
}

type sqlBuilder struct{}

func (sqlBuilder) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (sqlBuilder) GetPlaceholderGenerator() teasql.PlaceholderGenerator {
	return &placeholderGenerator{
		placeholder: atomic.Uint32{},
	}
}

type placeholderGenerator struct {
	placeholder atomic.Uint32
}

func (pg *placeholderGenerator) GeneratePlaceholder() string {
	placeholder := pg.placeholder.Add(1)
	return "$" + strconv.FormatUint(uint64(placeholder), 10)
}

func (pg *placeholderGenerator) Clone() teasql.PlaceholderGenerator {
	newPlaceholderGenerator := &placeholderGenerator{
		placeholder: atomic.Uint32{},
	}
	newPlaceholderGenerator.placeholder.Store(pg.placeholder.Load())
	return newPlaceholderGenerator
}
