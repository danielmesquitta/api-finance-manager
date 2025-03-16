INSERT INTO users (
    id,
    name,
    email,
    tier,
    avatar,
    subscription_expires_at,
    synchronized_at
  )
VALUES (
    'fdfdc888-da64-4988-8ad3-f739862c4ceb'::uuid,
    'John Doe',
    'johndoe@email.com',
    'PREMIUM',
    'https://avatar.iran.liara.run/public/15',
    NOW() + INTERVAL '1 month',
    '2025-03-11 00:00:00.000 -0300'
  ),
  (
    '5b4694a9-c810-41a2-bca6-74c3f3850fe7'::uuid,
    'Jane Doe',
    'janedoe@email.com',
    'FREE',
    'https://avatar.iran.liara.run/public/82',
    NULL,
    '2025-03-11 00:00:00.000 -0300'
  )