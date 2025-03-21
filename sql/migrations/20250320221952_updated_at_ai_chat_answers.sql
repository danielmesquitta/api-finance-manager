
-- Auto-generated trigger for table "ai_chat_answers" to update @updatedAt columns
CREATE OR REPLACE FUNCTION "ai_chat_answers_updated_at_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER "ai_chat_answers_updated_at_trigger"
BEFORE UPDATE ON "ai_chat_answers"
FOR EACH ROW
EXECUTE PROCEDURE "ai_chat_answers_updated_at_trigger"();
