BEGIN;

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    token_hash VARCHAR(72) NOT NULL,
    user_id INTEGER NOT NULL, -- to restrict number of simaltinious sesstions
    expires_at TIMESTAMPTZ NOT NULL -- mutable to logout
);

-- create unique index on token_hash
CREATE UNIQUE INDEX IF NOT EXISTS tokens_unique_idx_token_hash ON tokens (token_hash);

-- add user_id foreign key constraint
ALTER TABLE IF EXISTS tokens
    ADD CONSTRAINT tokens_fk_users
    FOREIGN KEY (user_id)
    REFERENCES users(id);

COMMIT;