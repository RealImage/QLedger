package models

import "github.com/stretchr/testify/assert"

func (ss *SearchSuite) TestSearchAccountsWithMustTerms() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "must": {
                "terms": [
                    {"customer_id": "C1"},
                    {"status": "active"}
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
                "terms": [
                    {"customer_id": "C2"},
                    {"status": "active"}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 0, len(accounts), "No account should exist for given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithMustTerms() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "must": {
                "terms": [
                    {"action": "setcredit"},
                    {"months": ["jan", "feb", "mar"]}
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
                "terms": [
                    {"action": "setcredit"},
                    {"months": ["oct", "nov", "dec"]}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 0, len(transactions), "No transaction should exist for given query")
}
