package models

import "github.com/stretchr/testify/assert"

func (ss *SearchSuite) TestSearchAccountsWithMustFields() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "acc1"}},
                    {"balance": {"gt": 0}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ := results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "acc2"}},
                    {"balance": {"gt": 0}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 0, len(accounts), "No account should exist for given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithMustFields() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "txn1"}},
                    {"timestamp": {"gte": "2017-08-08"}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ := results.([]*TransactionResult)
	assert.Equal(t, 1, len(transactions), "Transactions count doesn't match")
	assert.Equal(t, "txn1", transactions[0].ID, "Transaction ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "txn2"}},
                    {"timestamp": {"lt": "2017-08-08"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 0, len(transactions), "No transaction should exist for given query")
}

func (ss *SearchSuite) TestSearchAccountsWithFieldOperators() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"eq": "acc1"}}
                ]
            }
        }
    }`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ := results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"lt": "acc2"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"lte": "acc2"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 2, len(accounts), "Accounts count doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"gt": "acc1"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc2", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"gte": "acc1"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 2, len(accounts), "Accounts count doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"ne": "acc1"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc2", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"like": "%c1"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc1", accounts[0].ID, "Account ID doesn't match")

	query = `{
        "query": {
            "must": {
                "fields": [
                    {"id": {"notlike": "%c1"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
	assert.Equal(t, "acc2", accounts[0].ID, "Account ID doesn't match")
}
