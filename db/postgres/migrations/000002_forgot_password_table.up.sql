CREATE TABLE "forgot_password" (
                                 id bigserial PRIMARY KEY,
                                 user_id UUID NOT NULL REFERENCES "users"(id) ON DELETE CASCADE,
                                 reset_token text NOT NULL,
                                 is_used boolean NOT NULL DEFAULT FALSE,
                                 request_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                 expiration_timestamp TIMESTAMPTZ,
                                 deleted_at TIMESTAMPTZ DEFAULT NULL
);
