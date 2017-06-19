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

type TransactionResult struct {
	ID        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type AccountResult struct {
	ID      string          `json:"id"`
	Balance int             `json:"balance"`
	Data    json.RawMessage `json:"data"`
}

func NewSearchEngine(db *sql.DB, namespace string) (*SearchEngine, ledgerError.ApplicationError) {
	if !(namespace == "accounts" || namespace == "transactions") {
		return nil, SearchNamespaceInvalidError(namespace)
	}
	return &SearchEngine{db: db, namespace: namespace}, nil
}

func (engine *SearchEngine) Query(q string) (interface{}, ledgerError.ApplicationError) {
	rawQuery, err := NewSearchRawQuery(q)
	if err != nil {
		return nil, err
	}

	sqlQuery := rawQuery.ToSQLQuery(engine.namespace)
	log.Println("sqlQuery SQL:", sqlQuery.sql)
	log.Println("sqlQuery args:", sqlQuery.args)
	rows, derr := engine.db.Query(sqlQuery.sql, sqlQuery.args...)
	if derr != nil {
		return nil, DBError(derr)
	}
	defer rows.Close()

	switch engine.namespace {
	case "accounts":
		accounts := make([]*AccountResult, 0)
		for rows.Next() {
			acc := &AccountResult{}
			if err := rows.Scan(&acc.ID, &acc.Balance, &acc.Data); err != nil {
				return nil, DBError(err)
			}
			accounts = append(accounts, acc)
		}
		return accounts, nil
	case "transactions":
		transactions := make([]*TransactionResult, 0)
		for rows.Next() {
			txn := &TransactionResult{}
			if err := rows.Scan(&txn.ID, &txn.Timestamp, &txn.Data); err != nil {
				return nil, DBError(err)
			}
			transactions = append(transactions, txn)
		}
		return transactions, nil
	}
	return nil, nil
}

type SearchRawQuery struct {
	Query struct {
		ID         string                   `json:"id"`
		Terms      []map[string]interface{} `json:"terms"`
		RangeItems []map[string]interface{} `json:"range"`
	} `json:"query"`
}

type SearchSQLQuery struct {
	sql  string
	args []interface{}
}

func NewSearchRawQuery(q string) (*SearchRawQuery, ledgerError.ApplicationError) {
	var rawQuery *SearchRawQuery
	err := json.Unmarshal([]byte(q), &rawQuery)
	if err != nil {
		return nil, SearchQueryInvalidError(err)
	}
	return rawQuery, nil
}

func (rawQuery *SearchRawQuery) ToSQLQuery(namespace string) *SearchSQLQuery {
	var sql string
	var args []interface{}

	switch namespace {
	case "accounts":
		sql = "SELECT id, balance, data FROM current_balances"
	case "transactions":
		sql = "SELECT id, timestamp, data FROM transactions"
	default:
		return nil
	}
	if len(rawQuery.Query.ID) != 0 {
		sql = sql + " WHERE id = $1"
		return &SearchSQLQuery{sql: sql, args: []interface{}{rawQuery.Query.ID}}
	}
	if len(rawQuery.Query.Terms) == 0 && len(rawQuery.Query.RangeItems) == 0 {
		return &SearchSQLQuery{sql: sql}
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
	for _, term := range rawQuery.Query.Terms {
		var conditions []string
		for key, value := range term {
			conditions = append(
				conditions,
				fmt.Sprintf("data->'%s' @> $%d::jsonb", key, len(args)+1),
			)
			args = append(args, jsonify(value))
		}
		where = append(where, "("+strings.Join(conditions, " AND ")+")")
	}
	// Range queries
	/*
		-- numeric value
		SELECT id, data->'charge' FROM transactions WHERE data->>'charge' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'charge')::float >= 2000;
		-- other values
		SELECT id, data->'date' FROM transactions WHERE data->>'date' >= '2017-01-01' AND data->>'date' < '2017-06-01';
	*/
	for _, rangeItem := range rawQuery.Query.RangeItems {
		var conditions []string
		for key, comparison := range rangeItem {
			compItem, _ := comparison.(map[string]interface{})
			for op, value := range compItem {
				var condn string
				switch value.(type) {
				case int, int8, int16, int32, int64, float32, float64:
					condn = fmt.Sprintf(
						"data->>'%s' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'%s')::float %s $%d",
						key, key, sqlComparisonOp(op), len(args)+1,
					)
				default:
					condn = fmt.Sprintf("data->>'%s' %s $%d", key, sqlComparisonOp(op), len(args)+1)
				}
				conditions = append(conditions, condn)
				args = append(args, value)
			}
		}
		where = append(where, "("+strings.Join(conditions, " AND ")+")")
	}
	sql = sql + " WHERE " + strings.Join(where, " OR ")

	return &SearchSQLQuery{sql: sql, args: args}
}
