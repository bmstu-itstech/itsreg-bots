CREATE TABLE IF NOT EXISTS bots(
    id    BIGSERIAL     PRIMARY KEY,
    name  VARCHAR(256)  NOT NULL,
    token VARCHAR(256)  NOT NULL,
    start BIGINT        NOT NULL
);

CREATE TYPE BLOCK_TYPE AS ENUM ('message', 'question', 'selection');

CREATE TABLE IF NOT EXISTS blocks(
    type     BLOCK_TYPE NOT NULL,
    state    BIGINT     NOT NULL,
    bot_id   BIGINT     NOT NULL,
    default_ BIGINT,

    title  VARCHAR NOT NULL,
    text   VARCHAR NOT NULL,

    PRIMARY KEY ( state, bot_id ),

    CONSTRAINT bot_id_fk
        FOREIGN KEY ( bot_id )
            REFERENCES bots ( id )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS options(
    bot_id BIGINT  NOT NULL,
    state  BIGINT  NOT NULL,
    value_ VARCHAR NOT NULL,
    next   BIGINT  NOT NULL,

    CONSTRAINT block_id_fk
        FOREIGN KEY ( bot_id, state )
            REFERENCES blocks ( bot_id, state )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS participants(
    user_id BIGINT NOT NULL,
    bot_id  BIGINT NOT NULL,

    current BIGINT NOT NULL,

    PRIMARY KEY ( user_id, bot_id ),

    CONSTRAINT block_id_fk
        FOREIGN KEY ( bot_id, current )
            REFERENCES blocks ( bot_id, state )
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS answers(
    user_id BIGINT NOT NULL,
    bot_id  BIGINT NOT NULL,
    state   BIGINT NOT NULL,

    value_  VARCHAR NOT NULL,

    PRIMARY KEY ( user_id, bot_id, state ),

    CONSTRAINT block_id_fk
        FOREIGN KEY ( bot_id, state )
            REFERENCES blocks ( bot_id, state )
            ON DELETE CASCADE,

    CONSTRAINT participant_id_fk
        FOREIGN KEY ( bot_id, user_id )
            REFERENCES participants ( bot_id, user_id )
            ON DELETE CASCADE
);
