CREATE TABLE mailings (
    bot_uuid       VARCHAR(36)  NOT NULL,
    name           VARCHAR(128) NOT NULL,
    entry_key      VARCHAR(256) NOT NULL,
    required_state INTEGER      NOT NULL,

    CONSTRAINT fk_entry_key
        FOREIGN KEY ( bot_uuid, entry_key )
            REFERENCES entry_points ( bot_uuid, key )
                ON DELETE CASCADE
);
