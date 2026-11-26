import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'u/:id',
    name: PageName.USER,
    meta: { title: 'common.nav.profile' },
    props: true,
    component: () => import('./ui/UserPage.vue'),
  });
};
