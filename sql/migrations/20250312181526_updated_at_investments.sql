
-- Auto-generated trigger for table "investments" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "investments_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "investments_updated_at_trigger"
BEFORE UPDATE ON "investments"
FOR EACH ROW
EXECUTE PROCEDURE "investments_updated_at_trigger"();
