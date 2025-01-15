
-- Auto-generated trigger for table "transactions" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "transactions_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "transactions_updated_at_trigger"
BEFORE UPDATE ON "transactions"
FOR EACH ROW
EXECUTE PROCEDURE "transactions_updated_at_trigger"();
