import ui from '@nuxt/ui/vue-plugin';
import { useUserStoreOutside } from 'domains/entities/user';
import { setupI18n } from 'i18n/index';
import { setHttpHooks } from 'shared/api/http';
import { PageName } from 'shared/lib/types/app/pages';
import { createApp } from 'vue';

import App from './App.vue';
import { registerPages } from './providers/register-pages';
import { setupRouterGuards } from './providers/router/guards';
import { router } from './providers/router/router';
import { store } from './providers/store';

import './assets/main.css';

export async function initialize() {
  const app = createApp(App);
  app.use(store);
  registerPages(router);
  await setupI18n(app);
  app.use(ui);
  app.use(router);
  setupRouterGuards(router);

  // The transport knows nothing about the router or the user store: the app layer decides how to react to a lost session
  setHttpHooks({
    onSessionLost() {
      useUserStoreOutside().reset();
      router.replace({ name: PageName.SIGN_IN, query: { redirect: router.currentRoute.value.fullPath } });
    },
    onForbidden() {
      router.replace({ name: PageName.ERROR_403 });
    },
  });

  app.mount('#app');
}
