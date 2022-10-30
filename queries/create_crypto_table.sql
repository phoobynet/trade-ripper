CREATE TABLE if not exists 'crypto'
(
  sy  SYMBOL capacity 256 CACHE,
  s   float,
  p   float,
  tks text,
  b   text,
  q   text,
  t   TIMESTAMP
) timestamp(t) PARTITION BY DAY;