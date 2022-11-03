CREATE TABLE if not exists 'us_equity'
(
  ticker    SYMBOL capacity 12000 CACHE,
  size      float,
  price     float,
  timestamp TIMESTAMP
) timestamp(timestamp) PARTITION BY DAY;