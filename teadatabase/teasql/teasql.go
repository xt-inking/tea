package teasql

var IdentifierQuoter identifierQuoter

type identifierQuoter interface {
	QuoteIdentifier(Identifier) string
}

type Identifier string

func (identifier Identifier) Quote() string {
	return IdentifierQuoter.QuoteIdentifier(identifier)
}
