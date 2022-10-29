CREATE TABLE 'trades'
(
  sy     SYMBOL      capacity    12000    CACHE,
  s float,
  p float,
  t      TIMESTAMP
) timestamp(t) PARTITION BY DAY;
