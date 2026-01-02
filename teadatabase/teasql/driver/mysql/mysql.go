package mysql

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/tea-frame-go/tea/teadatabase/teasql"
)

func init() {
	teasql.RegisterSqlBuilder("mysql", SqlBuilder{})
}

type SqlBuilder struct{}

func (SqlBuilder) QuoteIdentifier(identifier string) string {
	return "`" + identifier + "`"
}

func (SqlBuilder) GetPlaceholderGenerator() teasql.PlaceholderGenerator {
	return PlaceholderGenerator{}
}

type PlaceholderGenerator struct{}

func (PlaceholderGenerator) GeneratePlaceholder() string {
	return "?"
}

func (PlaceholderGenerator) Clone() teasql.PlaceholderGenerator {
	return PlaceholderGenerator{}
}
