
-- Users Table
CREATE TABLE "users" (
                         username varchar PRIMARY KEY ,
                         email varchar NOT NULL,
                         full_name varchar NOT NULL,
                         hashed_password varchar NOT NULL,
                         is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Accounts Table
CREATE TABLE "accounts" (
                            id bigserial PRIMARY KEY,
                            owner varchar REFERENCES "users"(username),
                            balance bigint NOT NULL,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Verify Emails Table
CREATE TABLE "verify_emails" (
                                id bigserial PRIMARY KEY,
                                username varchar REFERENCES "users"(username),
                                email varchar NOT NULL,
                                secret_code varchar NOT NULL,
                                is_used bool NOT NULL DEFAULT false,
                                expired_at timestamptz NOT NULL,
                                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                deleted_at TIMESTAMPTZ DEFAULT NULL
);



-- Session table
CREATE TABLE "sessions" (
                            id UUID PRIMARY KEY,
                            username varchar REFERENCES "users"(username),
                            refresh_token TEXT NOT NULL ,
                            user_agent varchar NOT NULL,
                            client_ip varchar NOT NULL,
                            is_blocked boolean NOT NULL DEFAULT FALSE,
                            expired_at TIMESTAMPTZ NOT NULL ,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

);
