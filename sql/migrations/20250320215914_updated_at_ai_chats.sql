
-- Auto-generated trigger for table "ai_chats" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "ai_chats_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "ai_chats_updated_at_trigger"
BEFORE UPDATE ON "ai_chats"
FOR EACH ROW
EXECUTE PROCEDURE "ai_chats_updated_at_trigger"();
