BEGIN;
ALTER TABLE answers
    DROP CONSTRAINT fk_block,
    ADD  CONSTRAINT fk_block
        FOREIGN KEY ( bot_uuid, state )
            REFERENCES blocks ( bot_uuid, state );
COMMIT;
