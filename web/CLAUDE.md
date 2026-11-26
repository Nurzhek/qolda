# CLAUDE.md — qolda web

Vue 3 SPA (Vite, no SSR). `@nuxt/ui` 4 is installed as a Vue plugin. Package manager: pnpm
(`npx -y pnpm@10 …` if it isn't installed). There is no test runner.

```
pnpm dev          # http://localhost:5173, /api and /files are proxied to VITE_API_TARGET (default :8090)
pnpm lint         # eslint (@antfu); pnpm format applies autofixes
pnpm type-check   # vue-tsc --noEmit
pnpm build        # dist/, served by the backend Go binary
```

`lint` and `type-check` must be clean. Don't claim "tests passed" (there are none); say what you checked by hand.

## Layers (FSD)

```
app -> pages -> widgets -> features -> entities -> shared
```

- `src/app`: composition: plugins, router, guards, reaction to session loss (`setHttpHooks`).
- `src/pages/<page>`: the route (`<page>.router.ts`) plus a page that picks a widget/feature.
- `src/domains/widgets`: composition of features and entities (catalogue, listing details, my listings).
- `src/domains/features`: user scenarios that change data: sign-in, filters, listing form, status change, profile.
- `src/domains/entities`: `dict`, `user`, `listing`: dto, repository, store, read service, simple UI.
- `src/shared`: http, Executor, utils, shared stores (auth, locale), base UI states.

Imports go only downward, no horizontal imports within a layer, slices are imported through `index.ts`.
`shared` knows nothing about `domains` or `app`.

## Data: repository → service → store → component

- URLs live only in `api/*.repository.ts` (the exception is refresh in `shared/api/http.ts`).
- Requests go through `Executor.run` with a `loadingState` from `useLoading()`. It returns
  `{ ok: true, data } | { ok: false, error }`; the error toast is already shown (except 401/403 and field errors).
- Stores don't make requests. Reads go through the entity service, changes through the feature service.
- The API returns protojson: snake_case, int64 values (id, total_count) arrive **as strings**, so ids are typed `string`.
- Every data view renders `loading → error (with retry) → empty → success`.
- Forms: vee-validate + yup (`toTypedSchema`); `fields` from an `invalid_request` response go to fields via `setErrors`.

## UI

- Tokens in `src/app/assets/main.css`: `paper`, `sheet`, `ink`, `graphite`, `rule`, accent `moss-*`;
  fonts Onest (UI) and JetBrains Mono (prices and dates only). Hardcoded hex in components doesn't pass review.
- 2–4px radii, hairline borders instead of shadows, one accent. No gradients, emoji or all-caps labels.
- Icons: Phosphor (`Ph*`, auto-imported). Slice components are imported explicitly.
- User-facing strings only through i18n, in `src/i18n/en/` (English is the only locale for now; the
  language switcher shows up automatically once `LOCALE_LIST` has a second entry).
- Route names: `PageName` from `shared/lib/types/app/pages.ts`.
- `<script setup lang="ts">`, block order `template → script → style`, no `any`.
- Code comments are in English; explain *why*, not *what*.
