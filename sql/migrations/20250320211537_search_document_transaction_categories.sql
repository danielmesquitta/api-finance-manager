-- Computes and updates the search_document field for a given transaction category.
CREATE OR REPLACE FUNCTION compute_transaction_category_search_document(p_category_id UUID) RETURNS VOID AS $$ BEGIN
UPDATE transaction_categories
SET search_document = to_tsvector('portuguese', COALESCE(name, ''))
WHERE id = p_category_id;
END;
$$ LANGUAGE plpgsql;
-- Trigger function which calls the compute function after an insert or update.
CREATE OR REPLACE FUNCTION transaction_category_search_document_trigger() RETURNS trigger AS $$ BEGIN -- Prevent recursive trigger calls.
  IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
PERFORM compute_transaction_category_search_document(NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Create a trigger on the transaction_categories table to update the search_document for each row.
CREATE TRIGGER transaction_category_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON transaction_categories FOR EACH ROW EXECUTE FUNCTION transaction_category_search_document_trigger();