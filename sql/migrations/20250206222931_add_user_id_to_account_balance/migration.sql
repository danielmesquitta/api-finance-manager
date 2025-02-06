/*
  Warnings:

  - Added the required column `user_id` to the `account_balances` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "account_balances" ADD COLUMN     "user_id" UUID NOT NULL;

-- AddForeignKey
ALTER TABLE "account_balances" ADD CONSTRAINT "account_balances_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
