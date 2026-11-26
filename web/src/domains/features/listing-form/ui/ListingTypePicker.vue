<template>
  <fieldset :disabled="disabled">
    <legend class="mb-2 text-body font-medium">
      {{ t('listing.form.type') }}
    </legend>
    <div class="grid gap-2 sm:grid-cols-3">
      <label
        v-for="type in LISTING_TYPES"
        :key="type"
        class="flex cursor-pointer flex-col gap-0.5 rounded-[3px] border px-3 py-2.5 transition-colors has-[input:focus-visible]:outline-2 has-[input:focus-visible]:outline-moss-500 has-[input:disabled]:cursor-not-allowed"
        :class="modelValue === type
          ? 'border-moss-500 bg-moss-50'
          : 'border-rule bg-sheet hover:border-rule-strong'"
      >
        <input
          type="radio"
          name="listing-type"
          class="sr-only"
          :value="type"
          :checked="modelValue === type"
          @change="emit('update:modelValue', type)"
        >
        <span class="font-semibold" :class="modelValue === type && 'text-moss-700'">
          {{ t(`listing.form.typeAction.${type}`) }}
        </span>
        <span class="text-meta text-graphite">{{ t(`listing.form.typeHint.${type}`) }}</span>
      </label>
    </div>
    <p v-if="disabled" class="mt-2 text-meta text-graphite">
      {{ t('listing.form.typeLocked') }}
    </p>
  </fieldset>
</template>

<script setup lang="ts">
import type { ListingType } from 'domains/entities/listing';
import { LISTING_TYPES } from 'domains/entities/listing';
import { useI18n } from 'vue-i18n';

defineProps<{
  modelValue?: ListingType;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: ListingType];
}>();

const { t } = useI18n();
</script>
