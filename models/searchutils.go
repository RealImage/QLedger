package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// replaces $N to corresponding placeholder index ($1, $2,...)
func enumerateSQLPlacholder(msql string) (psql string) {
	splitItems := strings.Split(msql, "$N")
	for i, item := range splitItems {
		if i != len(splitItems)-1 {
			psql = psql + item + fmt.Sprintf("$%d", i+1)
		} else {
			psql = psql + item
		}
	}
	return
}

func jsonify(input interface{}) string {
	j, _ := json.Marshal(input)
	return string(j)
}

func sqlComparisonOp(op string) string {
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

func convertTermsToSQL(terms []map[string]interface{}) (where []string, args []interface{}) {
	// Sample terms
	/*
	   "terms": [
	       {"status": "completed", "active": true},
	       {"charge": 2000},
	       {"colors": ["red", "green"]},
	       {"products":{"qw":{"coupons":["x001"]}}}
	   ]
	*/
	// Corresponding SQL
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
	for _, term := range terms {
		var conditions []string
		for key, value := range term {
			conditions = append(
				conditions,
				fmt.Sprintf("data->'%s' @> $N::jsonb", key),
			)
			args = append(args, jsonify(value))
		}
		where = append(where, "("+strings.Join(conditions, " AND ")+")")
	}
	return
}

func convertRangesToSQL(ranges []map[string]map[string]interface{}) (where []string, args []interface{}) {
	// Sample ranges
	/*
	   "ranges": [
	       {"charge": {"gte": 2000, "lte": 4000}},
	       {"date": {"gt": "2017-01-01","lt": "2017-06-31"}}
	   ]
	*/
	// Corresponding SQL
	/*
	   -- numeric value
	   SELECT id, data->'charge' FROM transactions WHERE data->>'charge' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'charge')::float >= 2000;
	   -- other values
	   SELECT id, data->'date' FROM transactions WHERE data->>'date' >= '2017-01-01' AND data->>'date' < '2017-06-31';
	*/
	for _, rangeItem := range ranges {
		var conditions []string
		for key, comparison := range rangeItem {
			for op, value := range comparison {
				var condn string
				switch value.(type) {
				case int, int8, int16, int32, int64, float32, float64:
					condn = fmt.Sprintf(
						"data->>'%s' ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND (data->>'%s')::float %s $N",
						key, key, sqlComparisonOp(op),
					)
				default:
					condn = fmt.Sprintf("data->>'%s' %s $N", key, sqlComparisonOp(op))
				}
				conditions = append(conditions, condn)
				args = append(args, value)
			}
		}
		where = append(where, "("+strings.Join(conditions, " AND ")+")")
	}
	return
}

func convertFieldsToSQL(fields []map[string]map[string]interface{}) (where []string, args []interface{}) {
	// Sample ranges
	/*
	   "fields": [
	       {"id": {"eq": "ACME.CREDIT"}, "balance": {"lt": 0}},
	   ]
	*/
	// Corresponding SQL
	/*
	   -- numeric value
	   SELECT id, balance, data FROM accounts WHERE id = 'ACME.CREDIT' AND balance ~ '^([0-9]+[.]?[0-9]*|[.][0-9]+)$' AND balance::float >= 0;
	*/
	for _, field := range fields {
		var conditions []string
		for key, comparison := range field {
			for op, value := range comparison {
				condn := fmt.Sprintf("%s %s $N", key, sqlComparisonOp(op))
				conditions = append(conditions, condn)
				args = append(args, value)
			}
		}
		where = append(where, "("+strings.Join(conditions, " AND ")+")")
	}
	return
}
