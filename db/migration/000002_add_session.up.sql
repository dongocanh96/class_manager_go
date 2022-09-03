CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NUll,
    "user_agent" varchar NOT NUll,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NUll DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");