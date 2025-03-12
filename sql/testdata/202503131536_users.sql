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
    auth_id,
    open_finance_id
  )
VALUES (
    'fdfdc888-da64-4988-8ad3-f739862c4ceb'::uuid,
    'MOCK',
    'John Doe',
    'johndoe@email.com',
    true,
    'PREMIUM',
    'https://avatar.iran.liara.run/public/15',
    NOW() + INTERVAL '1 month',
    '2025-03-11 00:00:00.000 -0300',
    '6c2342aa-bdac-4efe-a31b-3a018072cff9',
    '5de68b44-d82b-482a-a3b6-b2189f201e6b'
  ),
  (
    '5b4694a9-c810-41a2-bca6-74c3f3850fe7'::uuid,
    'MOCK',
    'Jane Doe',
    'janedoe@email.com',
    true,
    'FREE',
    'https://avatar.iran.liara.run/public/82',
    NULL,
    '2025-03-11 00:00:00.000 -0300',
    '016aecbd-fae5-4ff0-9046-03b7eabf6a5c',
    '8d3ee953-7659-46a2-ab5f-2513ba960a05'
  )