<template>
  <form
    role="search"
    class="flex-1"
    @submit.prevent="commit"
  >
    <label for="listing-search" class="sr-only">{{ t('listing.filter.searchLabel') }}</label>
    <UInput
      id="listing-search"
      v-model="text"
      type="search"
      size="lg"
      :placeholder="t('listing.filter.search')"
      class="w-full"
      autocomplete="off"
    >
      <template #leading>
        <PhMagnifyingGlass :size="18" aria-hidden="true" />
      </template>
    </UInput>
  </form>
</template>

<script setup lang="ts">
import { watchDebounced } from '@vueuse/core';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const SEARCH_DEBOUNCE_MS = 400;

const { t } = useI18n();
const text = ref(props.modelValue);

// Back/forward in history changes the query: sync it into the field
watch(() => props.modelValue, (value) => {
  if (value !== text.value.trim())
    text.value = value;
});

function commit() {
  if (text.value.trim() !== props.modelValue)
    emit('update:modelValue', text.value.trim());
}

watchDebounced(text, commit, { debounce: SEARCH_DEBOUNCE_MS });
</script>
