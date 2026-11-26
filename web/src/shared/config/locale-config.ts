// ponytail: English only for now; add a locale here and an src/i18n/<code>/ folder to bring one back
export type Locale = 'en';

export const DEFAULT_LOCALE: Locale = 'en';

export const AVAILABLE_LOCALES: Locale[] = ['en'];

export const LOCALE_LIST: { key: Locale; title: string }[] = [
  { key: 'en', title: 'English' },
];

export function isLocale(value: string): value is Locale {
  return (AVAILABLE_LOCALES as string[]).includes(value);
}
