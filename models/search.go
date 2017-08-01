package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	ledgerError "github.com/RealImage/QLedger/errors"
)

type SearchEngine struct {
	db        *sql.DB
	namespace string
}

type TransactionResult struct {
	ID        string                   `json:"id"`
	Timestamp string                   `json:"timestamp"`
	Data      json.RawMessage          `json:"data"`
	Lines     []*TransactionLineResult `json:"lines"`
}

type TransactionLineResult struct {
	AccountID string `json:"account"`
	Delta     int    `json:"delta"`
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
			var rawAccounts, rawDelta string
			if err := rows.Scan(&txn.ID, &txn.Timestamp, &txn.Data, &rawAccounts, &rawDelta); err != nil {
				return nil, DBError(err)
			}
			var accounts []string
			var delta []int
			json.Unmarshal([]byte(rawAccounts), &accounts)
			json.Unmarshal([]byte(rawDelta), &delta)
			var lines []*TransactionLineResult
			for i, acc := range accounts {
				l := &TransactionLineResult{}
				l.AccountID = acc
				l.Delta = delta[i]
				lines = append(lines, l)
			}
			txn.Lines = lines
			transactions = append(transactions, txn)
		}
		return transactions, nil
	}
	return nil, nil
}

type QueryContainer struct {
	Fields     []map[string]map[string]interface{} `json:"fields"`
	Terms      []map[string]interface{}            `json:"terms"`
	RangeItems []map[string]map[string]interface{} `json:"ranges"`
}

type SearchRawQuery struct {
	Offset *int `json:"from,omitempty"`
	Limit  *int `json:"size,omitempty"`
	Query  struct {
		MustClause   QueryContainer `json:"must"`
		ShouldClause QueryContainer `json:"should"`
	} `json:"query"`
}

type SearchSQLQuery struct {
	sql  string
	args []interface{}
}

func hasValidKeys(items interface{}) bool {
	var validKey = regexp.MustCompile(`^[a-z_A-Z]+$`)
	switch t := items.(type) {
	case []map[string]interface{}:
		for _, item := range t {
			for key := range item {
				if !validKey.MatchString(key) {
					return false
				}
			}
		}
		return true
	case []map[string]map[string]interface{}:
		for _, item := range t {
			for key := range item {
				if !validKey.MatchString(key) {
					return false
				}
			}
		}
		return true
	default:
		return false
	}
}

func NewSearchRawQuery(q string) (*SearchRawQuery, ledgerError.ApplicationError) {
	var rawQuery *SearchRawQuery
	err := json.Unmarshal([]byte(q), &rawQuery)
	if err != nil {
		return nil, SearchQueryInvalidError(err)
	}
	if !(hasValidKeys(rawQuery.Query.MustClause.Fields) &&
		hasValidKeys(rawQuery.Query.MustClause.Terms) &&
		hasValidKeys(rawQuery.Query.MustClause.RangeItems) &&
		hasValidKeys(rawQuery.Query.ShouldClause.Fields) &&
		hasValidKeys(rawQuery.Query.MustClause.Terms) &&
		hasValidKeys(rawQuery.Query.MustClause.RangeItems)) {
		return nil, SearchQueryInvalidError(errors.New("Invalid key(s) in search query"))
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
		sql = `SELECT id, timestamp, data,
					array_to_json(ARRAY(
						SELECT lines.account_id FROM lines
							WHERE transaction_id=transactions.id
							ORDER BY lines.account_id
					)) AS account_array,
					array_to_json(ARRAY(
						SELECT lines.delta FROM lines
							WHERE transaction_id=transactions.id
							ORDER BY lines.account_id
					)) AS delta_array
			FROM transactions`
	default:
		return nil
	}

	// Process must queries
	var mustWhere []string
	mustClause := rawQuery.Query.MustClause
	fieldsWhere, fieldsArgs := convertFieldsToSQL(mustClause.Fields)
	mustWhere = append(mustWhere, fieldsWhere...)
	args = append(args, fieldsArgs...)
	termsWhere, termsArgs := convertTermsToSQL(mustClause.Terms)
	mustWhere = append(mustWhere, termsWhere...)
	args = append(args, termsArgs...)
	rangesWhere, rangesArgs := convertRangesToSQL(mustClause.RangeItems)
	mustWhere = append(mustWhere, rangesWhere...)
	args = append(args, rangesArgs...)

	// Process should queries
	var shouldWhere []string
	shouldClause := rawQuery.Query.ShouldClause
	fieldsWhere, fieldsArgs = convertFieldsToSQL(shouldClause.Fields)
	shouldWhere = append(shouldWhere, fieldsWhere...)
	args = append(args, fieldsArgs...)
	termsWhere, termsArgs = convertTermsToSQL(shouldClause.Terms)
	shouldWhere = append(shouldWhere, termsWhere...)
	args = append(args, termsArgs...)
	rangesWhere, rangesArgs = convertRangesToSQL(shouldClause.RangeItems)
	shouldWhere = append(shouldWhere, rangesWhere...)
	args = append(args, rangesArgs...)
	var offset = rawQuery.Offset
	var limit = rawQuery.Limit

	if len(mustWhere) == 0 && len(shouldWhere) == 0 {
		return &SearchSQLQuery{sql: sql, args: args}
	}

	sql = sql + " WHERE "
	if len(mustWhere) != 0 {
		sql = sql + strings.Join(mustWhere, " AND ")
		if len(shouldWhere) != 0 {
			sql = sql + " AND "
		}
	}
	if len(shouldWhere) != 0 {
		sql = sql + strings.Join(shouldWhere, " OR ")
	}
	if offset != nil {
		fmt.Println("offset", offset)
		sql = sql + " offset " + strconv.Itoa(*offset) + " "
	}

	if limit != nil {
		sql = sql + " limit " + strconv.Itoa(*limit)
	}

	sql = enumerateSQLPlacholder(sql)
	return &SearchSQLQuery{sql: sql, args: args}
}
