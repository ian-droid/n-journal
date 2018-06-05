BEGIN;

SELECT pay, count(*) FROM transactions GROUP BY pay;
SELECT income, count(*) FROM transactions GROUP BY income;
UPDATE transactions set pay = 1 where pay = 'True';
UPDATE transactions set pay = 0 where pay = 'False';
UPDATE transactions set income = 1 where income = 'True';
UPDATE transactions set income = 0 where income = 'False';

COMMIT;
