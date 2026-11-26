import type { Locale } from 'shared/config/locale-config';
import type { DefaultModule } from 'shared/lib/types/app/app';
import type { App } from 'vue';

import { AVAILABLE_LOCALES, DEFAULT_LOCALE } from 'shared/config/locale-config';
import { useLocaleStoreOutside } from 'shared/stores/locale.store';
import { createI18n } from 'vue-i18n';

const loadedLanguages: Locale[] = [];

export const i18nApp = createI18n({
  // you must set `false`, to use Composition API
  legacy: false,
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  availableLocales: AVAILABLE_LOCALES,
  missingWarn: false,
  fallbackWarn: false,
});

const { global: i18n } = i18nApp;
export { i18n };

export async function setupI18n(app: App) {
  app.use(i18nApp);

  // English is the fallback locale: always loaded so te()/fallback work
  await loadLanguageAsync(DEFAULT_LOCALE);
  await useLocaleStoreOutside().initLocale();
}

function setLanguage(lang: Locale) {
  i18n.locale.value = lang;
  document.documentElement.setAttribute('lang', lang);
  return lang;
}

export async function loadLanguageAsync(lang: Locale) {
  if (!loadedLanguages.includes(lang)) {
    const messages: DefaultModule = await import(`./${lang}/index.ts`);
    i18n.setLocaleMessage(lang, messages.default);
    loadedLanguages.push(lang);
  }

  return setLanguage(lang);
}
