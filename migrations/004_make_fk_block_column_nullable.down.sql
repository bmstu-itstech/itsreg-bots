-- migrate утилита не оборачивает файл в одну транзакцию.
BEGIN;
    ALTER TABLE options
        DROP CONSTRAINT fk_block;

    UPDATE options
        SET   next = 0
        WHERE next IS NULL;

    ALTER TABLE options
        ALTER COLUMN next SET NOT NULL;

    ALTER TABLE blocks
        DROP CONSTRAINT fk_next_state,
        ADD CONSTRAINT fk_bot_uuid
            FOREIGN KEY ( bot_uuid )
            REFERENCES bots ( uuid );

    UPDATE blocks
        SET    next_state = 0
        WHERE  next_state IS NULL;

    ALTER TABLE blocks
        ALTER COLUMN next_state SET NOT NULL;
END;
