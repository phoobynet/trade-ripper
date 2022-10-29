# Sipper Ripper CLI

## Danger here

This is a first stab (in Go) at doing this, are there a few untested things (there are no tests), and I need to make some QOL improvements, e.g. generating the table if it doesn't exist, etc.

**Use it at your own risk.**

## What is this?

If you pay for
an [Alpaca Markets data subscription](https://alpaca.markets/docs/api-references/market-data-api/stock-pricing-data/realtime/)
have access to all trades going through
the [Securities Information Processor](https://polygon.io/blog/understanding-the-sips/).

This app captures all trades (not quotes or bars), and writes the *Symbol*, *Price*, *Size* and *Timestamp* to
a [QuestDB](https://questdb.io/docs/) database. Fields such as trading conditions, exchange and tape are thrown away.

## Why QuestDB?

I tried a bunch of databases, SQLite, MongoDB (time series), Postgres, TimeScaleDB (postgres), but none were able to
perform as well as Quest with so little effort while being able to bulk insert and perform online queries at the same
time.

I run my QuestDB on Windows 11 Pro (yes, windows) AMD 5700g with 64GB of RAM and a PCI gen 4 NVME drive (but it's
limited to
PCI gen 3 speed because of the processor).

I've also had it running on linux in docker without any issues.

If you have an Apple Silicon Mac, I don't recommend using Docker (QuestDB image is not natively ARM). Use
the [binary download instead or install it with brew](https://questdb.io/docs/get-started/homebrew).

### Installation

Following the instructions on [QuestDB's](https://questdb.io/docs/get-started/binaries/) website.

One other thing you will need to do if you want to perform near real-time queries is to change
the [commit lag](https://questdb.io/docs/guides/out-of-order-commit-lag/).

This is a high throughput scenario, so well commit at least every 1 second, or if we have reach 10,000 uncommitted rows.  When the market opens, you can easily hit 30,000 - 40,000 trades a second.

Open `server.conf` and add the following lines

```
cairo.commit.lag=1000
cairo.max.uncommitted.rows=10000
```

Restart the QuestDB.

## Creating the table

Open a web browser on http://<questdb-host>:9000, and run the following to create the table.

```sql
CREATE TABLE 'crypto'
(
    sy     SYMBOL      capacity    256    CACHE,
    s float,
    p float,
    t      TIMESTAMP
) timestamp(t) PARTITION BY DAY;

```

```sql
CREATE TABLE 'us_equity'
(
    sy     SYMBOL      capacity    256    CACHE,
    s float,
    p float,
    t      TIMESTAMP
) timestamp(t) PARTITION BY DAY;

```

## Build and Run

Build it, then run it, making sure that port is
the [Influx DB inline protocol](https://questdb.io/docs/develop/insert-data/#influxdb-line-protocol) port (default
is `9009`)

```bash
go build

./trade-ripper --host my-questdb-host --port 9009
```





