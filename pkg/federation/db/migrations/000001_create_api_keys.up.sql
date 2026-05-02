CREATE TABLE IF NOT EXISTS api_keys
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id   TEXT    NOT NULL UNIQUE,
    secret_hash TEXT    NOT NULL,
    salt        TEXT    NOT NULL,
    created_at  TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')),
    revoked_at  TEXT
);

CREATE INDEX IF NOT EXISTS idx_api_keys_client_id ON api_keys (client_id);