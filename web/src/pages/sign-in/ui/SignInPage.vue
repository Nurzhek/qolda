<template>
  <main class="mx-auto grid min-h-dvh w-full max-w-5xl content-center gap-12 px-6 py-12 md:grid-cols-[minmax(0,1fr)_22rem] md:gap-16">
    <div class="flex flex-col gap-6">
      <AppWordmark />
      <h1 class="max-w-[16ch] text-display font-semibold tracking-[-0.025em] text-balance md:text-[2.75rem] md:leading-[1.05]">
        {{ t('auth.title') }}
      </h1>
      <p class="max-w-[52ch] text-lead text-graphite">
        {{ t('auth.lead') }}
      </p>
    </div>

    <div class="self-center rounded-[4px] border border-rule bg-sheet p-6">
      <div v-if="isLoading" class="flex flex-col gap-3">
        <USkeleton class="h-4 w-full" />
        <USkeleton class="h-11 w-full" />
      </div>
      <StateError v-else-if="hasError || !dict" @retry="fetchDict" />
      <SignInPanel
        v-else
        :redirect="redirect"
        :domains="dict.allowed_domains"
        :google-enabled="dict.google_login"
        :dev-enabled="dict.dev_login"
      />
    </div>
  </main>
</template>

<script setup lang="ts">
import { useDictService } from 'domains/entities/dict';
import { safeRedirect, SignInPanel } from 'domains/features/auth';
import { AppWordmark, StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const route = useRoute();
const { dict, hasError, isLoading, fetchDict } = useDictService();

const redirect = computed(() => safeRedirect(route.query.redirect));
</script>
