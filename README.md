# ledger
A general ledger, suitable for managing any system with financial events

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
