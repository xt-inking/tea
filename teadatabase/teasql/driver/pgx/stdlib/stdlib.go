package stdlib

import (
	"strconv"
	"sync/atomic"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/tea-frame-go/tea/teadatabase/teasql"
)

func init() {
	teasql.RegisterSqlBuilder("pgx", SqlBuilder{})
}

type SqlBuilder struct{}

func (SqlBuilder) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (SqlBuilder) GetPlaceholderGenerator() teasql.PlaceholderGenerator {
	return &PlaceholderGenerator{
		placeholder: atomic.Uint32{},
	}
}

type PlaceholderGenerator struct {
	placeholder atomic.Uint32
}

func (placeholderGenerator *PlaceholderGenerator) GeneratePlaceholder() string {
	placeholder := placeholderGenerator.placeholder.Add(1)
	return "$" + strconv.FormatUint(uint64(placeholder), 10)
}

func (placeholderGenerator *PlaceholderGenerator) Clone() teasql.PlaceholderGenerator {
	newPlaceholderGenerator := &PlaceholderGenerator{
		placeholder: atomic.Uint32{},
	}
	newPlaceholderGenerator.placeholder.Store(placeholderGenerator.placeholder.Load())
	return newPlaceholderGenerator
}
