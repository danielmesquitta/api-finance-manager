
-- Auto-generated trigger for table "payment_methods" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "payment_methods_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "payment_methods_updated_at_trigger"
BEFORE UPDATE ON "payment_methods"
FOR EACH ROW
EXECUTE PROCEDURE "payment_methods_updated_at_trigger"();
