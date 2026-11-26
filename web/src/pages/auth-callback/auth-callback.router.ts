import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute({
    path: '/auth/callback',
    name: PageName.AUTH_CALLBACK,
    meta: { public: true },
    component: () => import('./ui/AuthCallbackPage.vue'),
  });
};
