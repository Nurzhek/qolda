<template>
  <div class="flex flex-col gap-2">
    <h2 class="text-title font-semibold">
      {{ t('listing.detail.description') }}
    </h2>
    <p v-if="text" class="max-w-[68ch] text-lead whitespace-pre-line">
      {{ text }}
    </p>
    <p v-else class="text-body text-graphite">
      {{ t('listing.detail.noDescription') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

defineProps<{
  text: string;
}>();

const { t } = useI18n();
</script>
