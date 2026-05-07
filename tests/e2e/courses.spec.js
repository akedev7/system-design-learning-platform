const { test, expect } = require('@playwright/test');

test.describe('Course Browsing & Enrollment (US1, US2, US17, US19)', () => {
  test('US1: Browse available courses on homepage', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /course catalog/i })).toBeVisible();
    // Check course cards are displayed
    const courseCards = page.getByTestId('course-card');
    await expect(courseCards.first()).toBeVisible();
  });

  test('US19: Navigate hierarchical routes to course', async ({ page }) => {
    await page.goto('/');
    await page.getByTestId('course-card').first().click();
    await expect(page).toHaveURL(/\/courses\/[^\/]+$/);
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible();
  });

  test('US2: Enroll in a course', async ({ page }) => {
    await page.goto('/');
    const firstCourse = page.getByTestId('course-card').first();
    await firstCourse.click();

    const enrollButton = page.getByRole('button', { name: /enroll|start course/i });
    await expect(enrollButton).toBeVisible();
    await enrollButton.click();

    // Should show enrolled state
    await expect(page.getByText(/enrolled|progress/i)).toBeVisible();
  });

  test('US17: View course completion percentage', async ({ page }) => {
    await page.goto('/');
    const firstCourse = page.getByTestId('course-card').first();
    await firstCourse.click();

    // Check for completion percentage indicator
    const progressBar = page.getByRole('progressbar').or(page.getByTestId('completion-percentage'));
    await expect(progressBar).toBeVisible();
  });

  test('US19: Navigate to module from course', async ({ page }) => {
    await page.goto('/');
    await page.getByTestId('course-card').first().click();
    await page.getByTestId('module-card').first().click();
    await expect(page).toHaveURL(/\/modules\/[^\/]+$/);
  });

  test('US19: Navigate to lesson from module', async ({ page }) => {
    await page.goto('/');
    await page.getByTestId('course-card').first().click();
    await page.getByTestId('module-card').first().click();
    await page.getByTestId('lesson-card').first().click();
    await expect(page).toHaveURL(/\/lessons\/[^\/]+$/);
  });
});
