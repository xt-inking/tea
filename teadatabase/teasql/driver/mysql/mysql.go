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

func (SqlBuilder) GeneratePlaceholder() func() string {
	return func() string {
		return "?"
	}
}
