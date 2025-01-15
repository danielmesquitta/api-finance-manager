
-- Auto-generated trigger for table "institutions" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "institutions_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "institutions_updated_at_trigger"
BEFORE UPDATE ON "institutions"
FOR EACH ROW
EXECUTE PROCEDURE "institutions_updated_at_trigger"();
