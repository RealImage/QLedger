# QLedger
Systems that manage money do so by managing its movement - by tracking where it moved from, where it moved to, how much moved and why. QLedger is a service that provides APIs to manage the structured movement of money. 

The there are two primitives in the system: **accounts** and **transactions**. Money moves between accounts by means of a transaction. 

A **transaction** may have multiple lines - each line represents the change (*delta*) of money in one account. A valid transaction has a total delta of 0 - no money is created or destroyed, and all money moved out of any account(s) has moved in to other account(s). QLedger validates all transactions made via the API with a zero delta check.

> Phrased another way, the law of conversation of money is formalized by the rules of double entry bookkeeping - money debited from any account must be credited to another account (and vice versa), implying that all transactions must have at least two entries (double entry). QLedger makes it easy to follow these rules. 

Accounts do not need to be predefined - they are called into existence when they are first used. 

Both accounts and transactions can be tagged with an arbitrary number of key-value pairs. There is no enforcement as to the nature of the keys or values, but simple ASCII strings are recommended for easy querying and reporting. 

> When encoding information like dates into strings in the key-value pairs, they can be made sortable (and therefore queryable by range) by using placing the most significant bits first - like 2017-05-07T11:35:14Z for dates ([ISO 8601](https://en.wikipedia.org/wiki/ISO_8601)). While numbers are unlikely candidates for key-value storage, encoding them to be sortable / range-queryable without space-inefficient zero pading is [possible, but difficult](http://stackoverflow.com/questions/28413947/space-efficient-way-to-encode-numbers-as-sortable-strings).

All accounts and transactions are identified by a string identifier, which also acts an idempotency and an immutability key. Transactions once sent to the ledger cannot be changed - any 'modification' or reversal requires a new transaction. The safe recovery mechanism for all network errors is also a simple retry - as long as the identifier does not change the transaction will never be indavertently duplicated. 

#Accounts
  - id
  - data
  - tags[]
    - key
    - value
Accounts can be created on the fly 

#Transactions
  - id
  - timestamp
  - data
  - lines[]
    - account_id
    - delta
  - tags[]
    - key
    - value
   
   
To preserve double entry rules, the sum(delta) in any transaction must always = zero.

The `timestamp` and `lines[]` of a transaction are always immutable - if transaction needs to be reversed a new reversal transaction can be inserted. 

The `tags[]` of both accounts and transactions are mutable to reflect changing reporting requirements. And audit trail will be created for every overwrite. 

All calls are idempotent. The immutable properties of transactions will never be overwritten on multiple calls. Tags may be overwritten with an audit trail. 
