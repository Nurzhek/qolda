<template>
  <div class="flex flex-col gap-6">
    <fieldset class="flex flex-col gap-1">
      <legend class="mb-2 text-meta font-semibold text-graphite">
        {{ t('listing.filter.category') }}
      </legend>
      <button
        v-for="option in categoryOptions"
        :key="option.value ?? 'any'"
        type="button"
        class="-mx-2 rounded-[3px] px-2 py-1 text-left text-body transition-colors"
        :class="option.value === state.category
          ? 'bg-moss-50 font-medium text-moss-700'
          : 'text-ink hover:bg-paper-deep'"
        :aria-pressed="option.value === state.category"
        @click="emit('update', { category: option.value })"
      >
        {{ option.label }}
      </button>
    </fieldset>

    <UFormField :label="t('listing.filter.campus')">
      <USelect
        :model-value="state.campus ?? ANY"
        :items="campusOptions"
        class="w-full"
        @update:model-value="onCampus"
      />
    </UFormField>

    <fieldset>
      <legend class="mb-2 text-meta font-semibold text-graphite">
        {{ t('listing.filter.price') }}
      </legend>
      <div class="grid grid-cols-2 gap-2">
        <UInput
          :model-value="state.priceMin"
          type="number"
          inputmode="numeric"
          min="0"
          step="500"
          :placeholder="t('listing.filter.priceMin')"
          :aria-label="`${t('listing.filter.price')} ${t('listing.filter.priceMin')}`"
          @change="(e: Event) => emit('update', { priceMin: toInt(e) })"
        />
        <UInput
          :model-value="state.priceMax"
          type="number"
          inputmode="numeric"
          min="0"
          step="500"
          :placeholder="t('listing.filter.priceMax')"
          :aria-label="`${t('listing.filter.price')} ${t('listing.filter.priceMax')}`"
          @change="(e: Event) => emit('update', { priceMax: toInt(e) })"
        />
      </div>
    </fieldset>

    <fieldset>
      <legend class="text-meta font-semibold text-graphite">
        {{ t('listing.filter.dates') }}
      </legend>
      <p class="mb-2 text-meta text-graphite">
        {{ t('listing.filter.datesHint') }}
      </p>
      <div class="flex flex-col gap-2">
        <UInput
          :model-value="state.from"
          type="date"
          :min="today"
          :aria-label="`${t('listing.filter.dates')}: ${t('listing.filter.from')}`"
          @change="(e: Event) => emit('update', { from: toDate(e), type: state.type ?? 'rent' })"
        />
        <UInput
          :model-value="state.to"
          type="date"
          :min="state.from ?? today"
          :aria-label="`${t('listing.filter.dates')}: ${t('listing.filter.to')}`"
          @change="(e: Event) => emit('update', { to: toDate(e), type: state.type ?? 'rent' })"
        />
      </div>
    </fieldset>

    <div v-if="activeCount > 0">
      <UButton
        color="neutral"
        variant="link"
        class="px-0"
        :label="t('common.actions.reset')"
        @click="emit('reset')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { IListingFilterState } from '../model/use-listing-filter';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { toDateInput } from 'shared/lib/utils/format';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  state: IListingFilterState;
  categories: string[];
  campuses: string[];
  activeCount: number;
}>();

const emit = defineEmits<{
  update: [patch: Partial<IListingFilterState>];
  reset: [];
}>();

// USelect does not accept an empty string as an item value
const ANY = '__any';

const { t } = useI18n();
const { categoryLabel, campusLabel } = useDictLabels();

const today = toDateInput(new Date());

const categoryOptions = computed(() => [
  { value: undefined, label: t('listing.filter.anyCategory') },
  ...props.categories.map(value => ({ value, label: categoryLabel(value) })),
]);

const campusOptions = computed(() => [
  { value: ANY, label: t('listing.filter.anyCampus') },
  ...props.campuses.map(value => ({ value, label: campusLabel(value) })),
]);

function onCampus(value: unknown) {
  emit('update', { campus: value === ANY || typeof value !== 'string' ? undefined : value });
}

function toInt(event: Event) {
  const value = Number.parseInt((event.target as HTMLInputElement).value, 10);
  return Number.isFinite(value) && value >= 0 ? value : undefined;
}

function toDate(event: Event) {
  return (event.target as HTMLInputElement).value || undefined;
}
</script>
