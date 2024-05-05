CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "refresh_token" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT False,
    "user_agent" varchar NOT NULL,
    "user_ip" varchar NOT NULL,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT fk_user_id FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);