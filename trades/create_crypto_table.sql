CREATE TABLE if not exists 'crypto'
(
  ticker    SYMBOL capacity 256 CACHE,
  size      float,
  price     float,
  tks       text,
  t TIMESTAMP
) timestamp(t) PARTITION BY DAY;