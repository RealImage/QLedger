package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	ledgerError "github.com/RealImage/QLedger/errors"
)

type SearchEngine struct {
	db        *sql.DB
	namespace string
}

func NewSearchEngine(db *sql.DB, namespace string) (*SearchEngine, ledgerError.ApplicationError) {
	if !(namespace == "accounts" || namespace == "transactions") {
		return nil, SearchNamespaceInvalidError(namespace)
	}
	return &SearchEngine{db: db, namespace: namespace}, nil
}

func (engine *SearchEngine) Query(q string) (interface{}, ledgerError.ApplicationError) {
	searchQuery, err := NewSearchQuery(q)
	if err != nil {
		return nil, err
	}

	sqlQuery := searchQuery.ToSQL(engine.namespace)
	log.Println("sqlQuery:", sqlQuery)
	rows, derr := engine.db.Query(sqlQuery)
	if derr != nil {
		return nil, DBError(derr)
	}
	defer rows.Close()

	switch engine.namespace {
	case "accounts":
		var accounts []*Account
		for rows.Next() {
			acc := &Account{}
			if err := rows.Scan(&acc.Id, &acc.Balance); err != nil {
				return nil, DBError(err)
			}
			accounts = append(accounts, acc)
		}
		return accounts, nil
	case "transactions":
		var transactions []*Transaction
		for rows.Next() {
			txn := &Transaction{}
			if err := rows.Scan(&txn.ID, &txn.Timestamp); err != nil {
				return nil, DBError(err)
			}
			transactions = append(transactions, txn)
		}
		return transactions, nil
	}
	return nil, nil
}

type SearchQuery struct {
	Query struct {
		Terms      []map[string]interface{} `json:"terms"`
		RangeItems []map[string]interface{} `json:"range"`
	} `json:"query"`
}

func NewSearchQuery(q string) (*SearchQuery, ledgerError.ApplicationError) {
	var searchQuery *SearchQuery
	log.Println("q:", q)
	err := json.Unmarshal([]byte(q), &searchQuery)
	if err != nil {
		return nil, SearchQueryInvalidError(err)
	}
	return searchQuery, nil
}

func (searchQuery *SearchQuery) ToSQL(namespace string) string {
	var sqlQuery string
	switch namespace {
	case "accounts":
		sqlQuery = "SELECT id, balance FROM accounts"
	case "transactions":
		sqlQuery = "SELECT id, timestamp FROM transactions"
	default:
		return ""
	}

	jsonify := func(input interface{}) string {
		j, _ := json.Marshal(input)
		return string(j)
	}
	sqlComparisonOp := func(op string) string {
		switch op {
		case "gt":
			return ">"
		case "lt":
			return "<"
		case "gte":
			return ">="
		case "lte":
			return "<="
		}
		return "="
	}

	if len(searchQuery.Query.Terms) == 0 && len(searchQuery.Query.RangeItems) == 0 {
		return sqlQuery
	}
	var where []string

	// Term queries
	/*
		-- string value
		SELECT id FROM transactions WHERE data->'status' @> '"completed"'::jsonb;
		-- boolean value
		SELECT id FROM transactions WHERE data->'active' @> 'true'::jsonb;
		-- numeric value
		SELECT id FROM transactions WHERE data->'charge' @> '2000'::jsonb;
		-- array value
		SELECT id FROM transactions WHERE data->'colors' @> '["red", "green"]'::jsonb;
		-- object value
		SELECT id FROM transactions WHERE data->'products' @> '{"qw":{"coupons": ["x001"]}}'::jsonb;
	*/
	for _, term := range searchQuery.Query.Terms {
		var conditions []string
		for key, value := range term {
			conditions = append(
				conditions,
				fmt.Sprintf("data->'%s' @> '%s'", key, jsonify(value)))
		}
		where = append(where, strings.Join(conditions, " AND "))
	}
	// Range queries
	/*
		-- numeric value
		SELECT id, data->'charge' FROM transactions WHERE data->>'charge' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'charge')::float >= 2000;
		-- other values
		SELECT id, data->'date' FROM transactions WHERE data->>'date' >= '2017-01-01' AND data->>'date' < '2017-06-01';
	*/
	for _, rangeItem := range searchQuery.Query.RangeItems {
		var conditions []string
		for key, comparison := range rangeItem {
			compItem, _ := comparison.(map[string]interface{})
			for op, value := range compItem {
				var condn string
				switch value.(type) {
				case int, int8, int16, int32, int64, float32, float64:
					condn = fmt.Sprintf(
						"data->>'%s' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'%s')::float %s %v",
						key, key, sqlComparisonOp(op), value)
				default:
					condn = fmt.Sprintf("data->>'%s' %s '%s'", key, sqlComparisonOp(op), value)
				}
				conditions = append(conditions, condn)
			}
		}
		where = append(where, strings.Join(conditions, " AND "))
	}
	sqlQuery = sqlQuery + " WHERE " + strings.Join(where, " OR ")

	return sqlQuery
}
