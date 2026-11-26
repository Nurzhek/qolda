<template>
  <ul class="grid grid-cols-2 gap-x-4 gap-y-7 sm:grid-cols-3 xl:grid-cols-4">
    <template v-if="loading">
      <li
        v-for="i in skeletons"
        :key="i"
        class="flex flex-col gap-2"
        aria-hidden="true"
      >
        <USkeleton class="aspect-[4/3] w-full rounded-[3px]" />
        <USkeleton class="h-5 w-1/2" />
        <USkeleton class="h-4 w-4/5" />
      </li>
    </template>
    <template v-else>
      <li
        v-for="listing in items"
        :key="listing.id"
        class="min-w-0"
      >
        <slot :listing="listing" />
      </li>
    </template>
  </ul>
</template>

<script setup lang="ts" generic="T extends { id: string }">
withDefaults(defineProps<{
  items?: T[];
  loading?: boolean;
  skeletons?: number;
}>(), {
  items: () => [],
  loading: false,
  skeletons: 8,
});

defineSlots<{
  default: (props: { listing: T }) => unknown;
}>();
</script>
