<template>
  <div class="flex flex-col items-start gap-2 border-t border-rule py-8">
    <h2 class="text-title font-semibold">
      {{ title }}
    </h2>
    <p v-if="hint" class="max-w-prose text-body text-graphite">
      {{ hint }}
    </p>
    <div v-if="$slots.default" class="pt-2">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  hint?: string;
}>();
</script>
