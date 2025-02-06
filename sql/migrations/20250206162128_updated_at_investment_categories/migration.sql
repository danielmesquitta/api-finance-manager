
-- Auto-generated trigger for table "investment_categories" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "investment_categories_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "investment_categories_updated_at_trigger"
BEFORE UPDATE ON "investment_categories"
FOR EACH ROW
EXECUTE PROCEDURE "investment_categories_updated_at_trigger"();
