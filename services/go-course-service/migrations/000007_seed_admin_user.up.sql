-- Seed first admin user
-- To configure the first admin user:
-- Option 1: Update these values directly with your Auth0 sub after first login
-- Option 2: Set env var ADMIN_EMAIL before running migrations (see database.RunMigrations)
--
-- Default placeholder values - replace with actual Auth0 ID after setup
INSERT INTO users (auth0_id, email, name, role)
VALUES (
    '00000000-0000-0000-0000-000000000001'::uuid,
    'admin@example.com',
    'Admin User',
    'Admin'
)
ON CONFLICT (auth0_id) DO NOTHING;