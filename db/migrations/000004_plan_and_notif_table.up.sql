
CREATE TABLE "expenses_plans" (
            id bigserial PRIMARY KEY,
            user_id UUID NOT NULL REFERENCES "users"(id),
            title varchar NOT NULL ,
            description varchar default NULL,
            amount bigint not null,
            date DATE NOT NULL DEFAULT CURRENT_DATE,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE "notifications" (
            id bigserial PRIMARY KEY ,
            user_id UUID NOT NULL REFERENCES "users"(id),
            title varchar NOT NULL ,
            description varchar NOT NULL ,
            category varchar NOT NULL ,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            deleted_at TIMESTAMPTZ DEFAULT NULL
);