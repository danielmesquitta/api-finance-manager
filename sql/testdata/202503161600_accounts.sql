INSERT INTO accounts (
    id,
    external_id,
    "name",
    "type",
    user_institution_id
  )
VALUES (
    '8567ca77-ac20-4526-b3d9-dbf380a1c00d'::uuid,
    '1f0d4c9b-47bb-4098-84ad-6d11835a7c5c',
    'gold',
    'CREDIT',
    'd237bbc3-8f60-4a78-9282-8e3f1dbe1630'::uuid
  ),
  (
    'e5f31705-cb65-42a5-9072-2b9b59e338a8'::uuid,
    '26199ce7-eddc-448f-9d88-6b768fe23499',
    'Nu Pagamentos S.A. - Instituição de Pagamento',
    'BANK',
    'd237bbc3-8f60-4a78-9282-8e3f1dbe1630'::uuid
  ),
  (
    'ac4d82a0-9eff-4936-8a2e-8d12591c9d00'::uuid,
    '0a823f25-600b-47f6-8a4f-dfae990f8a30',
    'BTG Investimentos',
    'BANK',
    '101df06d-086d-42a9-a589-7fbe56e552c2'::uuid
  ),
  (
    'c3894f46-33a6-47cd-85cf-ceaf4bb10895'::uuid,
    '46f860c6-5010-41e1-aee3-e3905a51793a',
    'BTG Banking',
    'BANK',
    '101df06d-086d-42a9-a589-7fbe56e552c2'::uuid
  );