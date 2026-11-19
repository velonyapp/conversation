CREATE TABLE conversations (
    id          CHAR(36) PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL,
    create_time TIMESTAMP(6) NOT NULL,
    update_time TIMESTAMP(6) NOT NULL,
    version     INT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY
) ENGINE = InnoDB;

CREATE TABLE conversation_users (
    id              CHAR(36) PRIMARY KEY,
    conversation_id CHAR(36) NOT NULL,
    user_id         CHAR(36) NOT NULL,
    role            TINYINT UNSIGNED NOT NULL,
    create_time     TIMESTAMP(6) NOT NULL,
    update_time     TIMESTAMP(6) NOT NULL,
    version         INT NOT NULL,

    INDEX idx_conversation_users_conversation_id_role_update_time (conversation_id, role, update_time),

    UNIQUE KEY uq_conversation_users_user_id_conversation_id (user_id, conversation_id),

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

CREATE TABLE outbox_events (
    id             CHAR(36) PRIMARY KEY,
    type           TEXT NOT NULL,
    aggregate_id   TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    payload        JSON NOT NULL,
    occur_time     DATETIME(6) NOT NULL,

    INDEX idx_outbox_events_occur_time (occur_time)
) ENGINE = InnoDB;
