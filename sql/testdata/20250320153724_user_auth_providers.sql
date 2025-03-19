INSERT INTO user_auth_providers (
    id,
    user_id,
    external_id,
    provider,
    verified_email
  )
VALUES (
    'ca12d214-bbae-44ad-ae1a-dd07f063cf80'::uuid,
    'fdfdc888-da64-4988-8ad3-f739862c4ceb',
    '6c2342aa-bdac-4efe-a31b-3a018072cff9',
    'MOCK',
    true
  ),
  (
    '56b8d4b9-e654-4562-8301-5a7e56d3cf21'::uuid,
    '1698b17e-13a6-4ac9-9aa1-9408f29078af',
    '2a35fa25-2809-40d7-beeb-0d2766171b1d',
    'MOCK',
    true
  ),
  (
    '14e9d0fd-6777-4dd3-a4c1-12988af924c9'::uuid,
    '5b4694a9-c810-41a2-bca6-74c3f3850fe7',
    '016aecbd-fae5-4ff0-9046-03b7eabf6a5c',
    'MOCK',
    true
  );