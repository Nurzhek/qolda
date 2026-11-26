import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'my',
    name: PageName.MY_LISTINGS,
    meta: { title: 'common.nav.myListings' },
    component: () => import('./ui/MyListingsPage.vue'),
  });
};
