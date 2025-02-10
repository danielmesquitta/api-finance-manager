/*
  Warnings:

  - You are about to drop the column `external_id` on the `users` table. All the data in the column will be lost.
  - Added the required column `auth_id` to the `users` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "users" DROP COLUMN "external_id",
ADD COLUMN     "auth_id" TEXT NOT NULL,
ADD COLUMN     "open_finance_id" TEXT;
