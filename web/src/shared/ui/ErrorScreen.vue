<template>
  <main class="mx-auto flex min-h-dvh w-full max-w-xl flex-col justify-center gap-4 px-6">
    <RouterLink to="/" class="w-fit">
      <AppWordmark />
    </RouterLink>
    <p class="pt-6 font-mono text-display text-graphite tabular" aria-hidden="true">
      {{ code }}
    </p>
    <h1 class="text-heading font-semibold">
      {{ title }}
    </h1>
    <p class="text-lead text-graphite">
      {{ text }}
    </p>
    <div class="pt-2">
      <slot />
    </div>
  </main>
</template>

<script setup lang="ts">
import AppWordmark from './AppWordmark.vue';

defineProps<{
  code: string;
  title: string;
  text: string;
}>();
</script>
