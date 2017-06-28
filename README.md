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

A typical `data` object will be as follows:
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
> The key value formats here are just samples and they can be any valid JSON object.

The transactions can be created with `data` as follows:
##### POST `/v1/transactions`
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
##### POST `/v1/accounts`
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
##### PUT `/v1/transactions`
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

##### GET `/v1/transactions`
The transactions and accounts can be filtered from the endpoints `GET /v1/transactions` and `GET /v1/accounts` with the following query primitives in the payload.

**- `id` query**

Find row which exactly matches the specified `id`

Example: The following query matches a single transaction with ID `txn1`
`GET /v1/transactions`
```
{
  "query": {
    "id": "txn1"
  }
}
```

**- `terms` query**

Find rows where atleast one of the `terms` with all the specified key-value pairs exists.

Example: The following query matches all transactions which satisfies atleast one of the following terms:
- `status` is `completed` AND `active` is `true`
-  Values `jan`, `feb` AND `mar` in `months` array
-  Subset `{"qw": {"tax": 18.0}}` in `products` object

`GET /v1/transactions`
```
{
  "query": {
    "terms": [
      {"status": "completed", "active": true},
      {"months": ["jan", "feb", "mar"]},
      {
        "products": {
          "qw": {
              "tax": 18.0
          }
        }
      }
    ]
  }
}
```

**- `range` query**

Find rows where atleast one of the `range`  condition is satisfied.

Example: The following query matches all transactions which satisfies atleast one of the following `range` condition:

- `charge >= 2000` AND `charge <= 4000`
- `date > '2017-01-01'` AND `date < '2017-01-31'`

`GET /v1/transactions`
```
{
  "query": {
    "range": [
      {"charge": {"gte": 2000, "lte": 4000}},
      {"date": {"gt": "2017-01-01","lt": "2017-06-31"}}
    ]
  }
}
```

**Note:**

- The `GET /v1/transactions` and `GET /v1/accounts` follows a subset of [Elasticsearch querying](https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html) format.

-  Clients those doesn't support passing payload in the `GET` can use the `POST` alternate endpoints: `POST /v1/transactions/_search` and `POST /v1/accounts/_search`
