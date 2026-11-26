<template>
  <main class="mx-auto flex min-h-dvh w-full max-w-lg flex-col justify-center gap-6 px-6 py-12">
    <AppWordmark />
    <div class="flex flex-col gap-2">
      <h1 class="text-heading font-semibold">
        {{ t('auth.onboarding.title') }}
      </h1>
      <p class="text-body text-graphite">
        {{ t('auth.onboarding.lead') }}
      </p>
    </div>
    <ProfileForm
      v-if="me"
      :user="me"
      :campuses="campuses"
      :submit-label="t('auth.onboarding.submit')"
      @saved="router.replace(redirect)"
    />
  </main>
</template>

<script setup lang="ts">
import { useDictStore } from 'domains/entities/dict';
import { useUserStore } from 'domains/entities/user';
import { safeRedirect } from 'domains/features/auth';
import { ProfileForm } from 'domains/features/profile-edit';
import { storeToRefs } from 'pinia';
import { AppWordmark } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const { me } = storeToRefs(useUserStore());
const { campuses } = storeToRefs(useDictStore());

const redirect = computed(() => safeRedirect(route.query.redirect));
</script>
