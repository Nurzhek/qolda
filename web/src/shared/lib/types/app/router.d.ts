import 'vue-router';

declare module 'vue-router' {
  interface RouteMeta {
    // available without signing in
    public?: boolean;
    // guests only (sign-in screen)
    onlyLoggedOut?: boolean;
    // i18n key for the tab title
    title?: string;
  }
}
