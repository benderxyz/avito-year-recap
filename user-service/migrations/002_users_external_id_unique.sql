CREATE UNIQUE INDEX IF NOT EXISTS users_external_id_key
    ON users (external_id)
    WHERE external_id <> '';
