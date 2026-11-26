<template>
  <section class="flex flex-col gap-6">
    <div class="flex flex-wrap items-end justify-between gap-4">
      <h1 class="text-heading font-semibold">
        {{ t('listing.my.title') }}
        <span v-if="!isLoading && !hasError && total" class="ml-1 font-mono text-title font-normal text-graphite tabular">{{ total }}</span>
      </h1>
      <UButton :to="{ name: PageName.LISTING_NEW }" :label="t('common.nav.post')" />
    </div>

    <!-- loading -->
    <ul
      v-if="isLoading"
      class="flex flex-col gap-3"
      aria-hidden="true"
    >
      <li
        v-for="i in SKELETONS"
        :key="i"
        class="flex items-center gap-4"
      >
        <USkeleton class="h-16 w-24 rounded-[3px]" />
        <div class="flex flex-1 flex-col gap-2">
          <USkeleton class="h-4 w-1/2" />
          <USkeleton class="h-3 w-1/4" />
        </div>
      </li>
    </ul>

    <!-- error -->
    <StateError v-else-if="hasError" @retry="load" />

    <!-- empty -->
    <StateEmpty
      v-else-if="isEmpty"
      :title="t('listing.my.empty')"
      :hint="t('listing.my.emptyHint')"
    >
      <UButton :to="{ name: PageName.LISTING_NEW }" :label="t('common.nav.post')" />
    </StateEmpty>

    <!-- success -->
    <ul v-else class="divide-y divide-rule border-y border-rule">
      <li
        v-for="listing in items"
        :key="listing.id"
        class="flex items-center gap-4 py-3"
      >
        <ListingCover
          :src="listing.photos[0]"
          :title="listing.title"
          class="aspect-[4/3] w-20 shrink-0 sm:w-24 [&_span]:text-[0.75rem]"
          :class="listing.status !== 'active' && 'opacity-60'"
        />
        <div class="flex min-w-0 flex-1 flex-col gap-1">
          <RouterLink
            :to="{ name: PageName.LISTING, params: { id: listing.id } }"
            class="truncate font-medium underline-offset-2 hover:underline"
          >
            {{ listing.title }}
          </RouterLink>
          <p class="flex flex-wrap items-center gap-x-4 gap-y-1 text-meta text-graphite">
            <span>{{ t(`listing.type.${listing.type}`) }}</span>
            <span class="inline-flex items-center gap-1.5" :title="t(`listing.my.hint.${listing.status}`)">
              <span
                class="size-2 rounded-[1px]"
                :class="STATUS_DOT[listing.status]"
                aria-hidden="true"
              />
              {{ t(`listing.status.${listing.status}`) }}
            </span>
            <time :datetime="listing.created_at">{{ formatRelative(listing.created_at, locale) }}</time>
          </p>
        </div>
        <ListingPrice :listing="listing" class="hidden max-w-48 sm:flex" />
        <ListingStatusMenu
          :listing="listing"
          @edit="router.push({ name: PageName.LISTING_EDIT, params: { id: listing.id } })"
        />
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import type { ListingStatus } from 'domains/entities/listing';
import { ListingCover, ListingPrice, useListingService } from 'domains/entities/listing';
import { ListingStatusMenu } from 'domains/features/listing-status';
import { PageName } from 'shared/lib/types/app/pages';
import { formatRelative } from 'shared/lib/utils/format';
import { StateEmpty, StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';

// ponytail: one page of 100 listings; add pagination when someone has more
const MY_PAGE_SIZE = 100;
const SKELETONS = 4;

const STATUS_DOT: Record<ListingStatus, string> = {
  active: 'bg-moss-500',
  paused: 'border border-graphite',
  archived: 'bg-rule-strong',
};

const { t, locale } = useI18n();
const router = useRouter();
const { items, total, isEmpty, isLoading, hasError, fetchListings } = useListingService();

function load() {
  fetchListings({
    mine: true,
    list_params: { page: 0, page_size: MY_PAGE_SIZE, with_total_count: true, sort: ['-created_at'] },
  });
}

load();
</script>
