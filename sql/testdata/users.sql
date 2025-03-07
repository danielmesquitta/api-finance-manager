INSERT INTO users (
    id,
    auth_id,
    open_finance_id,
    provider,
    name,
    email,
    verified_email,
    tier,
    avatar,
    subscription_expires_at,
    synchronized_at,
    created_at,
    updated_at,
    deleted_at
  )
VALUES (
    '0a3ea07c-1f83-409f-87a8-5bbc23b1647d',
    '6c2342aa-bdac-4efe-a31b-3a018072cff9',
    NULL,
    'MOCK',
    'John Doe',
    'johndoe@email.com',
    TRUE,
    'PREMIUM',
    NULL,
    NOW() + INTERVAL '30 days',
    NULL,
    NOW(),
    NOW(),
    NULL
  )