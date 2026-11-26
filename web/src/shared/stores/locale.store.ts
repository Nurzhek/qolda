import type { Locale } from 'shared/config/locale-config';
import { store } from 'app/providers/store';
import { loadLanguageAsync } from 'i18n/index';
import { defineStore } from 'pinia';
import { DEFAULT_LOCALE, isLocale } from 'shared/config/locale-config';

const STORE_ID = 'locale';

function preferredLocale(): Locale {
  const match = window.navigator.languages
    .map(lang => lang.split('-')[0])
    .find(isLocale);
  return match ?? DEFAULT_LOCALE;
}

export const useLocaleStore = defineStore(
  STORE_ID,
  () => {
    const locale = ref<Locale>(preferredLocale());

    async function setLocale(value: Locale) {
      locale.value = await loadLanguageAsync(value);
      return locale.value;
    }

    async function initLocale() {
      await setLocale(isLocale(locale.value) ? locale.value : DEFAULT_LOCALE);
    }

    return {
      locale,
      setLocale,
      initLocale,
    };
  },
  {
    persist: {
      key: 'qolda-locale',
    },
  },
);

// Need to be used outside the setup
export function useLocaleStoreOutside() {
  return useLocaleStore(store);
}
