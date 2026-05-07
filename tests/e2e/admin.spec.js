const { test, expect } = require('@playwright/test');

test.describe('Admin Dashboard (US9, US10, US11, US12, US13)', () => {
  test.beforeEach(async ({ context }) => {
    // Mock admin session
    await context.addCookies([{
      name: 'appSession',
      value: 'mock-admin-session',
      domain: 'localhost',
      path: '/'
    }]);
  });

  test('US9: Admin can view admin dashboard', async ({ page }) => {
    await page.goto('/admin');
    await expect(page.getByRole('heading', { name: /admin|dashboard/i })).toBeVisible();
    await expect(page.getByRole('link', { name: /courses|content/i })).toBeVisible();
  });

  test('US9: Admin can create a new course', async ({ page }) => {
    await page.goto('/admin');

    await page.getByRole('button', { name: /create course|new course/i }).click();

    await page.getByLabel(/title/i).fill('Test Course');
    await page.getByLabel(/description/i).fill('Test Description');
    await page.getByRole('button', { name: /save|create/i }).click();

    await expect(page.getByText(/success|created|test course/i)).toBeVisible();
  });

  test('US10: Drag-and-drop block editor for lessons', async ({ page }) => {
    await page.goto('/admin/lessons/new');

    // Drag text block to editor
    const textBlockPalette = page.getByTestId('block-palette-text');
    const editorCanvas = page.getByTestId('content-editor');

    await textBlockPalette.dragTo(editorCanvas);

    // Block should appear in editor
    await expect(editorCanvas.getByTestId('text-block-editor')).toBeVisible();
  });

  test('US11: Build quiz block with custom questions', async ({ page }) => {
    await page.goto('/admin/lessons/new');

    // Add quiz block
    await page.getByTestId('block-palette-quiz').dragTo(page.getByTestId('content-editor'));

    // Configure quiz
    const quizEditor = page.getByTestId('quiz-block-editor');
    await quizEditor.getByRole('button', { name: /add question/i }).click();
    await quizEditor.getByLabel(/question/i).fill('What is system design?');
    await quizEditor.getByLabel(/option/i).first().fill('A design pattern');
    await quizEditor.getByRole('radio', { name: /correct/i }).first().click();

    await expect(quizEditor.getByText('What is system design?')).toBeVisible();
  });

  test('US12: Mini React Flow canvas for diagram answers', async ({ page }) => {
    await page.goto('/admin/lessons/new');

    // Add diagram block
    await page.getByTestId('block-palette-diagram').dragTo(page.getByTestId('content-editor'));

    // Use mini React Flow to define answer
    const diagramEditor = page.getByTestId('diagram-block-editor');
    await diagramEditor.getByTestId('mini-react-flow-canvas').waitFor();

    // Drag expected nodes
    const nodePalette = diagramEditor.getByTestId('answer-node-palette');
    const answerCanvas = diagramEditor.getByTestId('answer-canvas');

    await nodePalette.getByTestId('palette-load-balancer').dragTo(answerCanvas);

    await expect(answerCanvas.getByTestId('answer-node')).toBeVisible();
  });

  test('US13: Upload images to cloud storage', async ({ page }) => {
    await page.goto('/admin/lessons/new');

    // Add image block
    await page.getByTestId('block-palette-image').dragTo(page.getByTestId('content-editor'));

    const imageEditor = page.getByTestId('image-block-editor');

    // Upload image
    const fileInput = imageEditor.getByLabel(/upload|image/i);
    await fileInput.setInputFiles({
      name: 'test-image.png',
      mimeType: 'image/png',
      buffer: Buffer.from(new ArrayBuffer(1024))
    });

    // Should show upload progress/success
    await expect(imageEditor.getByText(/uploading|uploaded|success/i)).toBeVisible();
  });
});
