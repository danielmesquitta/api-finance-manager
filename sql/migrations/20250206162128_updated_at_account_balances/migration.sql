
-- Auto-generated trigger for table "account_balances" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "account_balances_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "account_balances_updated_at_trigger"
BEFORE UPDATE ON "account_balances"
FOR EACH ROW
EXECUTE PROCEDURE "account_balances_updated_at_trigger"();
