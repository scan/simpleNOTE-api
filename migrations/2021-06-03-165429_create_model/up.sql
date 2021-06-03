CREATE TABLE IF NOT EXISTS accounts (
    id uuid NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    encrypted_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    last_login TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS account_tokens (
    token uuid NOT NULL,
    account_id uuid NOT NULL,
    user_agent VARCHAR(255) NULL DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (token),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_account_tokens_last_used ON account_tokens (account_id, last_used_at DESC);

CREATE TABLE IF NOT EXISTS notes (
    id uuid NOT NULL,
    account_id uuid NOT NULL,
    title VARCHAR(255) NOT NULL,
    contents TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_notes_created_at ON notes (account_id, deleted_at, created_at DESC);
