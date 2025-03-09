INSERT INTO accounts (
    id,
    external_id,
    "name",
    "type",
    created_at,
    updated_at,
    deleted_at,
    user_id,
    institution_id
  )
VALUES (
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    '1f0d4c9b-47bb-4098-84ad-6d11835a7c5c',
    'gold',
    'CREDIT',
    '2025-03-09 14:11:19.681211-03',
    '2025-03-09 14:11:19.681211-03',
    NULL,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid
  ),
  (
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    '26199ce7-eddc-448f-9d88-6b768fe23499',
    'Nu Pagamentos S.A. - Instituição de Pagamento',
    'BANK',
    '2025-03-09 14:11:19.681211-03',
    '2025-03-09 14:11:19.681211-03',
    NULL,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid
  ),
  (
    'c09ca3e4-1630-4a8d-935b-6fc30171e34d'::uuid,
    '0a823f25-600b-47f6-8a4f-dfae990f8a30',
    'BTG Investimentos',
    'BANK',
    '2025-03-09 14:11:19.745342-03',
    '2025-03-09 14:11:19.745342-03',
    NULL,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '1202269c-ed03-4dfe-bbcd-c61d615a17b5'::uuid
  ),
  (
    '1ea748ad-b9b6-4b6e-b110-6fc86057ad12'::uuid,
    '46f860c6-5010-41e1-aee3-e3905a51793a',
    'BTG Banking',
    'BANK',
    '2025-03-09 14:11:19.745342-03',
    '2025-03-09 14:11:19.745342-03',
    NULL,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '1202269c-ed03-4dfe-bbcd-c61d615a17b5'::uuid
  );