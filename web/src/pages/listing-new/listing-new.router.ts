import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'item/new',
    name: PageName.LISTING_NEW,
    meta: { title: 'listing.form.titleNew' },
    component: () => import('./ui/ListingNewPage.vue'),
  });
};
