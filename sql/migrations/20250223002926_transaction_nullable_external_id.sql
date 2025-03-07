-- DropIndex
DROP INDEX "transactions_external_id_key";

-- AlterTable
ALTER TABLE "transactions" ALTER COLUMN "external_id" DROP NOT NULL;
