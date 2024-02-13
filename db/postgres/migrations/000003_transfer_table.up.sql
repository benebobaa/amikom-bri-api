
CREATE TABLE "entries" (
                           id bigserial PRIMARY KEY,
                           account_id bigint NOT NULL REFERENCES "accounts" (id) ON DELETE CASCADE,
                           date DATE NOT NULL DEFAULT CURRENT_DATE,
                           amount bigint NOT NULL,
                           entry_type varchar NOT NULL,
                           created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                           deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE "transfers" (
                             id bigserial PRIMARY KEY,
                             from_account_id bigint NOT NULL REFERENCES "accounts" (id) ON DELETE CASCADE,
                             to_account_id bigint NOT NULL REFERENCES "accounts" (id) ON DELETE CASCADE,
                             amount bigint NOT NULL,
                             created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                             deleted_at TIMESTAMPTZ DEFAULT NULL
);

ALTER TABLE "users" ADD COLUMN "hashed_pin" varchar DEFAULT NULL;


