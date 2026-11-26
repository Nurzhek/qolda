import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: '',
    name: PageName.CATALOGUE,
    meta: { title: 'common.nav.browse' },
    component: () => import('./ui/CataloguePage.vue'),
  });
};
