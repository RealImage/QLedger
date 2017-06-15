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

#### GET `/v1/accounts?id=alice`
```
{
  "id": "alice",
  "balance": -100
}
```


#### Metadata for Querying and Reports

Transactions and accounts can have arbitrary number of key-value pairs maintained as a single JSON `data` which helps in grouping and filtering them by one or more criteria.

For transactions, the `data` can be available either while creation or can later be added. But for the accounts it can be updated with `data` only after it is created through a transaction.

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
  "data": {
    "christmas-offer": "",
    "status": "completed",
    "products": {
      "qw": {
          "coupons": ["x001"]
      }
    },
    "months": ["jan", "feb"],
    "date": "2017-01-01"
  }
}
```

The transactions or accounts can be updated with `data` using endpoint `POST /v1/:model/:itemID/data`.

The transaction with ID `abcd1234` is updated with `data` as follows:
#### POST `/v1/transactions/abcd1234/data`
```
{
  "data": {
    "hold-on": "",
    "active": true,
    "products": {
      "qw": {
          "coupons": ["y001", "z001"]
      }
    },
    "months": ["jan", "mar"],
    "charge": 2000
  }
}
```

So after the above initial creation and update, the `data` of transaction `abcd1234` will look as follows:

```
{
  "data": {
    "christmas-offer": "",
    "hold-on": "",
    "status": "completed",
    "active": true,
    "products": {
      "qw": {
          "coupons": ["x001", "y001", "z001"]
      }
    },
    "months": ["jan", "feb", "mar"],
    "date": "2017-01-01",
    "charge": 2000
  }
}
```

The transactions and accounts can be filtered from the endpoint `GET /v1/search/:namespace` using the following query primitives.

- `terms` query
Find rows where atleast one of the `terms` with the specified key-value pairs exists.

- `range` query
Find rows where atleast one of the item with the specified `range` exists.

> The `/v1/search/:namespace` follows a subset of [Elasticsearch querying](https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html) format.

Here are samples how transactions can be filtered.

#### POST `/v1/search/transactions`
```
{
  "query": {
    "terms": [
      {"christmas-offer": ""},
      {"status": "completed", "active": true},
      {"months": ["jan", "feb", "mar"]},
      {
        "products": {
          "qw": {
              "coupons": "x001"
          }
        },
        "charge": 2000
      },
      {
        "products": {
          "qw": {
              "coupons": "y001"
          }
        },
        "charge": 2000
      }
    ],
    "range": [
      {"date": {"gte": "2017-01-01","lte": "2017-06-31"}},
      {
        "charge": {"gt": 2000},
        "year": {"gt": 2014}
      }
    ]
  }
}
```
