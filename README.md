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

Transactions and accounts can be tagged with arbitrary number of key-value pairs which helps in grouping and filtering by one or more criteria.

The keys are usually a plain string and the values can be one of string, numbers, boolean, arrays(of string, numbers or boolean) or an object. The object values can be composite of all the above types in recursive manner.

For transactions, the tags can be available either while creation or can later be tagged. But the accounts can be tagged only after it is created through a transaction.

Both transactions and accounts can be re-tagged multiple times. When a same key is re-tagged, the old value will be overwritten by the new value. When a new key is available while re-tagging, the key & value will be added to the existing list of tags.

A typical tags object will be as follows:
```
{
  "key1": null,
  "key2": null,
  "key3": "val1",
  "key4": ["val1, "val2", "val3"],
  "key5": {
    "key51": {
      "key511": "val1"
    },
    "key52": ["val1", "val2"]
  }
}
```

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
  "tags": {
    "christmas-offer": null,
    "status": "completed",
    "active": true,
    "months": ["jan", "feb"],
    "products": {
      "qw": {
        "coupons": ["x001"],
        "year": 2017
      }
    }
  }
}
```

The transactions or accounts can be re-tagged using API `POST /v1/:model/:itemID/tags`.

The transaction with ID `abcd1234` can be re-tagged as follows:
#### POST `/v1/transactions/abcd1234/tags`
```
{
  "hold-on": null,
  "active": true,
  "months": ["jan, "feb", "mar"],
  "products": {
    "qw": {
      "coupons": ["x001", "y001"],
      "year": 2017
    },
    "jt": {
      "coupons": ["x001", "y001", "z001"],
      "year": 2016
    }
  }
}
```

So after the above initial tagging and re-tagging, the tags list of transaction `abcd1234` will look as follows:

```
{
  "christmas-offer": null,
  "hold-on": null,
  "status": "completed",
  "active": true,
  "months": ["jan, "feb", "mar"],
  "products": {
    "qw": {
      "coupons": ["x001"],
      "year": 2017
    },
    "jt": {
      "coupons": ["x001", "y001", "z001"],
      "year": 2016
    }
  }
}
```

The transactions and accounts can be filtered using the tags filter API `GET /v1/:model/tags/:tag1/:tag2/.../:tagN`

Here are samples how transactions can be filtered with tags:

#### GET `/v1/transactions/tags/christmas-offer`
All transactions with tag `{"christmas-offer": V}`  will be listed.
> V can be of any value even `null`

#### GET `/v1/transactions/tags/christmas-offer/hold-on`
All transactions with BOTH tags `{"christmas-offer": V, "hold-on": V}` will be listed.
> V can be of any value even `null`

#### GET `/v1/transactions/tags/status:completed`
All transactions with tag `{"status": "completed"}` will be listed.

#### GET `/v1/transactions/tags/months:jan,feb`
All transactions with tag `{"months": V}` will be listed.
> V should be atleast one of the values of `"jan"`, `"feb"`, `["jan",..]`, `["feb",..]`

#### GET `/v1/transactions/tags/months:jan,feb/active:true`
All transactions with BOTH tags `{"months": V1}`  AND `{"active": V2}` will be listed.
> V1 should be atleast one of the values of `"jan"`, `"feb"`, `["jan",..]`, `["feb",..]`
> V2 can be either `true` or `"true"`

#### GET `/v1/transactions/tags/products:qw:year:2017`
All transactions with tag `{"products": {"qw": {"year": V,..},..},..}` will be listed.
> V can be either `2017` or `"2017"`

#### GET `/v1/transactions/tags/products:jt:coupons:x001,y001`
All transactions with tag `{"products": {"jt": {"coupons": V,..},..},..}` will be listed.
> V1 should be atleast one of the values of `"x001"`, `"y001"`, `["x001",..]`, `["y001",..]`
