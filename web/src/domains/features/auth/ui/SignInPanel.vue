<template>
  <div class="flex flex-col gap-6">
    <p class="text-body text-graphite">
      {{ domains.length ? t('auth.onlyUniversity', { domains: domainList }) : t('auth.onlyUniversityAny') }}
    </p>

    <UButton
      v-if="googleEnabled"
      size="xl"
      block
      :label="t('auth.google')"
      @click="signInWithGoogle(redirect)"
    >
      <template #leading>
        <PhGoogleLogo
          :size="20"
          weight="bold"
          aria-hidden="true"
        />
      </template>
    </UButton>

    <div v-if="devEnabled" :class="googleEnabled && 'border-t border-rule pt-6'">
      <DevSignInForm :redirect="redirect" :domains="domains" />
    </div>

    <p v-if="!googleEnabled && !devEnabled" class="text-body">
      {{ t('auth.noMethods') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useAuth } from '../model/use-auth';
import DevSignInForm from './DevSignInForm.vue';

const props = defineProps<{
  redirect: string;
  domains: string[];
  googleEnabled: boolean;
  devEnabled: boolean;
}>();

const { t } = useI18n();
const { signInWithGoogle } = useAuth();

const domainList = computed(() => props.domains.map(d => `@${d}`).join(', '));
</script>
