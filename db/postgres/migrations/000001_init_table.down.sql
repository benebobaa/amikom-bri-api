-- Drop Foreign Key Constraints

-- Drop foreign key constraint in "accounts" table
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "fk_accounts_owner";

-- Drop foreign key constraint in "verify_emails" table
ALTER TABLE "verify_emails" DROP CONSTRAINT IF EXISTS "fk_verify_emails_username";


-- Drop Tables

-- Drop Users Sessions table
DROP TABLE IF EXISTS "sessions";

-- Drop "verify_emails" table
DROP TABLE IF EXISTS "verify_emails";

-- Drop "accounts" table
DROP TABLE IF EXISTS "accounts";

-- Drop "users" table
DROP TABLE IF EXISTS "users";

