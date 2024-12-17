-- AlterTable
ALTER TABLE "accounts" ALTER COLUMN "balance" SET DATA TYPE DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "budget_categories" ALTER COLUMN "amount" SET DATA TYPE DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "budgets" ALTER COLUMN "amount" SET DATA TYPE DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "credit_cards" ALTER COLUMN "limit" SET DATA TYPE DOUBLE PRECISION,
ALTER COLUMN "available_limit" SET DATA TYPE DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "investments" ALTER COLUMN "amount" SET DATA TYPE DOUBLE PRECISION,
ALTER COLUMN "rate" SET DATA TYPE DOUBLE PRECISION;

-- AlterTable
ALTER TABLE "transactions" ALTER COLUMN "amount" SET DATA TYPE DOUBLE PRECISION;