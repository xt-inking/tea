package mysql

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/tea-frame-go/tea/teadatabase/teasql"
)

func init() {
	teasql.RegisterSqlBuilder("mysql", sqlBuilder{})
}

type sqlBuilder struct{}

func (sqlBuilder) QuoteIdentifier(identifier string) string {
	return "`" + identifier + "`"
}

func (sqlBuilder) GetPlaceholderGenerator() teasql.PlaceholderGenerator {
	return placeholderGenerator{}
}

type placeholderGenerator struct{}

func (placeholderGenerator) GeneratePlaceholder() string {
	return "?"
}

func (placeholderGenerator) Clone() teasql.PlaceholderGenerator {
	return placeholderGenerator{}
}
