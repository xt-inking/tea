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

func (SqlBuilder) GeneratePlaceholder() func() string {
	placeholder := atomic.Uint32{}
	return func() string {
		placeholder := placeholder.Add(1)
		return "$" + strconv.FormatUint(uint64(placeholder), 10)
	}
}
