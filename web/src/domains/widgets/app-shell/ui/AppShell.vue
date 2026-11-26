<template>
  <div class="flex min-h-dvh flex-col">
    <a
      href="#main"
      class="sr-only z-50 rounded-[3px] bg-ink px-3 py-2 text-paper focus:not-sr-only focus:fixed focus:top-2 focus:left-2"
    >
      {{ t('common.nav.skip') }}
    </a>
    <AppHeader />
    <main id="main" class="mx-auto w-full max-w-7xl flex-1 px-4 pt-6 pb-20 sm:px-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import AppHeader from './AppHeader.vue';

const { t } = useI18n();
</script>
