import type { Router } from 'vue-router';
import { PageName } from 'shared/lib/types/app/pages';

export default (router: Router) => {
  router.addRoute({
    path: '/welcome',
    name: PageName.ONBOARDING,
    meta: { title: 'common.nav.profile' },
    component: () => import('./ui/OnboardingPage.vue'),
  });
};
