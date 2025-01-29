-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    session_uuid UUID PRIMARY KEY UNIQUE
);

CREATE TABLE IF NOT EXISTS chats(
    chat_uuid UUID PRIMARY KEY UNIQUE,
    session_uuid UUID,
    read_only BOOLEAN DEFAULT FALSE,
    ttl INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_chats_session_uuid FOREIGN KEY (session_uuid) REFERENCES users (session_uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages(
    message_uuid UUID PRIMARY KEY UNIQUE,
    session_uuid UUID,
    chat_uuid UUID,
    text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_messages_chat_uuid FOREIGN KEY (chat_uuid) REFERENCES chats (chat_uuid) ON DELETE CASCADE,
    CONSTRAINT fk_messages_session_uuid FOREIGN KEY (session_uuid) REFERENCES users (session_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
