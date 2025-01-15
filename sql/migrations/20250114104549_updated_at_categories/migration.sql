
-- Auto-generated trigger for table "categories" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "categories_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "categories_updated_at_trigger"
BEFORE UPDATE ON "categories"
FOR EACH ROW
EXECUTE PROCEDURE "categories_updated_at_trigger"();
