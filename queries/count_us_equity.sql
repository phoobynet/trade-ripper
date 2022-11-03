select count(*) from 'us_equity' where timestamp >= timestamp_floor('d', now())
