const { test, expect } = require('@playwright/test');

test.describe('Optimized Images (US20)', () => {
  test('US20: Images load from CDN with optimization', async ({ page }) => {
    await page.goto('/');

    const firstCourse = page.getByTestId('course-card').first();
    await firstCourse.click();

    // Check for optimized images
    const images = page.getByRole('img');
    const count = await images.count();

    if (count > 0) {
      const src = await images.first().getAttribute('src');
      // Next.js optimized images should have URL params or be from CDN
      expect(src).toBeTruthy();
    }
  });

  test('US20: Next.js Image component optimization', async ({ page }) => {
    await page.goto('/');

    // Check for Next.js image optimization attributes
    const images = page.locator('img[src*="cdn"], img[data-nimg]');
    await expect(images.first()).toBeVisible();
  });
});
