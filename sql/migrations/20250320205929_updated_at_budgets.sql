
-- Auto-generated trigger for table "budgets" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "budgets_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "budgets_updated_at_trigger"
BEFORE UPDATE ON "budgets"
FOR EACH ROW
EXECUTE PROCEDURE "budgets_updated_at_trigger"();
