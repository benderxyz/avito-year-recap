export const SUPPORTED_LOCALES = ['en', 'ru'] as const;

export type SupportedLocale = (typeof SUPPORTED_LOCALES)[number];

export const DEFAULT_LOCALE: SupportedLocale = 'en';

export const LOCALES_LOAD_PATH =
  import.meta.env.VITE_LOCALES_URL ?? '/locales/{{lng}}.json';
