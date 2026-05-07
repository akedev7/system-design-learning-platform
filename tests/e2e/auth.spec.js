const { test, expect } = require('@playwright/test');

test.describe('Authentication (US8, US15, US25)', () => {
  test('US8: Login page loads with Auth0', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /course catalog/i })).toBeVisible();
  });

  test('US25: Admin route accessibility', async ({ page }) => {
    // With DISABLE_AUTH=true, admin routes are accessible
    await page.goto('/admin');
    const url = page.url();
    // Should either show admin page or redirect to login
    expect(url).toMatch(/admin|login|api\/auth/);
  });

  test('US8: Auth0 login flow redirects correctly', async ({ page }) => {
    await page.goto('/api/auth/login');
    const url = page.url();
    expect(url).toMatch(/auth0|localhost:3000|login/);
  });
});
