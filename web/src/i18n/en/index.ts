import { getLocaleModules } from 'i18n/i18n-utils';

const messages = import.meta.glob('./*.ts', {
  eager: true,
});

export default {
  ...getLocaleModules(messages),
};
