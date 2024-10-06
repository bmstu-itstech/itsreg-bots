-- migrate утилита не оборачивает файл в одну транзакцию.
BEGIN;
    ALTER TABLE mailings
        DROP CONSTRAINT IF EXISTS fk_required_state;

    UPDATE mailings
        SET   required_state = 0
        WHERE required_state IS NULL;

    ALTER TABLE mailings
        ALTER COLUMN required_state SET NOT NULL;

    ALTER TABLE participants
        DROP CONSTRAINT IF EXISTS fk_block_state,
        ADD CONSTRAINT fk_bot_uuid
            FOREIGN KEY ( bot_uuid )
                REFERENCES bots ( uuid );

    UPDATE participants
        SET   state = 0
        WHERE state IS NULL;

    ALTER TABLE participants
        ALTER COLUMN state SET NOT NULL;

    ALTER TABLE options
        DROP CONSTRAINT fk_next_state;

    UPDATE options
        SET   next = 0
        WHERE next IS NULL;

    ALTER TABLE options
        ALTER COLUMN next SET NOT NULL;

    ALTER TABLE blocks
        DROP CONSTRAINT fk_bot_uuid,
        DROP CONSTRAINT fk_next_state;

    UPDATE blocks
        SET    next_state = 0
        WHERE  next_state IS NULL;

    ALTER TABLE blocks
        ALTER COLUMN next_state SET NOT NULL;
END;
