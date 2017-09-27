[![CircleCI](https://circleci.com/gh/RealImage/QLedger.svg?style=svg)](https://circleci.com/gh/RealImage/QLedger)

# QLedger
Systems that manage money do so by managing its movement - by tracking where it moved from, where it moved to, how much moved and why. QLedger is a service that provides APIs to manage the structured movement of money.

The there are two primitives in the system: **accounts** and **transactions**. Money moves between accounts by means of a transaction.

A **transaction** may have multiple *lines* - each line represents the change (*delta*) of money in one account. A valid transaction has a total delta of zero - no money is created or destroyed, and all money moved out of any account(s) has moved in to other account(s). QLedger validates all transactions made via the API with a zero delta check.

> Phrased another way, the law of conversation of money is formalized by the rules of double entry bookkeeping - money debited from any account must be credited to another account (and vice versa), implying that all transactions must have at least two entries (double entry) with a zero sum delta. QLedger makes it easy to follow these rules.

Accounts do not need to be predefined - they are called into existence when they are first used.

All accounts and transactions are identified by a string identifier, which also acts an idempotency and an immutability key. Transactions once sent to the ledger cannot be changed - any 'modification' or reversal requires a new transaction. The safe recovery mechanism for all network errors is also a simple retry - as long as the identifier does not change the transaction will never be indavertently duplicated.

#### POST `/v1/transactions`
```
{
  "id": "abcd1234",
  "lines": [
    {
      "account": "alice",
      "delta": -100
    },
    {
      "account": "bob",
      "delta": 100
    }
  ],
  ...
}
```
> Transactions with a total delta not equal to zero will result in a `400 BAD REQUEST` error.


#### Metadata for Querying and Reports

Transactions and accounts can have arbitrary number of key-value pairs maintained as a single JSON `data` which helps in grouping and filtering them by one or more criteria.

Both transactions and accounts can be updated multiple times with `data`. The existing `data` is always overwritten with the new `data` value.

The `data` can be arbitrary JSON value as follows:
```
{
  "data": {
    "k1": "",
    "k2": "strval",
    "k3": ["av1", "av2", "av3"],
    "k4": {
      "nest1": {
        "nest2": "val"
      }
    },
    "k5": 2017,
    "k6": "2017-12-01"
  }
}
```

The transactions can be created with `data` as follows:

`POST /v1/transactions`
```
{
  "id": "abcd1234",
  "lines": [
    {
      "account": "alice",
      "delta": -100
    },
    {
      "account": "bob",
      "delta": 100
    }
  ],
  "data": {
    "christmas-offer": "",
    "status": "completed",
    "products": {
      "qw": {
          "tax": 14.5
      }
    },
    "months": ["jan", "feb"],
    "date": "2017-01-01"
  }
}
```

The accounts can be created with `data` as follows:

`POST /v1/accounts`
```
{
  "id": "alice",
  "data": {
    "product": "qw",
    "date": "2017-01-01"
  }
}
```

The transactions or accounts can be updated with `data` using endpoints `PUT /v1/transactions` and `PUT /v1/accounts`

The transaction with ID `abcd1234` is updated with `data` as follows:

`PUT /v1/transactions`
```
{
  "id": "abcd1234",
  "data": {
    "christmas-offer": "",
    "hold-on": "",
    "status": "completed",
    "active": true,
    "products": {
      "qw": {
          "tax": 18.0
      }
    },
    "months": ["jan", "feb", "mar"],
    "date": "2017-01-01",
    "charge": 2000
  }
}
```

#### Searching of accounts and transactions

The transactions and accounts can be filtered from the endpoints `GET /v1/transactions` and `GET /v1/accounts` with the search query formed using the bool clauses(`must` and `should`) and query types(`fields`, `terms` and `ranges`).

##### Query types:

###### - `fields` query

Find items where the specified column exists with the specified value in the specified range.

Example fields:
- Field `{"id": {"eq": "ACME.CREDIT"}}` filters items where the column `id` is equal to `ACME.CREDIT`
- Field `{"balance": {"lt": 0}}` filters items where the column `balance` is less than `0`
- Field `{"timestamp": {"gte": "2017-01-01T05:30"}}` filters items where `timestamp` is greater than or equal to `2017-01-01T05:30`
- Field `{"id": {"ne": "ACME.CREDIT"}}` filters items where the column `id` is not equal to `ACME.CREDIT`
- Field `{"id": {"like": "%.DEBIT"}}` filters items where the column `id` ends with `.DEBIT`
- Field `{"id": {"notlike": "%.DEBIT"}}` filters items where the column `id` doesn't ends with `.DEBIT`

> The supported range operators are `lt`(less than), `lte`(less than or equal), `gt`(greater than), `gte`(greater than or equal), `eq`(equal), `ne`(not equal), `like`(like patterns), `notlike`(not like patterns).

###### - `terms` query

Filters items where the specified key-value pairs in a term exists in the `data` JSON.

Example terms:
- Term `{"status": "completed", "active": true}` filters items where `data.status` is `completed` AND `data.active` is `true`
- Term `{"months": ["jan", "feb", "mar"]}` filters items where values `jan`, `feb` AND `mar` in `data.months` array
- Term `{"products":{"qw":{"tax":18.0}}}` filters items where subset `{"qw": {"tax": 18.0}}` in `products` object


###### - `range` query

Filters items which the specified key in `data` JSON exists in the specified range of values.

Example range:
- Range `{"charge": {"gte": 2000, "lte": 4000}}` filters items where `data.charge >= 2000` AND `data.charge <= 4000`
- Range `{"date": {"gt": "2017-01-01","lt": "2017-06-31"}}` filters items where `data.date > '2017-01-01'` AND `data.date < '2017-01-31'`


##### Bool clauses:
The following bool clauses determine whether all or any of the queries needs to be satisfied.

###### - `must` clause
All of the query items in the `must` clause must be satisfied to get results.

> The `must` clause can be equated with boolean `AND`

Example: The following query matches requests to match accounts which satisfies **ALL** of the following items:

- Field `balance > 0`
- Term `data.type` is `credit` AND `data.active` is `true`
- Term `data.months` with values `jan`, `feb` AND `mar`
- Range `data.coupon >= 2000` AND `data.coupon >= 4000`
- Range `data.date > '2017-01-01'` AND `data.date < '2017-06-31'`

`GET /v1/accounts`
```
{
  "query": {
      "must": {
        "fields": [
            {"balance": {"gt": 0}}
        ],
        "terms": [
            {"type": "credit", "active": true},
            {"months": ["jan", "feb", "mar"]}
        ],
        "ranges": [
            {"coupon": {"gte": 2000, "lte": 4000}},
            {"date": {"gt": "2017-01-01","lt": "2017-06-31"}}
        ]
      }
  }
}
```



###### - `should` clause
Any of the query items in the `should` clause should be satisfied to get results.

> The `should` clause can be equated with boolean `OR`
>
Example: The following query matches requests to match transactions which satisfies **ANY** of the following items:

- Field `id = '2017-06-31T05:00:45'`
- Term `data.type` is `company.credit` AND `order_id` is `001`
- Range `data.timestamp >= '2017-01-01T05:30'`

`GET /v1/accounts`
```
{
  "query": {
      "should": {
        "fields": [
            {"id": {"eq": "intent_QW_001"}}
        ],
        "terms": [
            {"type": "company.credit", "order_id": "001"}
        ],
        "ranges": [
            {"timestamp": {"gte": "2017-01-01T05:30"}}
        ]
      }
  }
}
```


**Note:**

- This search API follows a subset of [Elasticsearch querying](https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html) format.

-  Clients those doesn't support passing search payload in the `GET`, can alternatively use the `POST`  endpoints: `POST /v1/transactions/_search` and `POST /v1/accounts/_search`.

- A search query can have both `must` and `should` clauses.
