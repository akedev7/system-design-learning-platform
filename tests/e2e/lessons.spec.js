const { test, expect } = require('@playwright/test');

test.describe('Lesson Content & Progress (US3, US4, US5, US6, US7, US16)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    const firstCourse = page.getByTestId('course-card').first();
    await firstCourse.click();
    await page.getByTestId('module-card').first().click();
    await page.getByTestId('lesson-card').first().click();
  });

  test('US3: View lesson with mixed content blocks', async ({ page }) => {
    // Check for various content block types
    await expect(page.getByTestId('text-block')).toBeVisible();
    await expect(page.getByTestId('quiz-block')).toBeVisible();
    await expect(page.getByTestId('diagram-block')).toBeVisible();
  });

  test('US4: Take multiple choice quiz', async ({ page }) => {
    const quizBlock = page.getByTestId('quiz-block').first();
    await quizBlock.getByRole('radio').first().click();
    await quizBlock.getByRole('button', { name: /submit|check/i }).click();

    // Should show feedback
    await expect(quizBlock.getByText(/correct|incorrect|score/i)).toBeVisible();
  });

  test('US4: Take true/false quiz', async ({ page }) => {
    const quizBlock = page.getByTestId('quiz-block').nth(1);
    await quizBlock.getByRole('radio', { name: /true|false/i }).first().click();
    await quizBlock.getByRole('button', { name: /submit|check/i }).click();

    await expect(quizBlock.getByText(/correct|incorrect|score/i)).toBeVisible();
  });

  test('US5: Build diagram with React Flow drag-and-drop', async ({ page }) => {
    const diagramBlock = page.getByTestId('diagram-block').first();

    // Drag a node from palette to canvas
    const paletteNode = diagramBlock.getByTestId('palette-node').first();
    const canvas = diagramBlock.getByTestId('react-flow-canvas');

    await paletteNode.dragTo(canvas);

    // Node should appear on canvas
    await expect(diagramBlock.getByTestId('flow-node')).toBeVisible();
  });

  test('US6: Get instant feedback on diagram', async ({ page }) => {
    const diagramBlock = page.getByTestId('diagram-block').first();

    // Build a correct diagram (simplified - add required nodes)
    const paletteNode = diagramBlock.getByTestId('palette-node').first();
    const canvas = diagramBlock.getByTestId('react-flow-canvas');
    await paletteNode.dragTo(canvas);

    await diagramBlock.getByRole('button', { name: /check|submit|validate/i }).click();

    // Should show instant feedback
    await expect(diagramBlock.getByText(/correct|incorrect|feedback/i)).toBeVisible();
  });

  test('US16: Client-side diagram validation', async ({ page }) => {
    const diagramBlock = page.getByTestId('diagram-block').first();

    // Submit empty diagram - should show client-side validation
    await diagramBlock.getByRole('button', { name: /submit|check/i }).click();

    // Should show validation message without server call
    await expect(diagramBlock.getByText(/incomplete|add nodes|required/i)).toBeVisible();
  });

  test('US7: Track progress at lesson level', async ({ page }) => {
    // Complete a quiz to trigger progress
    const quizBlock = page.getByTestId('quiz-block').first();
    await quizBlock.getByRole('radio').first().click();
    await quizBlock.getByRole('button', { name: /submit/i }).click();

    // Check progress indicator updates
    const progress = page.getByTestId('lesson-progress');
    await expect(progress).toBeVisible();
  });
});
