CREATE TABLE conversations (
    id                       CHAR(36) PRIMARY KEY,
    title                    TEXT NOT NULL,
    description              TEXT NOT NULL,
    create_time              TIMESTAMP(6) NOT NULL,
    update_time              TIMESTAMP(6) NOT NULL,
    last_message_id          CHAR(36),
    last_message_create_time TIMESTAMP(6),
    version                  INT NOT NULL,

    INDEX idx_conversations_last_activity_time ((COALESCE(last_message_create_time, update_time)))
) ENGINE = InnoDB;

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY
) ENGINE = InnoDB;

CREATE TABLE conversation_users (
    id              CHAR(36) PRIMARY KEY,
    conversation_id CHAR(36) NOT NULL,
    user_id         CHAR(36) NOT NULL,
    role            VARCHAR(16) NOT NULL,
    create_time     TIMESTAMP(6) NOT NULL,
    update_time     TIMESTAMP(6) NOT NULL,
    version         INT NOT NULL,

    INDEX idx_conversation_users_conversation_id_create_time (conversation_id, create_time),

    UNIQUE KEY uq_conversation_users_conversation_id_user_id (conversation_id, user_id),

    CONSTRAINT fk_conversation_users_conversation_id
        FOREIGN KEY (conversation_id)
        REFERENCES conversations (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_conversation_users_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE messages (
    id              CHAR(36) PRIMARY KEY,
    conversation_id CHAR(36) NOT NULL,
    user_id         CHAR(36),
    content         MEDIUMTEXT NOT NULL,
    create_time     TIMESTAMP(6) NOT NULL,
    update_time     TIMESTAMP(6) NOT NULL,
    version         INT NOT NULL,

    INDEX idx_messages_conversation_id_create_time (conversation_id, create_time),

    CONSTRAINT fk_messages_conversation_id
        FOREIGN KEY (conversation_id)
        REFERENCES conversations (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_messages_user_id
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE SET NULL
) ENGINE = InnoDB;

ALTER TABLE conversations
    ADD CONSTRAINT fk_conversations_last_message_id
        FOREIGN KEY (last_message_id)
        REFERENCES messages (id)
        ON DELETE SET NULL;

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     DATETIME(6) NOT NULL,

    INDEX idx_outbox_events_occur_time (occur_time)
) ENGINE = InnoDB;
