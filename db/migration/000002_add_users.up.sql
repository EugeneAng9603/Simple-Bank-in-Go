CREATE TABLE "users" (
  "username_of_user" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username_of_user");

-- CREATE UNIQUE INDEX ON "accounts" ("username", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "username_of_user_currency_key" UNIQUE ("username", "currency");

