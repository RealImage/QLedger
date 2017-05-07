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

#### Roadmap

#### Metadata for Querying and Reports
Both accounts and transactions can be tagged with an arbitrary number of key-value pairs. There is no enforcement as to the nature of the keys or values, but simple ASCII strings are recommended for easy querying and reporting. 

> When encoding information like dates into strings in the key-value pairs, they can be made sortable (and therefore queryable by range) by using placing the most significant bits first - like 2017-05-07T11:35:14Z for dates ([ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)). While numbers are unlikely candidates for key-value storage, encoding them to be sortable / range-queryable without space-inefficient zero pading is [possible, but difficult](http://stackoverflow.com/questions/28413947/space-efficient-way-to-encode-numbers-as-sortable-strings).
