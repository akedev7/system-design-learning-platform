const { test, expect } = require('@playwright/test');

test.describe('Course Detail Page', () => {
  test('should display course details and modules', async ({ page }) => {
    // Go to course catalog
    await page.goto('http://localhost:3000');
    
    // Click on first course
    await page.waitForSelector('a[href*="/courses"]');
    await page.click('a[href*="/courses"]:first-of-type');
    
    // Wait for course page to load
    await page.waitForURL('**/courses/**');
    
    // Check course title is visible
    await expect(page.locator('h1')).toContainText('System Design Fundamentals');
    
    // Wait for modules to load
    await page.waitForSelector('a[href*="/modules"]', { timeout: 10000 });
    
    // Check modules are displayed
    await expect(page.locator('text=Introduction to System Design')).toBeVisible();
    await expect(page.locator('text=Client-Server Architecture')).toBeVisible();
  });

  test('should navigate to module page', async ({ page }) => {
    // Go to course 2
    await page.goto('http://localhost:3000/courses/2');
    
    // Wait for modules to load
    await page.waitForSelector('a[href*="/modules"]', { timeout: 10000 });
    
    // Click on first module
    await page.click('a[href*="/modules"]:first-of-type');
    
    // Should navigate to module page
    await page.waitForURL('**/modules/**');
    
    // Check module title is visible
    await expect(page.locator('h1')).toContainText('Introduction to System Design');
  });

  test('should display lesson in module', async ({ page }) => {
    // Go to module 2 (Introduction to System Design)
    await page.goto('http://localhost:3000/courses/2/modules/2');
    
    // Wait for lessons to load
    await page.waitForSelector('a[href*="/lessons"]', { timeout: 10000 });
    
    // Check lessons are displayed
    await expect(page.locator('text=What is System Design?')).toBeVisible();
    await expect(page.locator('text=Key Components')).toBeVisible();
  });
});
