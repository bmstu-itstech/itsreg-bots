DO $$ BEGIN
    CREATE TYPE BOT_STATUS AS ENUM ('started', 'stopped', 'failed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS bots (
    uuid       VARCHAR(36)  PRIMARY KEY,
    name       VARCHAR(256) NOT NULL,
    token      VARCHAR(256) NOT NULL,
    status     BOT_STATUS   NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL
);

DO $$ BEGIN
    CREATE TYPE BLOCK_TYPE AS ENUM ('message', 'question', 'selection');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS blocks (
    bot_uuid   VARCHAR(36) NOT NULL,
    state      INTEGER     NOT NULL,

    type       BLOCK_TYPE   NOT NULL,
    next_state INTEGER      NOT NULL,
    title      VARCHAR(256) NOT NULL,
    text       TEXT         NOT NULL,

    PRIMARY KEY ( bot_uuid, state ),

    CONSTRAINT fk_bot_uuid
        FOREIGN KEY ( bot_uuid )
            REFERENCES bots ( uuid )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS entry_points (
    bot_uuid VARCHAR(36)  NOT NULL,
    key      VARCHAR(256) NOT NULL,
    state    INTEGER      NOT NULL,

    PRIMARY KEY ( bot_uuid, key ),

    CONSTRAINT fk_block
        FOREIGN KEY ( bot_uuid, state )
            REFERENCES blocks ( bot_uuid, state )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS options (
    bot_uuid VARCHAR(36)  NOT NULL,
    state    INTEGER      NOT NULL,
    next     INTEGER      NOT NULL,
    text     VARCHAR(256) NOT NULL,

    CONSTRAINT fk_block
    FOREIGN KEY ( bot_uuid, state )
        REFERENCES blocks ( bot_uuid, state )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS participants (
    bot_uuid VARCHAR(36) NOT NULL,
    user_id  BIGINT      NOT NULL,
    state    INTEGER     NOT NULL,

    PRIMARY KEY ( bot_uuid, user_id ),

    CONSTRAINT fk_bot_uuid
        FOREIGN KEY ( bot_uuid )
            REFERENCES bots ( uuid )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS answers (
    bot_uuid VARCHAR(36) NOT NULL,
    user_id  BIGINT      NOT NULL,
    state    INTEGER     NOT NULL,
    text     TEXT        NOT NULL,

    PRIMARY KEY ( bot_uuid, user_id, state ),

    CONSTRAINT fk_participant
        FOREIGN KEY ( bot_uuid, user_id )
            REFERENCES participants ( bot_uuid, user_id )
            ON DELETE CASCADE,

    CONSTRAINT fk_block
        FOREIGN KEY ( bot_uuid, state )
            REFERENCES blocks ( bot_uuid, state )
            ON DELETE NO ACTION
);
