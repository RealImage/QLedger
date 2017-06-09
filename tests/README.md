## Scripts

#### CSV Tests runner

CSV tests are run to ensure that the server with specified load works correctly in all cases. Different cases such as the transactions executed sequentially, in parallel and same transactions repeated in parallel are executed and account balances are verified.

```
go run transaction_checker.go [-endpoint ENDPOINT] -filename FILENAME -load LOAD
```

where
-  `endpoint` - API endpoint
-  `filename` - Transactions CSV file (default "transactions.csv")
-  `load` - Load count for repeating the tests (default 10)


The CSV file should be of the following format.
```
transaction_id,account_id,delta
100,alice,100
100,bob,-100
101,alice,100
101,bob,-50
101,carly,-50
```

Here are sample executions:

- Run in a local server
```
go run transaction_checker.go -endpoint=http://127.0.0.1:7000 -load=10 -filename=transactions.csv
```

- Run in a remote Heroku endpoint
```
go run transaction_checker.go -endpoint=https://qubeledger.herokuapp.com -load=10 -filename=transactions.csv
```

- Run in a local test server created on-demand:
```
go run transaction_checker.go -load=10 -filename=transactions.csv
```