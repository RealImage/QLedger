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

Transactions and accounts can be tagged with arbitrary number of key-value pairs which helps in grouping and filtering by one or more criteria. The keys and values are limited to be simple ASCII strings for now.

For transactions, the tags can be available either while creation or can later be tagged. But the accounts can be tagged only after it is created through a transaction.

Both transactions and accounts can be re-tagged multiple times. Duplicate key-value pairs are simply ignored.

A typical tags list will be as follows:
```
{
  tags: [
    {"key": "k1", "value": ""},             // Key only
    {"key": "k2", "value": "v2"},           // Key Value
    {"key": "k3.n1.n2", "value": "v3"},     // Nested key using (.)
    {"key": "k4:n1:n2", "value": "v4"},     // Nested key using (:)
    {"key": "k5", "value": "1234"},         // Numerical value as string
    {"key": "k6", "value": "2017-12-01"},   // Date value as string
    {"key": "k7", "value": "av1"},          // Array values "av1", "av2", "av3"
    {"key": "k7", "value": "av2"},
    {"key": "k7", "value": "av3"}
  ]
}
```
> The key value formats here are just samples and they can be any possible ASCII combination.

The transactions can be tagged while creation as follows:
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
  tags: [
    {"key": "christmas-offer", "value": ""},
    {"key": "status", "value": "completed"},
    {"key": "products.qw.coupons", "value": "x001"},
    {"key": "products.qw.year", "value": "2017"},
    {"key": "months", "value": "jan"},
    {"key": "months", "value": "feb"}
  ]
}
```

The transactions or accounts can be re-tagged using API `POST /v1/:model/:itemID/tags`.

The transaction with ID `abcd1234` can be re-tagged as follows:
#### POST `/v1/transactions/abcd1234/tags`
```
{
  tags: [
    {"key": "hold-on", "value": ""},
    {"key": "active", "value": "true"},
    {"key": "products.qw.coupons", "value": "y001"},
    {"key": "products.qw.coupons", "value": "z001"},
    {"key": "months", "value": "jan"},
    {"key": "months", "value": "mar"}
  ]
}
```

So after the above initial tagging and re-tagging, the tags list of transaction `abcd1234` will look as follows:

```
{
  tags: [
    {"key": "christmas-offer", "value": ""},
    {"key": "hold-on", "value": ""},
    {"key": "active", "value": "true"},
    {"key": "status", "value": "completed"},
    {"key": "products.qw.coupons", "value": "x001"},
    {"key": "products.qw.coupons", "value": "y001"},
    {"key": "products.qw.coupons", "value": "z001"},
    {"key": "products.qw.year", "value": "2017"},
    {"key": "months", "value": "jan"},
    {"key": "months", "value": "feb"},
    {"key": "months", "value": "mar"}
  ]
}
```

The transactions and accounts can be filtered using the tags filter API `GET /v1/:model/tags/:tag1/:tag2/.../:tagN`

Here are samples how transactions can be filtered with tags:

#### GET `/v1/transactions/tags/christmas-offer`
All transactions with tag `{"key": "christmas-offer", "value": V}`  will be listed.
> V can be of any value

#### GET `/v1/transactions/tags/christmas-offer/hold-on`
All transactions with BOTH tags `{"key": "christmas-offer", "value": V1}` and `{"key": "hold-on", "value": V2}` will be listed.
> V1 and V2 can be of any value

#### GET `/v1/transactions/tags/status:completed`
All transactions with tag `{"key": "status", "value": "completed"}` will be listed.

#### GET `/v1/transactions/tags/months:jan,feb`
All transactions with tag `{"key": "months", "value": V}` will be listed.
> V should be ATLEAST ONE of the values of `"jan"`, `"feb"`

#### GET `/v1/transactions/tags/months:jan,feb/active:true`
All transactions with BOTH tags `{"key": "months", "value": V}`  AND `{"key": "active", "value": "true"}` will be listed.
> V should be ATLEAST ONE of the values of `"jan"`, `"feb"`

#### GET `/v1/transactions/tags/products.qw.year:2017`
All transactions with tag `{"key": "products.qw.year", "value": "2017"}` will be listed.
> Warning: The nested key can also be formatted with (:), so care must be taken while implementation.

#### GET `/v1/transactions/tags/products.qw.coupons:x001,y001`
All transactions with tag `{"key": "products.qw.coupons", "value": V}` will be listed.
> V should be atleast one of the values of `"x001"`, `"y001"`
