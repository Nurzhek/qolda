<template>
  <div class="relative overflow-hidden rounded-[3px] border border-rule bg-paper-deep">
    <img
      v-if="src"
      :src="src"
      :alt="title"
      loading="lazy"
      decoding="async"
      class="size-full object-cover"
    >
    <!-- No photo: a typographic placeholder, like a handwritten notice -->
    <div
      v-else
      class="flex size-full items-end p-3"
      aria-hidden="true"
    >
      <span class="line-clamp-3 text-[1.25rem] leading-[1.05] sm:text-[1.625rem] font-bold tracking-[-0.03em] break-words text-ink/75">
        {{ title }}
      </span>
    </div>
    <span v-if="!src" class="sr-only">{{ t('listing.noPhoto') }}</span>
    <div v-if="$slots.default" class="absolute top-2 left-2">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

defineProps<{
  src?: string;
  title: string;
}>();

const { t } = useI18n();
</script>
