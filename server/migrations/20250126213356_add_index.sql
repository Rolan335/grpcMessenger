-- +goose Up
-- +goose StatementBegin

CREATE INDEX idx_chats_created_at ON chats(created_at);

CREATE INDEX idx_message_chat_uuid_created_at ON messages(chat_uuid, created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chats_created_at;
DROP INDEX IF EXISTS idx_message_chat_uuid_created_at;
-- +goose StatementEnd
