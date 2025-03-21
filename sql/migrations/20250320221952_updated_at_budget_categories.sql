
-- Auto-generated trigger for table "budget_categories" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "budget_categories_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "budget_categories_updated_at_trigger"
BEFORE UPDATE ON "budget_categories"
FOR EACH ROW
EXECUTE PROCEDURE "budget_categories_updated_at_trigger"();
