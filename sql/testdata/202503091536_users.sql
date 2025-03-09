INSERT INTO users (
    id,
    provider,
    "name",
    email,
    verified_email,
    tier,
    avatar,
    subscription_expires_at,
    synchronized_at,
    created_at,
    updated_at,
    deleted_at,
    auth_id,
    open_finance_id
  )
VALUES (
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'MOCK',
    'John Doe',
    'johndoe@email.com',
    true,
    'PREMIUM',
    'https://avatars.githubusercontent.com/u/60039311',
    NOW() + INTERVAL '1 month',
    '2025-03-07 21:00:00-03',
    '2025-03-09 14:11:19.322842-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '6c2342aa-bdac-4efe-a31b-3a018072cff9',
    '5de68b44-d82b-482a-a3b6-b2189f201e6b'
  );