
-- Auto-generated trigger for table "users" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "users_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "users_updated_at_trigger"
BEFORE UPDATE ON "users"
FOR EACH ROW
EXECUTE PROCEDURE "users_updated_at_trigger"();
