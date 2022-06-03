# chaos-btcusd

This is a project that fetch BTC-USD price from public API and stores them in a database. It also serves endpoints to get the price.

<br />


## To run in docker

```
$ docker build -t test-server .
$ docker run --rm -p 80:80 test-server
```

<br />

## To call the endpoints

1. Get last price of BTC-USD
```
$ curl localhost:80/price
```

1. Get a price given timestamp
```
$ curl localhsot:80/price?timestamp=2022-06-01T18:39:47Z
```

3. Get average price in a time range
```
$ curl localhsot:80/price?from=2022-06-01T18:39:04Z&to=2023-06-01T18:47:47Z
```
