# QLedger
Systems that manage money do so by managing it's movement - by tracking where it moved from, where it moved to, how much moved and why. QLedger is a service that provides APIs to manage the structured movement of money. 

The there are two primitives in the system: **accounts** and **transactions**. 

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
