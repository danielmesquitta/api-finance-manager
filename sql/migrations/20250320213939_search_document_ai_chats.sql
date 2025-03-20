-- Computes and updates the search_document for a given AI Chat.
-- It concatenates the AI chat title (if any) and the text from all related
-- ai_chat_messages and ai_chat_answers, converting each part to a tsvector.
CREATE OR REPLACE FUNCTION compute_ai_chat_search_document(p_ai_chat_id UUID) RETURNS VOID AS $$ BEGIN
UPDATE ai_chats
SET search_document = to_tsvector('portuguese', COALESCE(title, '')) || (
    SELECT COALESCE(
        to_tsvector(
          'portuguese',
          string_agg(aggregated_msgs.message, ' ')
        ),
        to_tsvector('portuguese', '')
      )
    FROM (
        -- Get all messages from this chat
        SELECT message
        FROM ai_chat_messages
        WHERE ai_chat_id = p_ai_chat_id
        UNION ALL
        -- Get all answer messages for messages in this chat
        SELECT message
        FROM ai_chat_answers
        WHERE ai_chat_message_id IN (
            SELECT id
            FROM ai_chat_messages
            WHERE ai_chat_id = p_ai_chat_id
          )
      ) AS aggregated_msgs
  )
WHERE id = p_ai_chat_id;
END;
$$ LANGUAGE plpgsql;
-----------------------------------------------------------
-- Trigger function for ai_chat_messages inserts/updates.
-- When an AI chat message is created or updated, it updates the
-- parent ai_chat's search_document.
CREATE OR REPLACE FUNCTION ai_chat_message_search_document_trigger() RETURNS trigger AS $$ BEGIN -- Prevent recursive trigger calls
  IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
PERFORM compute_ai_chat_search_document(NEW.ai_chat_id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-----------------------------------------------------------
-- Trigger function for ai_chat_answers inserts/updates.
-- When an answer is created or updated, it fetches the related ai_chat_id
-- (via the associated ai_chat_message) and updates the corresponding
-- ai_chat search_document.
CREATE OR REPLACE FUNCTION ai_chat_answer_search_document_trigger() RETURNS trigger AS $$
DECLARE v_ai_chat_id UUID;
BEGIN IF pg_trigger_depth() > 1 THEN RETURN NEW;
END IF;
SELECT ai_chat_id INTO v_ai_chat_id
FROM ai_chat_messages
WHERE id = NEW.ai_chat_message_id;
PERFORM compute_ai_chat_search_document(v_ai_chat_id);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-----------------------------------------------------------
-- Trigger on ai_chat_messages: fire after insert or update.
CREATE TRIGGER ai_chat_message_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON ai_chat_messages FOR EACH ROW EXECUTE FUNCTION ai_chat_message_search_document_trigger();
-----------------------------------------------------------
-- Trigger on ai_chat_answers: fire after insert or update.
CREATE TRIGGER ai_chat_answer_search_document_trigger
AFTER
INSERT
  OR
UPDATE ON ai_chat_answers FOR EACH ROW EXECUTE FUNCTION ai_chat_answer_search_document_trigger();