import type { Router } from 'vue-router';
import { useDictService } from 'domains/entities/dict';
import { useUserService, useUserStoreOutside } from 'domains/entities/user';
import { i18n } from 'i18n/index';
import { PageName } from 'shared/lib/types/app/pages';
import { Logger } from 'shared/lib/utils/logger';
import { useAuthStoreOutside } from 'shared/stores/auth.store';

const _logger = new Logger('Router');

function createAuthGuard(router: Router) {
  const authStore = useAuthStoreOutside();
  const userStore = useUserStoreOutside();
  const { ensureDict } = useDictService();
  const { fetchMe } = useUserService();

  router.beforeEach(async (to) => {
    _logger.log(`Going to ${String(to.name)}`);

    // The sign-in screen also needs the dictionary (which sign-in methods are enabled); an error doesn't block navigation
    await ensureDict();

    if (to.meta.public) {
      if (to.meta.onlyLoggedOut && authStore.isLoggedIn)
        return { name: PageName.CATALOGUE, replace: true };
      return;
    }

    if (!authStore.isLoggedIn)
      return { name: PageName.SIGN_IN, query: { redirect: to.fullPath }, replace: true };

    if (!userStore.me) {
      const ok = await fetchMe();
      if (!ok) {
        // A 401 has already logged out and sent to sign-in; otherwise the backend is unavailable
        return authStore.isLoggedIn ? { name: PageName.ERROR_500, replace: true } : false;
      }
    }

    // The catalogue is useless without a campus: one short form first
    if (!userStore.me?.campus && to.name !== PageName.ONBOARDING)
      return { name: PageName.ONBOARDING, query: { redirect: to.fullPath }, replace: true };
  });
}

function createTitleGuard(router: Router) {
  router.afterEach((to) => {
    const title = to.meta.title ? i18n.t(to.meta.title) : '';
    document.title = title ? `${title} — qolda` : `qolda — ${i18n.t('common.tagline')}`;
  });
}

export function setupRouterGuards(router: Router) {
  createAuthGuard(router);
  createTitleGuard(router);
}
