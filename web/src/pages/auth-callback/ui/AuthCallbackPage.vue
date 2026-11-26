<template>
  <main class="mx-auto flex min-h-dvh w-full max-w-md flex-col justify-center gap-4 px-6">
    <AppWordmark />
    <template v-if="error">
      <h1 class="pt-4 text-heading font-semibold">
        {{ t('auth.callback.failed') }}
      </h1>
      <p class="text-lead text-graphite" role="alert">
        {{ t(getErrorMessage(error)) }}
      </p>
      <div class="pt-2">
        <UButton :to="{ name: PageName.SIGN_IN }" :label="t('auth.callback.backToSignIn')" />
      </div>
    </template>
    <p
      v-else
      class="pt-4 text-lead text-graphite"
      aria-live="polite"
    >
      {{ t('auth.callback.signingIn') }}
    </p>
  </main>
</template>

<script setup lang="ts">
import { parseCallbackHash, useAuth } from 'domains/features/auth';
import { getErrorMessage } from 'shared/api/error-handler';
import { PageName } from 'shared/lib/types/app/pages';
import { AppWordmark } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { finishSignIn } = useAuth();
const error = shallowRef<string>();

onMounted(async () => {
  const result = parseCallbackHash(window.location.hash);
  // Tokens must not stay in the address bar or history
  history.replaceState(null, '', window.location.pathname);

  if (!result.ok) {
    error.value = result.error;
    return;
  }
  if (!(await finishSignIn(result.access, result.refresh, result.redirect)))
    error.value = 'service_not_available';
});
</script>
