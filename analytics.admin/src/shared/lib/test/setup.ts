import '@testing-library/jest-dom/vitest';
import { i18n } from '@/shared/i18n/config';
import en from '@/shared/i18n/locales/en.json';

i18n.addResourceBundle('en', 'translation', en, true, true);
i18n.options.react = { ...i18n.options.react, useSuspense: false };
void i18n.changeLanguage('en');
