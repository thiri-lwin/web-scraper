CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar,
  "last_name" varchar,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "search_results" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "keyword" varchar NOT NULL,
  "results" jsonb,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "search_results" ("user_id");

CREATE INDEX ON "search_results" ("user_id", "keyword");

ALTER TABLE "search_results" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");