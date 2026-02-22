// i18n configuration — supported locales, default locale, and type-safe helpers.

export const defaultLocale = 'en' as const;
export const locales = ['en', 'es'] as const;
export type Locale = (typeof locales)[number];

/** Map locale codes to their display names (in their own language). */
export const localeNames: Record<Locale, string> = {
  en: 'English',
  es: 'Español',
};

/** Check if a string is a valid locale. */
export function isLocale(value: string): value is Locale {
  return (locales as readonly string[]).includes(value);
}
