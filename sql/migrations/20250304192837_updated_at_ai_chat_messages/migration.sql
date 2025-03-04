
-- Auto-generated trigger for table "ai_chat_messages" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "ai_chat_messages_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "ai_chat_messages_updated_at_trigger"
BEFORE UPDATE ON "ai_chat_messages"
FOR EACH ROW
EXECUTE PROCEDURE "ai_chat_messages_updated_at_trigger"();
