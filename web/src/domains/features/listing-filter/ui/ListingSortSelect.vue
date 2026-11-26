<template>
  <USelect
    :model-value="modelValue"
    :items="items"
    :aria-label="t('listing.filter.sort')"
    class="w-48"
    @update:model-value="onUpdate"
  />
</template>

<script setup lang="ts">
import type { SortKey } from '../model/use-listing-filter';
import { useI18n } from 'vue-i18n';
import { SORTS } from '../model/use-listing-filter';

defineProps<{
  modelValue: SortKey;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: SortKey];
}>();

const { t } = useI18n();

function onUpdate(value: unknown) {
  if (typeof value === 'string' && value in SORTS)
    emit('update:modelValue', value as SortKey);
}

const items = computed(() => (Object.keys(SORTS) as SortKey[]).map(value => ({ value, label: t(`listing.sort.${value}`) })));
</script>
