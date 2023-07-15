ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "username_of_user_currency_key";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_username_fkey";

DROP TABLE IF EXISTS "users";
