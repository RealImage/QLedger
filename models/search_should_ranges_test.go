package models

import "github.com/stretchr/testify/assert"

func (ss *SearchSuite) TestSearchAccountsWithShouldRanges() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "accounts")

	query := `{
        "query": {
            "should": {
                "ranges": [
                    {"created": {"gte": "2017-01-01", "lte": "2017-06-30"}},
                    {"created": {"gte": "2017-07-01", "lte": "2017-12-30"}}
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
                "ranges": [
                    {"created": {"gte": "2017-07-01", "lte": "2017-12-30"}},
                    {"created": {"gte": "2018-01-01", "lte": "2018-06-30"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	accounts, _ = results.([]*AccountResult)
	assert.Equal(t, 0, len(accounts), "No account should exist for given query")
}

func (ss *SearchSuite) TestSearchTransactionsWithShouldRanges() {
	t := ss.T()
	engine, _ := NewSearchEngine(ss.db, "transactions")

	query := `{
        "query": {
            "should": {
                "ranges": [
                    {"expiry": {"gte": "2018-01-01", "lte": "2018-01-30"}},
                    {"expiry": {"gte": "2018-06-01", "lte": "2018-06-30"}}
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
                "ranges": [
                    {"expiry": {"gte": "2018-06-01", "lte": "2018-06-30"}},
                    {"expiry": {"gte": "2018-07-01", "lte": "2018-07-30"}}
                ]
            }
        }
    }`
	results, err = engine.Query(query)
	assert.Equal(t, nil, err, "Error in building search query")
	transactions, _ = results.([]*TransactionResult)
	assert.Equal(t, 0, len(transactions), "No transaction should exist for given query")
}
