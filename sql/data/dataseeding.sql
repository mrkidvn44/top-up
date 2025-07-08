insert into "provider" (created_at, updated_at, deleted_at,name, code, logo_url, status) VALUES
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 'Viettel', 'VTL', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null,'MobiFone', 'MBF', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 'Vinaphone', 'VNP', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null,'Vietnamobile', 'VNM', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null,'Wintel', 'WNT', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null,'Itel', 'ITL', 'example.com.vn', 'active'),
(CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null,'Gmobile', 'GML', 'example.com.vn', 'inactive');
INSERT INTO cash_back (code, type, value, created_at, updated_at)
VALUES
  -- Percentage-based cashback
  ('CB5P', 'percentage', 5, NOW(), NOW()),      -- 5%
  ('CB10P', 'percentage', 10, NOW(), NOW()),    -- 10%
  ('CB15P', 'percentage', 15, NOW(), NOW()),    -- 15%
  ('CB20P', 'percentage', 20, NOW(), NOW()),    -- 20%
  ('CB25P', 'percentage', 25, NOW(), NOW()),    -- 25%
  ('CB30P', 'percentage', 30, NOW(), NOW()),    -- 30%
  ('CB50P', 'percentage', 50, NOW(), NOW()),    -- 50%
  -- Fixed-amount cashback
  ('CB1F', 'fixed', 1000, NOW(), NOW()),        -- 1,000 VND
  ('CB2F', 'fixed', 2000, NOW(), NOW()),        -- 2,000 VND
  ('CB5F', 'fixed', 5000, NOW(), NOW()),        -- 5,000 VND
  ('CB10F', 'fixed', 10000, NOW(), NOW()),      -- 10,000 VND
  ('CB20F', 'fixed', 20000, NOW(), NOW()),      -- 20,000 VND
  ('CB50F', 'fixed', 50000, NOW(), NOW()),      -- 50,000 VND
  ('CB100F', 'fixed', 100000, NOW(), NOW());    -- 100,000 VND

INSERT INTO sku (provider_code, cash_back_code, price, created_at, updated_at)
VALUES
  -- Viettel: all prices, mix of cashback
  ('VTL', NULL, 10000, NOW(), NOW()),
  ('VTL', NULL, 20000, NOW(), NOW()),
  ('VTL', 'CB5P', 30000, NOW(), NOW()),
  ('VTL', 'CB10P', 50000, NOW(), NOW()),
  ('VTL', 'CB15P', 100000, NOW(), NOW()),
  ('VTL', 'CB20F', 200000, NOW(), NOW()),
  ('VTL', 'CB25P', 300000, NOW(), NOW()),
  ('VTL', 'CB50F', 500000, NOW(), NOW()),
  ('VTL', 'CB30P', 1000000, NOW(), NOW()),
  ('VTL', 'CB100F', 2000000, NOW(), NOW()),
  ('VTL', 'CB50P', 3000000, NOW(), NOW()),
  ('VTL', 'CB100F', 5000000, NOW(), NOW()),
  -- MobiFone: only fixed cashback, 6 prices
  ('MBF', NULL, 10000, NOW(), NOW()),
  ('MBF', NULL, 20000, NOW(), NOW()),
  ('MBF', 'CB1F', 30000, NOW(), NOW()),
  ('MBF', 'CB2F', 50000, NOW(), NOW()),
  ('MBF', 'CB5F', 100000, NOW(), NOW()),
  ('MBF', 'CB10F', 200000, NOW(), NOW()),
  -- Vinaphone: 5 prices, only percentage cashback
  ('VNP', NULL, 10000, NOW(), NOW()),
  ('VNP', NULL, 20000, NOW(), NOW()),
  ('VNP', 'CB5P', 30000, NOW(), NOW()),
  ('VNP', 'CB10P', 50000, NOW(), NOW()),
  ('VNP', 'CB15P', 100000, NOW(), NOW()),
  -- Vietnamobile: 7 prices, mix
  ('VNM', NULL, 10000, NOW(), NOW()),
  ('VNM', NULL, 20000, NOW(), NOW()),
  ('VNM', 'CB1F', 30000, NOW(), NOW()),
  ('VNM', 'CB5P', 50000, NOW(), NOW()),
  ('VNM', 'CB10F', 100000, NOW(), NOW()),
  ('VNM', 'CB15P', 200000, NOW(), NOW()),
  ('VNM', 'CB2F', 300000, NOW(), NOW()),
  -- Wintel: all prices, only fixed cashback
  ('WNT', NULL, 10000, NOW(), NOW()),
  ('WNT', NULL, 20000, NOW(), NOW()),
  ('WNT', 'CB1F', 30000, NOW(), NOW()),
  ('WNT', 'CB2F', 50000, NOW(), NOW()),
  ('WNT', 'CB5F', 100000, NOW(), NOW()),
  ('WNT', 'CB10F', 200000, NOW(), NOW()),
  ('WNT', 'CB20F', 300000, NOW(), NOW()),
  ('WNT', 'CB50F', 500000, NOW(), NOW()),
  ('WNT', 'CB100F', 1000000, NOW(), NOW()),
  ('WNT', 'CB100F', 2000000, NOW(), NOW()),
  ('WNT', 'CB100F', 3000000, NOW(), NOW()),
  ('WNT', 'CB100F', 5000000, NOW(), NOW()),
  -- Itel: 5 prices, only percentage cashback
  ('ITL', NULL, 10000, NOW(), NOW()),
  ('ITL', NULL, 20000, NOW(), NOW()),
  ('ITL', 'CB5P', 30000, NOW(), NOW()),
  ('ITL', 'CB10P', 50000, NOW(), NOW()),
  ('ITL', 'CB15P', 100000, NOW(), NOW()),
  -- Gmobile: 6 prices, only fixed cashback
  ('GML', NULL, 10000, NOW(), NOW()),
  ('GML', NULL, 20000, NOW(), NOW()),
  ('GML', 'CB1F', 30000, NOW(), NOW()),
  ('GML', 'CB2F', 50000, NOW(), NOW()),
  ('GML', 'CB5F', 100000, NOW(), NOW()),
  ('GML', 'CB10F', 200000, NOW(), NOW()); 
INSERT INTO "user" (created_at, updated_at, deleted_at, first_name, last_name, password, phone_number )
VALUES
  ( CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, null, 'John', 'Doe', '$2a$14$QpdzQDCUaWutaB0WPJFz..LJPffpLAdsqZv/O8xPkwk77/9nRxeFG', '0358571286')