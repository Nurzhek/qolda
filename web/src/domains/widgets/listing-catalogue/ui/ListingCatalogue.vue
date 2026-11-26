<template>
  <section class="flex flex-col gap-5">
    <h1 class="sr-only">
      {{ t('common.tagline') }}
    </h1>

    <div class="flex items-center gap-2">
      <ListingSearch :model-value="state.q" @update:model-value="(q) => update({ q })" />
      <UButton
        class="lg:hidden"
        color="neutral"
        variant="outline"
        size="lg"
        @click="filtersOpen = true"
      >
        <PhSlidersHorizontal :size="18" aria-hidden="true" />
        {{ t('listing.filter.open') }}
        <span v-if="activeCount" class="font-mono text-meta tabular">{{ activeCount }}</span>
      </UButton>
    </div>

    <ListingTypeSwitch :model-value="state.type" @update:model-value="(type) => update({ type })" />

    <div class="grid gap-10 lg:grid-cols-[13.5rem_minmax(0,1fr)]">
      <aside class="hidden lg:block">
        <ListingFilterPanel
          v-bind="panelProps"
          @update="update"
          @reset="reset"
        />
      </aside>

      <div class="flex min-w-0 flex-col gap-5">
        <div class="flex min-h-9 items-center justify-between gap-3">
          <p class="text-meta text-graphite" aria-live="polite">
            <template v-if="!isLoading && !hasError">
              {{ t('listing.filter.results', total) }}
            </template>
          </p>
          <ListingSortSelect :model-value="state.sort" @update:model-value="(sort) => update({ sort })" />
        </div>

        <!-- loading -->
        <ListingGrid v-if="isLoading" loading />

        <!-- error -->
        <StateError v-else-if="hasError" @retry="load" />

        <!-- empty -->
        <StateEmpty
          v-else-if="isEmpty && hasAnyFilter"
          :title="t('listing.catalogue.emptyFiltered')"
          :hint="t('listing.catalogue.emptyFilteredHint')"
        >
          <UButton
            color="neutral"
            variant="outline"
            :label="t('common.actions.reset')"
            @click="resetAll"
          />
        </StateEmpty>
        <StateEmpty
          v-else-if="isEmpty"
          :title="t('listing.catalogue.empty')"
          :hint="t('listing.catalogue.emptyHint')"
        >
          <UButton :to="{ name: PageName.LISTING_NEW }" :label="t('common.nav.post')" />
        </StateEmpty>

        <!-- success -->
        <template v-else>
          <ListingGrid v-slot="{ listing }" :items="items">
            <ListingCard :listing="listing" :to="{ name: PageName.LISTING, params: { id: listing.id } }" />
          </ListingGrid>

          <UPagination
            v-if="total > PAGE_SIZE"
            class="pt-4"
            :page="state.page + 1"
            :total="total"
            :items-per-page="PAGE_SIZE"
            @update:page="goToPage"
          />
        </template>
      </div>
    </div>

    <USlideover
      v-model:open="filtersOpen"
      side="left"
      :title="t('listing.filter.title')"
    >
      <template #body>
        <ListingFilterPanel
          v-bind="panelProps"
          @update="update"
          @reset="reset"
        />
      </template>
      <template #footer>
        <UButton
          block
          size="lg"
          :label="t('common.actions.apply')"
          @click="filtersOpen = false"
        />
      </template>
    </USlideover>
  </section>
</template>

<script setup lang="ts">
import { useDictStore } from 'domains/entities/dict';
import { ListingCard, ListingGrid, useListingService } from 'domains/entities/listing';
import {
  ListingFilterPanel,
  ListingSearch,
  ListingSortSelect,
  ListingTypeSwitch,
  PAGE_SIZE,
  useListingFilter,
} from 'domains/features/listing-filter';
import { storeToRefs } from 'pinia';
import { PageName } from 'shared/lib/types/app/pages';
import { StateEmpty, StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';

// The widget is the only place where the filter feature and the listing entity meet.
const { t } = useI18n();
const router = useRouter();
const { categories, campuses } = storeToRefs(useDictStore());
const { state, params, activeCount, hasAnyFilter, update, reset } = useListingFilter();
const { items, total, isEmpty, isLoading, hasError, fetchListings } = useListingService();

const filtersOpen = ref(false);

const panelProps = computed(() => ({
  state: state.value,
  categories: categories.value,
  campuses: campuses.value,
  activeCount: activeCount.value,
}));

function load() {
  fetchListings(params.value);
}

function resetAll() {
  router.replace({ query: {} });
}

function goToPage(page: number) {
  update({ page: page - 1 });
  window.scrollTo({ top: 0 });
}

// The request starts in setup so the first render is already in the loading state
watch(params, load, { immediate: true });
</script>
