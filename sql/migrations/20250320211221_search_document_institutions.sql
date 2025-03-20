-- Computes and updates the search_document field for a given institution.
CREATE OR REPLACE FUNCTION compute_institution_search_document(p_institution_id UUID) RETURNS VOID AS $$ BEGIN
UPDATE institutions
SET search_document = to_tsvector('portuguese', COALESCE(name, ''))
WHERE id = p_institution_id;
END;
$$ LANGUAGE plpgsql;
-- Trigger function which calls the compute function after an insert or update.
CREATE OR REPLACE FUNCTION institution_search_document_trigger() RETURNS trigger AS $$ BEGIN -- Prevent recursive trigger calls
  IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
PERFORM compute_institution_search_document(NEW.id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Create a trigger on the institutions table to update the search_document for each row.
CREATE TRIGGER institution_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON institutions FOR EACH ROW EXECUTE FUNCTION institution_search_document_trigger();