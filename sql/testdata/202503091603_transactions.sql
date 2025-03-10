INSERT INTO transactions (
    id,
    external_id,
    "name",
    amount,
    is_ignored,
    "date",
    created_at,
    updated_at,
    deleted_at,
    payment_method_id,
    user_id,
    account_id,
    institution_id,
    category_id
  )
VALUES (
    '61edb9dd-c137-4e7e-8361-bc78a7ef864b'::uuid,
    '601c813f-756f-4d7f-84d6-ef4084656341',
    'Transferência enviada|LE POSTICHE',
    -94999,
    false,
    '2025-01-09 18:35:32.556-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '02701aac-b8db-4c7e-834c-6d4f4eab3399'::uuid
  ),
  (
    '571c215a-1ee2-4b1a-a316-3ffdd971340d'::uuid,
    '829b9140-f741-4d44-a418-4e60998221cb',
    'Resgate RDB',
    250000,
    false,
    '2025-02-06 17:47:43.203-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '262f50e1-a751-4184-9427-90a23f485482'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '029fc6cb-edcf-414c-9c81-9dd69c34e629'::uuid
  ),
  (
    '26657ab2-19cc-47a0-8af6-160f12737e14'::uuid,
    '557dae4e-dd1b-477b-8e10-93db882a7db7',
    'Resgate de Cashback',
    257,
    false,
    '2024-06-10 21:00:00.001-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '262f50e1-a751-4184-9427-90a23f485482'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '03bd0abc-7186-4eb3-9871-e4f624c535b8'::uuid
  ),
  (
    'a5707214-415e-4c0c-b8da-e1e225365151'::uuid,
    'c48750f4-d76a-4eb4-b37b-b4cbe260c6c1',
    'Transferência enviada',
    -27110,
    false,
    '2024-03-09 12:42:17.198-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '0c84d0a3-7336-4089-bc3d-756ce31c679a'::uuid
  ),
  (
    'f274a4be-1150-4542-896d-88239378b828'::uuid,
    'fda53877-70fa-46d6-b238-766c6b1afb75',
    'Carrefour Com',
    26867,
    false,
    '2024-08-29 08:32:14.001-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '12deb35c-0ce5-4d23-87a4-2f68fd77f019'::uuid
  ),
  (
    '8df11353-b6ec-42c5-9fec-84ae140d85cb'::uuid,
    'd6c663b3-e7c7-43e1-8d9b-52d4f8c2c83d',
    'Oficial*Oficialfarma 9/9',
    3083,
    false,
    '2024-09-08 00:00:00-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '1bd7db5b-5b8a-4ac1-82a1-c75e418a25c0'::uuid
  ),
  (
    '27d30d16-c585-49da-8370-bdd77c278295'::uuid,
    '48b3239e-b4b0-42d8-ac96-0278ac5a4df6',
    'PIX CASH IN EXTERNO',
    100000,
    false,
    '2025-03-03 07:44:58.714-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '1ea748ad-b9b6-4b6e-b110-6fc86057ad12'::uuid,
    '1202269c-ed03-4dfe-bbcd-c61d615a17b5'::uuid,
    '65583cfa-b72d-4fab-9de1-4ca9dfe11a4e'::uuid
  ),
  (
    'f68cadfa-b54c-4e37-857c-51db6bb0c465'::uuid,
    '95ae3369-2e64-4110-983d-3c8b45a6ef05',
    'Transferência enviada',
    -38961,
    false,
    '2024-03-07 10:25:00.082-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '70c89492-9977-42b4-a28c-b1e261c59615'::uuid
  ),
  (
    'eb3c0fc8-77bd-4130-83b6-af815d1a2956'::uuid,
    'cc318e61-9226-4f0d-8662-7b5de0b50e8e',
    'Crédito de "Google Storage"',
    -299,
    false,
    '2024-07-24 09:00:00.001-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    '9db43714-3025-494f-9578-4feb5a69681e'::uuid
  ),
  (
    '79260d65-66bb-476e-85db-1fce518b6aae'::uuid,
    'a6b11a51-aaa9-43d2-acb5-7bbd32bcbe1e',
    'Transferência enviada|Julia Fernandes Avelar',
    -41000,
    true,
    '2024-12-07 08:32:45.66-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    'a910d4f6-2904-4b1e-a76d-aa04515eb966'::uuid
  ),
  (
    'cad1e583-f48c-460f-8a46-a3a86abbb2fa'::uuid,
    '51868d15-bed5-4485-bc6d-9627b6285236',
    'Uber *Uber *Trip',
    722,
    false,
    '2024-11-01 20:24:30.001-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    'b06e0c42-4053-4fad-b289-be0cfc22502c'::uuid
  ),
  (
    '18d326f3-13f2-43c3-ab33-920bc9caefb2'::uuid,
    '7977be03-7d6c-4982-9802-91841dfd2b30',
    'Transferência enviada|RECEITA FEDERAL',
    -125359,
    false,
    '2024-11-14 07:49:21.156-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    'c198bbd9-3f06-42cd-b2be-b212404d83fc'::uuid
  ),
  (
    '3319a062-50cd-4ea3-afbe-8edc18b21686'::uuid,
    '9b2a203d-309b-42fb-bcdd-7d063a5cf94d',
    'IOF de compra internacional',
    150,
    false,
    '2024-07-06 23:09:27.153-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    'd62f49b0-7aa6-4346-a7e4-bb156b0a99d4'::uuid
  ),
  (
    'f8309ff6-f457-485e-abd5-6c8df4f20ceb'::uuid,
    'd529515e-dbf3-4bc4-8bf9-1f2ab4b02a93',
    'Transferência Recebida|DANIEL SANTOS DE MESQUITA',
    27000,
    false,
    '2024-12-05 12:00:56.022-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    '2158b0b6-844f-44b6-b487-282d0c1b045c'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    '9a86591f-506c-4c54-8abb-3496b30ea57d'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    'f1f7c2ee-f797-4014-ae7e-ab40bee5afdd'::uuid
  ),
  (
    '37204747-eedf-4407-8618-ce0e24c9a36a'::uuid,
    '246471c4-4521-472d-b4ec-0f99a3c063e5',
    'Lauton *Lauton.Com.Br 3/6',
    3993,
    false,
    '2024-08-08 00:00:00-03',
    '2025-03-09 14:11:20.418619-03',
    '2025-03-09 14:11:20.418619-03',
    NULL,
    'abbedc1f-0812-4ed1-9ec9-f51ca13e1069'::uuid,
    'dd696788-412f-48e6-bdfb-f0db7c497bf6'::uuid,
    'df25c07b-2db4-407c-a3b6-f8b1406b7a58'::uuid,
    'ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e'::uuid,
    NULL
  );