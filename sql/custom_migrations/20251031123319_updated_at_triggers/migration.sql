-- CreateFunction
CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "transactions" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "categories" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "accounts" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "user_accounts" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "budgets" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
-- CreateTrigger
CREATE TRIGGER update_modified_time BEFORE
UPDATE ON "budget_categories" FOR EACH ROW EXECUTE PROCEDURE update_modified_column();