import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'item/:id',
    name: PageName.LISTING,
    meta: { title: 'common.nav.browse' },
    props: true,
    component: () => import('./ui/ListingPage.vue'),
  });
};
