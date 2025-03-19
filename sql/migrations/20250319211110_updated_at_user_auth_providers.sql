
-- Auto-generated trigger for table "user_auth_providers" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "user_auth_providers_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "user_auth_providers_updated_at_trigger"
BEFORE UPDATE ON "user_auth_providers"
FOR EACH ROW
EXECUTE PROCEDURE "user_auth_providers_updated_at_trigger"();
