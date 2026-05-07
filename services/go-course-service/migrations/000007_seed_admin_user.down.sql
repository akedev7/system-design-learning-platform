-- Remove admin seed user (if it was the seeded one)
DELETE FROM users 
WHERE auth0_id = '00000000-0000-0000-0000-000000000001'::uuid
  AND email = 'admin@example.com';