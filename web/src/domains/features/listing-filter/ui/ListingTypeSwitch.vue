<template>
  <div
    class="flex gap-5 border-b border-rule"
    role="group"
    :aria-label="t('listing.filter.title')"
  >
    <button
      v-for="option in options"
      :key="option.value ?? 'all'"
      type="button"
      class="-mb-px border-b-2 pt-1 pb-2 text-body font-medium transition-colors"
      :class="option.value === modelValue
        ? 'border-moss-500 text-ink'
        : 'border-transparent text-graphite hover:text-ink'"
      :aria-pressed="option.value === modelValue"
      @click="emit('update:modelValue', option.value)"
    >
      {{ option.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
import type { ListingType } from 'domains/entities/listing';
import { LISTING_TYPES } from 'domains/entities/listing';
import { useI18n } from 'vue-i18n';

defineProps<{
  modelValue?: ListingType;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: ListingType | undefined];
}>();

const { t } = useI18n();

const options = computed(() => [
  { value: undefined, label: t('listing.filter.allTypes') },
  ...LISTING_TYPES.map(value => ({ value, label: t(`listing.type.${value}`) })),
]);
</script>
