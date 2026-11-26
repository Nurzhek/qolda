<template>
  <span
    class="inline-flex items-center rounded-[2px] border border-rule bg-sheet px-1.5 py-px text-meta font-medium text-ink"
  >
    {{ t(`listing.type.${type}`) }}
  </span>
</template>

<script setup lang="ts">
import type { ListingType } from '../model/listing.dto';
import { useI18n } from 'vue-i18n';

defineProps<{
  type: ListingType;
}>();

const { t } = useI18n();
</script>
