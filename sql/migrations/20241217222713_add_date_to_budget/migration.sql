/*
  Warnings:

  - Added the required column `date` to the `budgets` table without a default value. This is not possible if the table is not empty.

*/
-- DropIndex
DROP INDEX "budgets_user_id_key";

-- AlterTable
ALTER TABLE "budgets" ADD COLUMN     "date" TIMESTAMPTZ NOT NULL;
