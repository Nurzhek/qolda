import type { DefaultModule } from 'shared/lib/types/app/app';
import type { Router } from 'vue-router';

interface AppPage {
  router?: (router: Router) => void;
}

// Pages layer: each slice exposes its route through its index.ts.
// Don't search deeper than one level: an index.ts inside a slice's ui/model is not a route.
const pages = import.meta.glob<DefaultModule<AppPage>>('pages/*/index.ts', {
  eager: true,
});

export function registerPages(router: Router) {
  Object.keys(pages)
    .sort()
    .forEach((pageKey) => {
      const page = pages[pageKey];

      if (page.default && page.default.router)
        page.default.router(router);
    });
}
