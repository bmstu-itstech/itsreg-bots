CREATE TABLE IF NOT EXISTS bots(
    id    BIGSERIAL     PRIMARY KEY,
    name  VARCHAR(256)  NOT NULL,
    token VARCHAR(256)  NOT NULL,
    start BIGSERIAL     NOT NULL
);

CREATE TYPE BLOCK_TYPE AS ENUM ('message', 'question', 'selection');

CREATE TABLE IF NOT EXISTS blocks(
    type         BLOCK_TYPE NOT NULL,
    state        BIGSERIAL  NOT NULL,
    bot_id       BIGSERIAL  NOT NULL,
    default_next BIGSERIAL,

    PRIMARY KEY ( state, bot_id ),

    CONSTRAINT bot_id_fk
    FOREIGN KEY ( bot_id )
        REFERENCES bots ( id )
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS options(
    bot_id BIGSERIAL NOT NULL,
    _value VARCHAR   NOT NULL,
    next   BIGSERIAL NOT NULL,

    CONSTRAINT bot_id_fk
        FOREIGN KEY ( bot_id )
        REFERENCES bots ( id )
        ON DELETE CASCADE
);
