## Environment Variables

#### Server Port: [Optional]

QLedger server by default runs in port `7000`, which can be overridden by the following:
```
export PORT=7000
```

#### Authentication Token:

QLedger API requests are authenticated using the secret token, which can be set using the following:
```
export LEDGER_AUTH_TOKEN=XXXXX
```

#### Database URL:

QLedger uses PostgreSQL database to store the accounts and transactions.

The PostgreSQL database URL can be set using:
```
export DATABASE_URL="postgres://localhost/ledgerdb?sslmode=disable"
```

For the purpose of running test cases, a separate database URL can be set using:
```
export TEST_DATABASE_URL="postgres://localhost/qw_ledger_test?sslmode=disable"
```

**Note:**

- The database URL can be in one of the mentioned formats here:
https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING

#### Sharing Load Balancer/Domain Name: [Optional]

In staging/production environments, the services are usually deployed in the same domain, differentiated and routed using the definite path prefixes.

To access all QLedger APIs with prefix `/qledger/api`, set the following:
```
export HOST_PREFIX=/qledger/api
```
