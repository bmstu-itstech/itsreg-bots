-- migrate утилита не оборачивает файл в единую транзакцию.
BEGIN;
    ALTER TABLE blocks
        ALTER COLUMN next_state
            DROP NOT NULL;

    UPDATE blocks
        SET   next_state = NULL
        WHERE next_state = 0;

    ALTER TABLE blocks
        ADD CONSTRAINT fk_next_state
            FOREIGN KEY ( bot_uuid, next_state )
                REFERENCES blocks ( bot_uuid, state )
                ON DELETE CASCADE;

    ALTER TABLE options
        ALTER COLUMN next
            DROP NOT NULL;

    UPDATE options
        SET   next = NULL
        WHERE next = 0;

    ALTER TABLE options
        ADD CONSTRAINT fk_next_state
            FOREIGN KEY ( bot_uuid, next )
                REFERENCES blocks ( bot_uuid, state )
                ON DELETE CASCADE;

    ALTER TABLE participants
        ALTER COLUMN state
            DROP NOT NULL;

    UPDATE participants
        SET   state = NULL
        WHERE state = 0;

    ALTER TABLE participants
        DROP CONSTRAINT IF EXISTS fk_bot_uuid,
        ADD CONSTRAINT fk_block_state
            FOREIGN KEY ( bot_uuid, state )
                REFERENCES blocks ( bot_uuid, state )
                ON DELETE CASCADE;

    ALTER TABLE mailings
        ALTER COLUMN required_state DROP NOT NULL;

    UPDATE mailings
        SET   required_state = NULL
        WHERE required_state = 0;

    ALTER TABLE mailings
        ADD CONSTRAINT fk_required_state
            FOREIGN KEY ( bot_uuid, required_state )
                REFERENCES blocks ( bot_uuid, state );
END;
