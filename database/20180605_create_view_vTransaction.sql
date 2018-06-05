CREATE VIEW vTransaction AS
    SELECT t.oid AS oid,
           t.date,
           t.item,
           CASE WHEN t.description IS NULL THEN '' ELSE t.description END AS description,
           CASE
             WHEN t.pay AND NOT t.income THEN 'Pay'
             WHEN t.income AND NOT t.pay THEN 'Income'
             ELSE 'TBD'
           END AS direction,
           c.name AS currency,
           t.amount,
           p.name AS payment,
           b.name AS bank
      FROM transactions t,
           currency c,
           payment p,
           bank b
     WHERE t.currency = c.id AND
           t.payment = p.id AND
           t.bank = b.id;
