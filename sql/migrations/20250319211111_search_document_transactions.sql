-- Computes and updates the search_document field for a specified transaction,
-- combining text vectors from the transaction name, its institution name,
-- its category name, and its payment method name.
CREATE OR REPLACE FUNCTION compute_transaction_search_document(p_transaction_id UUID) RETURNS VOID AS $$ BEGIN
UPDATE transactions
SET search_document = to_tsvector('portuguese', COALESCE(name, '')) || to_tsvector(
    'portuguese',
    COALESCE(
      (
        SELECT name
        FROM institutions
        WHERE institutions.id = transactions.institution_id
      ),
      ''
    )
  ) || to_tsvector(
    'portuguese',
    COALESCE(
      (
        SELECT name
        FROM transaction_categories
        WHERE transaction_categories.id = transactions.category_id
      ),
      ''
    )
  ) || to_tsvector(
    'portuguese',
    COALESCE(
      (
        SELECT name
        FROM payment_methods
        WHERE payment_methods.id = transactions.payment_method_id
      ),
      ''
    )
  )
WHERE id = p_transaction_id;
END;
$$ LANGUAGE plpgsql;
-- When an institution record is updated, this function finds all transactions
-- linked to that institution and recomputes their search_document.
CREATE OR REPLACE FUNCTION institution_trigger_for_transaction_search_document() RETURNS trigger AS $$
DECLARE rec RECORD;
BEGIN FOR rec IN
SELECT id
FROM transactions
WHERE institution_id = NEW.id LOOP PERFORM compute_transaction_search_document(rec.id);
END LOOP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Calls the above function after every update on the institutions table.
CREATE TRIGGER institution_trigger_for_transaction_search_document
AFTER
UPDATE ON institutions FOR EACH ROW EXECUTE FUNCTION institution_trigger_for_transaction_search_document();
-- When a transaction category is updated, this function finds all transactions
-- associated with that category and recomputes their search_document.
CREATE OR REPLACE FUNCTION category_trigger_for_transaction_search_document() RETURNS trigger AS $$
DECLARE rec RECORD;
BEGIN FOR rec IN
SELECT id
FROM transactions
WHERE category_id = NEW.id LOOP PERFORM compute_transaction_search_document(rec.id);
END LOOP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Calls the above function after every update on the transaction_categories table.
CREATE TRIGGER category_trigger_for_transaction_search_document
AFTER
UPDATE ON transaction_categories FOR EACH ROW EXECUTE FUNCTION category_trigger_for_transaction_search_document();
-- When a payment method is updated, this function finds all transactions
-- that use that payment method and recomputes their search_document.
CREATE OR REPLACE FUNCTION payment_method_trigger_for_transaction_search_document() RETURNS trigger AS $$
DECLARE rec RECORD;
BEGIN FOR rec IN
SELECT id
FROM transactions
WHERE payment_method_id = NEW.id LOOP PERFORM compute_transaction_search_document(rec.id);
END LOOP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Calls the above function after every update on the payment_methods table.
CREATE TRIGGER payment_method_trigger_for_transaction_search_document
AFTER
UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION payment_method_trigger_for_transaction_search_document();
-- When a transaction is inserted or updated, this function recomputes its search_document.
-- It also prevents recursive trigger calls by checking pg_trigger_depth.
CREATE OR REPLACE FUNCTION transaction_search_document_trigger() RETURNS trigger AS $$ BEGIN IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
PERFORM compute_transaction_search_document(NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Calls the above function after every insert or update on the transactions table.
CREATE TRIGGER transaction_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION transaction_search_document_trigger();