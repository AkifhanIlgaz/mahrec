CREATE TABLE IF NOT EXISTS clients (
    id         TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS goals (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id     TEXT NOT NULL REFERENCES clients (id) ON DELETE CASCADE,
    title         TEXT NOT NULL,
    target_count  INTEGER NOT NULL,
    current_count INTEGER NOT NULL DEFAULT 0,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    completed_at  DATETIME
);

CREATE INDEX IF NOT EXISTS idx_goals_client_id ON goals (client_id);
