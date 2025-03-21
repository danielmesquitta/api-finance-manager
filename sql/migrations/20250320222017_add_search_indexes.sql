-- Create IMMUTABLE wrapper for unaccent
CREATE OR REPLACE FUNCTION indexed_unaccent(text) RETURNS text AS $$
SELECT public.unaccent($1);
$$ LANGUAGE sql IMMUTABLE;
-- CreateIndex
CREATE INDEX idx_transactions_name_unaccent_trgm ON transactions USING gin (indexed_unaccent(name) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_payment_methods_name_unaccent_trgm ON payment_methods USING gin (indexed_unaccent(name) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_transaction_categories_name_unaccent_trgm ON transaction_categories USING gin (indexed_unaccent(name) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_institutions_name_unaccent_trgm ON institutions USING gin (indexed_unaccent(name) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_ai_chats_title_unaccent_trgm ON ai_chats USING gin (indexed_unaccent(title) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_ai_chat_messages_message_unaccent_trgm ON ai_chat_messages USING gin (indexed_unaccent(message) gin_trgm_ops);
-- CreateIndex
CREATE INDEX idx_ai_chat_answers_message_unaccent_trgm ON ai_chat_answers USING gin (indexed_unaccent(message) gin_trgm_ops);