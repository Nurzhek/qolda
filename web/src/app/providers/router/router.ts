import { PageName } from 'shared/lib/types/app/pages';
import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, saved) {
    if (saved)
      return saved;
    // Catalogue filter changes only touch the query: don't scroll the page
    if (to.name === from.name && to.path === from.path)
      return false;
    return { top: 0 };
  },
  routes: [
    {
      path: '/',
      name: PageName.SHELL,
      component: () => import('domains/widgets/app-shell').then(m => m.AppShell),
      children: [],
    },
  ],
});
