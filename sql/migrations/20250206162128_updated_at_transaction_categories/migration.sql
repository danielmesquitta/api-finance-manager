
-- Auto-generated trigger for table "transaction_categories" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "transaction_categories_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "transaction_categories_updated_at_trigger"
BEFORE UPDATE ON "transaction_categories"
FOR EACH ROW
EXECUTE PROCEDURE "transaction_categories_updated_at_trigger"();
