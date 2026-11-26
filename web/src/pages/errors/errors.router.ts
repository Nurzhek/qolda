import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute({
    path: '/forbidden',
    name: PageName.ERROR_403,
    meta: { public: true, title: 'errors.page.forbidden' },
    props: { code: '403' },
    component: () => import('./ui/ErrorPage.vue'),
  });
  router.addRoute({
    path: '/server-error',
    name: PageName.ERROR_500,
    meta: { public: true, title: 'errors.page.server' },
    props: { code: '500' },
    component: () => import('./ui/ErrorPage.vue'),
  });
  router.addRoute({
    path: '/:pathMatch(.*)*',
    name: PageName.ERROR_404,
    meta: { public: true, title: 'errors.page.notFound' },
    props: { code: '404' },
    component: () => import('./ui/ErrorPage.vue'),
  });
};
