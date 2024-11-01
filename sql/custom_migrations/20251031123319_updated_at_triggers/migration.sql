-- CreateFunction
CREATE OR REPLACE FUNCTION update_modified_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- CreateTrigger
DO $$
DECLARE tbl_name TEXT;
BEGIN FOR tbl_name IN
SELECT unnest(
    ARRAY [
            'users',
            'transactions',
            'categories',
            'accounts',
            'institutions',
            'budgets',
            'budget_categories'
        ]
  ) LOOP EXECUTE format(
    'CREATE TRIGGER update_modified_time
      BEFORE UPDATE ON %I
      FOR EACH ROW
      EXECUTE PROCEDURE update_modified_column();',
    tbl_name
  );
END LOOP;
END;
$$;
