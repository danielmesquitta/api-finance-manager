/*
  Warnings:

  - You are about to drop the column `description` on the `transactions` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "accounts" ALTER COLUMN "type" DROP DEFAULT;

-- AlterTable
ALTER TABLE "investments" ALTER COLUMN "type" DROP DEFAULT,
ALTER COLUMN "rateType" DROP DEFAULT;

-- AlterTable
ALTER TABLE "transactions" DROP COLUMN "description",
ALTER COLUMN "payment_method" DROP DEFAULT;

-- AlterTable
ALTER TABLE "users" ALTER COLUMN "tier" DROP DEFAULT;
