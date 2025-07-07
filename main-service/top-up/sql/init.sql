DO $$
BEGIN
    DROP TYPE IF EXISTS cash_back_type CASCADE;
    DROP TYPE IF EXISTS provider_status CASCADE;
    DROP TYPE IF EXISTS purchase_history_status CASCADE;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cash_back_type') THEN
        CREATE TYPE cash_back_type AS ENUM ('percentage','fixed');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'provider_status') THEN
        CREATE TYPE provider_status AS ENUM ('active','inactive');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'purchase_history_status') THEN
        CREATE TYPE purchase_history_status AS ENUM ('pending','confirm','success','failed');
    END IF;
END $$;

