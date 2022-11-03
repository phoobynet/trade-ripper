CREATE TABLE if not exists 'crypto'
(
  pair symbol capacity 256 CACHE,
  size      float,
  price     float,
  tks       text,
  timestamp TIMESTAMP
) timestamp(timestamp) PARTITION BY DAY;