
-- Auto-generated trigger for table "user_institutions" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "user_institutions_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "user_institutions_updated_at_trigger"
BEFORE UPDATE ON "user_institutions"
FOR EACH ROW
EXECUTE PROCEDURE "user_institutions_updated_at_trigger"();
