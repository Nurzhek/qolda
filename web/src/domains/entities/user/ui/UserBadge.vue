<template>
  <div class="flex items-center gap-3">
    <span
      class="grid size-10 shrink-0 place-items-center rounded-sm bg-paper-deep text-title font-semibold text-ink"
      aria-hidden="true"
    >
      {{ initial }}
    </span>
    <div class="min-w-0">
      <RouterLink
        v-if="to"
        :to="to"
        class="block truncate font-medium underline-offset-2 hover:underline"
      >
        {{ name }}
      </RouterLink>
      <p v-else class="truncate font-medium">
        {{ name }}
      </p>
      <p v-if="meta" class="truncate text-meta text-graphite">
        {{ meta }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { RouteLocationRaw } from 'vue-router';

const props = defineProps<{
  name: string;
  meta?: string;
  to?: RouteLocationRaw;
}>();

const initial = computed(() => props.name.trim().charAt(0).toUpperCase() || '?');
</script>
