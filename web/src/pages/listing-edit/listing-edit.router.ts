import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute(PageName.SHELL, {
    path: 'item/:id/edit',
    name: PageName.LISTING_EDIT,
    meta: { title: 'listing.form.titleEdit' },
    props: true,
    component: () => import('./ui/ListingEditPage.vue'),
  });
};
