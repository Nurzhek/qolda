import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'profile',
    name: PageName.PROFILE,
    meta: { title: 'common.nav.profile' },
    component: () => import('./ui/ProfilePage.vue'),
  });
};
