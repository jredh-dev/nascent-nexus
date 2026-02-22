import { test, expect } from '@playwright/test';

// Full browser-driven login cycle test.
// Uses the seeded demo account (demo@demo.com / demo).
//
// Run:
//   BASE_URL=https://portal.jredh.com npx playwright test
//   BASE_URL=http://localhost:4321 npx playwright test

test.describe('Login cycle', () => {
  test('shows login page', async ({ page }) => {
    await page.goto('/en/login');
    await expect(page.locator('h1')).toContainText('Welcome Back');
    await expect(page.locator('form')).toBeVisible();
  });

  test('bad credentials show error on login page', async ({ page }) => {
    await page.goto('/en/login');

    await page.fill('input[name="email"]', 'nobody@example.com');
    await page.fill('input[name="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');

    // Should stay on /en/login with an error message.
    await page.waitForURL(/\/en\/login/);
    await expect(page.locator('.notification.is-danger')).toBeVisible();
    await expect(page.locator('.notification.is-danger')).toContainText('Invalid');
  });

  test('valid credentials redirect to dashboard with real data', async ({ page }) => {
    await page.goto('/en/login');

    await page.fill('input[name="email"]', 'demo@demo.com');
    await page.fill('input[name="password"]', 'demo');
    await page.click('button[type="submit"]');

    // Should redirect to /en/dashboard.
    await page.waitForURL(/\/en\/dashboard/, { timeout: 15_000 });
    await expect(page.locator('h1')).toContainText('Dashboard');

    // Should show the demo user's profile data.
    await expect(page.locator('text=demo@demo.com')).toBeVisible();

    // Should have at least one active session.
    await expect(page.locator('table').first()).toBeVisible();
  });

  test('logout redirects to home', async ({ page }) => {
    // Login first.
    await page.goto('/en/login');
    await page.fill('input[name="email"]', 'demo@demo.com');
    await page.fill('input[name="password"]', 'demo');
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/en\/dashboard/, { timeout: 15_000 });

    // Click logout.
    await page.click('a[href="/logout"]');

    // Should redirect to / which redirects to /en/.
    await page.waitForURL(/\/en\/?$/);
  });

  test('bare /login redirects to /en/login', async ({ page }) => {
    await page.goto('/login');
    await page.waitForURL(/\/en\/login/);
    await expect(page.locator('h1')).toContainText('Welcome Back');
  });
});

test.describe('Signup page', () => {
  test('shows signup page', async ({ page }) => {
    await page.goto('/en/signup');
    await expect(page.locator('h1')).toContainText('Create Account');
    await expect(page.locator('form')).toBeVisible();
  });

  test('missing fields show error', async ({ page }) => {
    await page.goto('/en/signup');

    // Fill all required fields but use an invalid/duplicate email to trigger
    // a server-side error from the portal backend.
    await page.fill('input[name="username"]', 'testuser_' + Date.now());
    await page.fill('input[name="email"]', 'demo@demo.com'); // duplicate
    await page.fill('input[name="phone"]', '5551234567');
    await page.fill('input[name="password"]', 'password123');

    await page.click('button[type="submit"]');

    // Should redirect back to /en/signup with an error.
    await page.waitForURL(/\/en\/signup/, { timeout: 10_000 });
    await expect(page.locator('.notification.is-danger')).toBeVisible();
  });
});

test.describe('i18n', () => {
  test('Spanish login page renders translated content', async ({ page }) => {
    await page.goto('/es/login');
    await expect(page.locator('h1')).toContainText('Bienvenido de nuevo');
    await expect(page.locator('button[type="submit"]')).toContainText('Iniciar sesión');
  });

  test('Spanish signup page renders translated content', async ({ page }) => {
    await page.goto('/es/signup');
    await expect(page.locator('h1')).toContainText('Crear cuenta');
  });

  test('language switcher navigates between locales', async ({ page }) => {
    await page.goto('/en/about');
    await expect(page.locator('h1')).toContainText('Jared Hooper');

    // Click the ES language switcher link
    await page.click('.lang-switcher-link');
    await page.waitForURL(/\/es\/about/);

    // Should show Spanish content
    await expect(page.locator('.subtitle.is-6.is-muted').first()).toContainText('Ingeniero de Software');
  });

  test('dashboard timestamps are formatted client-side', async ({ page }) => {
    // Login first.
    await page.goto('/en/login');
    await page.fill('input[name="email"]', 'demo@demo.com');
    await page.fill('input[name="password"]', 'demo');
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/en\/dashboard/, { timeout: 15_000 });

    // Wait for client-side script to run.
    await page.waitForTimeout(500);

    // time elements should have been reformatted from ISO strings.
    // Check that at least one <time> element has a human-readable format.
    const timeEl = page.locator('time[data-iso]').first();
    await expect(timeEl).toBeVisible();
    const text = await timeEl.textContent();
    // Should no longer be a raw ISO string (those contain 'T' and 'Z').
    expect(text).not.toContain('T');
  });
});
