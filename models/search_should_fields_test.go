package models

import "github.com/stretchr/testify/assert"

func (ss *SearchSuite) TestSearchAccountsWithShouldFields() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "should": {
                "fields": [
                    {"id": {"eq": "acc1"}},
                    {"id": {"eq": "acc2"}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ := results.([]*AccountResult)
	assert.Equal(t, 2, len(accounts), "Accounts count doesn't match")

	query = `{
        "query": {
            "should": {
                "fields": [
                    {"id": {"eq": "acc3"}},
                    {"id": {"eq": "acc4"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 0, len(accounts), "No account should exist for the given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithShouldFields() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "should": {
                "fields": [
                    {"id": {"eq": "txn1"}},
                    {"id": {"eq": "txn2"}},
                    {"id": {"eq": "txn3"}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ := results.([]*TransactionResult)
	assert.Equal(t, 3, len(transactions), "Transactions count doesn't match")

	query = `{
        "query": {
            "should": {
                "fields": [
                    {"id": {"eq": "txn4"}},
                    {"id": {"eq": "txn5"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 0, len(transactions), "No transaction should exist for the given query")
}
