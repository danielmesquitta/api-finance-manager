-- Computes and updates the search_document field for a given payment method.
CREATE OR REPLACE FUNCTION compute_payment_method_search_document(p_payment_method_id UUID) RETURNS VOID AS $$ BEGIN
UPDATE payment_methods
SET search_document = to_tsvector('portuguese', COALESCE(unaccent(name), ''))
WHERE id = p_payment_method_id;
END;
$$ LANGUAGE plpgsql;
-- Trigger function which calls the compute function after an insert or update.
CREATE OR REPLACE FUNCTION payment_method_search_document_trigger() RETURNS trigger AS $$ BEGIN -- Prevent recursive trigger calls
  IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
PERFORM compute_payment_method_search_document(NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Create a trigger on the payment_methods table to update the search_document for each row.
CREATE TRIGGER payment_method_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION payment_method_search_document_trigger();