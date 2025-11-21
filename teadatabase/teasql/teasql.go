package teasql

import (
	"sync"
)

var identifierQuotersMu sync.RWMutex

var identifierQuoters = make(map[string]identifierQuoter)

func Register(name string, identifierQuoter identifierQuoter) {
	identifierQuotersMu.Lock()
	defer identifierQuotersMu.Unlock()
	if identifierQuoter == nil {
		panic("teasql: Register identifierQuoter is nil")
	}
	if _, dup := identifierQuoters[name]; dup {
		panic("teasql: Register called twice for identifierQuoter " + name)
	}
	identifierQuoters[name] = identifierQuoter
}

type identifierQuoter interface {
	QuoteIdentifier(identifier string) string
}
