import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute({
    path: '/sign-in',
    name: PageName.SIGN_IN,
    meta: { public: true, onlyLoggedOut: true, title: 'auth.dev.submit' },
    component: () => import('./ui/SignInPage.vue'),
  });
};
