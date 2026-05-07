const { test, expect } = require('@playwright/test');

test.describe('System Design Learning Platform', () => {
  test('should load course catalog', async ({ page }) => {
    await page.goto('http://localhost:3000');
    
    // Check title
    await expect(page).toHaveTitle('Course Catalog');
  });

  test('should display seeded courses', async ({ page }) => {
    await page.goto('http://localhost:3000');
    
    // Wait for courses to load
    await page.waitForSelector('a[href*="/courses"]', { timeout: 10000 });
    
    // Check for seeded courses
    await expect(page.locator('text=System Design Fundamentals')).toBeVisible();
    await expect(page.locator('text=Scalability & Load Balancing')).toBeVisible();
    await expect(page.locator('text=Database Choices')).toBeVisible();
  });

  test('should navigate to course details', async ({ page }) => {
    await page.goto('http://localhost:3000');
    
    // Wait for courses to load
    await page.waitForSelector('a[href*="/courses"]');
    
    // Click on first course
    await page.click('a[href*="/courses"]:first-of-type');
    
    // Should navigate to course page
    await page.waitForURL('**/courses/**');
    await expect(page).toHaveURL(/.*\/courses\/\d+/);
  });

  test('API should have CORS headers', async ({ request }) => {
    const response = await request.get('http://localhost:8080/api/v1/courses', {
      headers: {
        'Origin': 'http://localhost:3000'
      }
    });
    
    // Check CORS header (case-insensitive check)
    const headers = response.headers();
    const corsHeader = headers['access-control-allow-origin'] || headers['Access-Control-Allow-Origin'];
    expect(corsHeader).toBe('http://localhost:3000');
    
    // Check response structure
    const body = await response.json();
    expect(body.status).toBe('success');
    expect(Array.isArray(body.data)).toBeTruthy();
  });
});
