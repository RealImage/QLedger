package models

import "github.com/stretchr/testify/assert"

func (ss *SearchSuite) TestSearchAccountsWithMustRanges() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "must": {
                "ranges": [
                    {"created": {"gte": "2017-01-01"}},
                    {"created": {"lte": "2017-02-01"}}
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
                "ranges": [
                    {"created": {"gte": "2017-07-01"}},
                    {"created": {"lte": "2017-12-30"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 0, len(accounts), "No account should exist for given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithMustRanges() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "must": {
                "ranges": [
                    {"expiry": {"gte": "2018-01-01"}},
                    {"expiry": {"lte": "2018-01-02"}}
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
                "ranges": [
                    {"expiry": {"gte": "2018-02-01"}},
                    {"expiry": {"lte": "2018-02-05"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 0, len(transactions), "No transaction should exist for given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithIsOperator() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	// Test IS operator
	query := `{
		"query": {
			"must": {
				"ranges": [
					{"type": {"is": null}}
				]
			}
		}
	}`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ := results.([]*TransactionResult)
	assert.Equal(t, 3, len(transactions), "Transactions count doesn't match")

	// Test IS NOT operator
	query = `{
		"query": {
			"must": {
				"ranges": [
					{"action": {"isnot": null}}
				]
			}
		}
	}`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 3, len(transactions), "Transactions count doesn't match")
}

func (ss *SearchSuite) TestSearchAccountsWithInOperator() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	// Test IS operator
	query := `{
		"query": {
			"must": {
				"ranges": [
					{"customer_id": {"in": ["C1", "C2", "C3"]}}
				]
			}
		}
	}`
	results, err := engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ := results.([]*AccountResult)
	assert.Equal(t, 2, len(accounts), "Accounts count doesn't match")

	// Test IS NOT operator
	query = `{
		"query": {
			"must": {
				"ranges": [
					{"customer_id": {"in": ["C2", "C3", "C4"]}}
				]
			}
		}
	}`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 1, len(accounts), "Accounts count doesn't match")
}
