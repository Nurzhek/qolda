<template>
  <div
    role="alert"
    class="flex flex-col items-start gap-3 rounded-sm border border-rule bg-sheet p-4"
  >
    <p class="text-body text-graphite">
      {{ message ?? t('common.states.loadError') }}
    </p>
    <UButton
      color="neutral"
      variant="outline"
      :label="t('common.actions.retry')"
      @click="emit('retry')"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

defineProps<{
  message?: string;
}>();

const emit = defineEmits<{
  retry: [];
}>();

const { t } = useI18n();
</script>
