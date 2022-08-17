CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "is_teacher" boolean NOT NULL DEFAULT false
);

CREATE TABLE "homeworks" (
  "id" bigserial PRIMARY KEY,
  "teacher_id" bigint NOT NULL,
  "subject" varchar(256) NOT NULL,
  "title" varchar(256) NOT NULL,
  "file_name" varchar UNIQUE NOT NULL,
  "saved_path" varchar NOT NULL,
  "is_closed" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "closed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "solutions" (
  "id" bigserial PRIMARY KEY,
  "problem_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "file_name" varchar UNIQUE NOT NULL,
  "saved_path" varchar NOT NULL,
  "submited_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "from_user_id" bigint NOT NULL,
  "to_user_id" bigint NOT NULL,
  "content" varchar(5000) NOT NULL,
  "is_read" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "read_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "homeworks" ("id");

CREATE INDEX ON "homeworks" ("teacher_id");

CREATE INDEX ON "solutions" ("id");

CREATE INDEX ON "solutions" ("problem_id");

CREATE INDEX ON "solutions" ("user_id");

CREATE INDEX ON "solutions" ("problem_id", "user_id");

CREATE INDEX ON "messages" ("id");

CREATE INDEX ON "messages" ("from_user_id");

CREATE INDEX ON "messages" ("to_user_id");

CREATE INDEX ON "messages" ("from_user_id", "to_user_id");

ALTER TABLE "homeworks" ADD FOREIGN KEY ("teacher_id") REFERENCES "users" ("id");

ALTER TABLE "solutions" ADD FOREIGN KEY ("problem_id") REFERENCES "homeworks" ("id");

ALTER TABLE "solutions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "messages" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");
