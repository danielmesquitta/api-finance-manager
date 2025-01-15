
-- Auto-generated trigger for table "accounts" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "accounts_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "accounts_updated_at_trigger"
BEFORE UPDATE ON "accounts"
FOR EACH ROW
EXECUTE PROCEDURE "accounts_updated_at_trigger"();
