CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE "users" (
                         id  UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                         username varchar NOT NULL,
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
                            user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
                            balance bigint NOT NULL,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Verify Emails Table
CREATE TABLE "verify_emails" (
                                id bigserial PRIMARY KEY,
                                user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
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
                            user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
                            refresh_token TEXT NOT NULL ,
                            user_agent varchar NOT NULL,
                            client_ip varchar NOT NULL,
                            is_blocked boolean NOT NULL DEFAULT FALSE,
                            expired_at TIMESTAMPTZ NOT NULL ,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMPTZ DEFAULT NULL

);
