-- migrate утилита не оборачивает файл в одну транзакцию.
BEGIN;
    ALTER TABLE blocks
        ALTER COLUMN next_state
            DROP NOT NULL;

    UPDATE blocks
        SET   next_state = NULL
        WHERE next_state = 0;

    ALTER TABLE blocks
        DROP CONSTRAINT fk_bot_uuid,
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
END;
