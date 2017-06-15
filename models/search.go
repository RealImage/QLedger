package models

import (
	"errors"

	ledgerError "github.com/RealImage/QLedger/errors"
)

type SearchEngine struct {
	namespace string
}

func NewSearchEngine(namespace string) (*SearchEngine, ledgerError.ApplicationError) {
	if !(namespace == "accounts" || namespace == "transactions") {
		return nil, SearchNamespaceInvalidError(namespace)
	}
	return &SearchEngine{namespace: namespace}, nil
}

func (engine *SearchEngine) Query(q string) (interface{}, ledgerError.ApplicationError) {
	//TODO: Implement querying
	return nil, SearchQueryInvalidError(errors.New("not implemented"))
}
