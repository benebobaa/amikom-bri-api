-- Drop the added column in the "users" table
ALTER TABLE "users" DROP COLUMN IF EXISTS "hashed_pin";

-- Drop the "transfers" table
DROP TABLE IF EXISTS "transfers";

-- Drop the "entries" table
DROP TABLE IF EXISTS "entries";
