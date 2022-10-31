select count(*) from 'us_equity' where t >= timestamp_floor('d', now())
