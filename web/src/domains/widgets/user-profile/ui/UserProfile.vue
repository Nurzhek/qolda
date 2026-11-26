<template>
  <!-- loading -->
  <div v-if="userLoading" class="flex flex-col gap-3">
    <USkeleton class="h-8 w-64" />
    <USkeleton class="h-4 w-40" />
  </div>

  <!-- error -->
  <StateError v-else-if="userError" @retry="load" />

  <!-- success -->
  <div v-else-if="user" class="flex flex-col gap-10">
    <header class="flex flex-wrap items-end justify-between gap-4 border-b border-rule pb-6">
      <div class="flex flex-col gap-1">
        <h1 class="text-display font-semibold tracking-[-0.02em]">
          {{ user.name }}
        </h1>
        <p class="flex flex-wrap gap-x-4 text-body text-graphite">
          <span v-if="user.campus">{{ campusLabel(user.campus) }}</span>
          <span>{{ t('profile.memberSince', { date: formatMonthYear(user.created_at, locale) }) }}</span>
        </p>
      </div>
      <UButton
        v-if="isMe"
        color="neutral"
        variant="outline"
        :to="{ name: PageName.PROFILE }"
        :label="t('common.actions.edit')"
      />
      <UButton
        v-else-if="contactUrl"
        :to="contactUrl"
        target="_blank"
        rel="noopener noreferrer"
      >
        <PhTelegramLogo :size="18" aria-hidden="true" />
        {{ t('profile.messageOnTelegram') }}
      </UButton>
    </header>

    <section class="flex flex-col gap-5">
      <h2 class="text-title font-semibold">
        {{ t('profile.publicListings', { name: user.name }) }}
      </h2>

      <ListingGrid
        v-if="listingsLoading"
        loading
        :skeletons="4"
      />
      <StateError v-else-if="listingsError" @retry="loadListings" />
      <StateEmpty v-else-if="isEmpty" :title="t('profile.noListings', { name: user.name })" />
      <ListingGrid
        v-else
        v-slot="{ listing }"
        :items="items"
      >
        <ListingCard :listing="listing" :to="{ name: PageName.LISTING, params: { id: listing.id } }" />
      </ListingGrid>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ListingCard, ListingGrid, useListingService } from 'domains/entities/listing';
import { telegramLink, useUserService, useUserStore } from 'domains/entities/user';
import { storeToRefs } from 'pinia';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { PageName } from 'shared/lib/types/app/pages';
import { formatMonthYear } from 'shared/lib/utils/format';
import { StateEmpty, StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  id: string;
}>();

const PROFILE_PAGE_SIZE = 48;

const { t, locale } = useI18n();
const { campusLabel } = useDictLabels();
const { me } = storeToRefs(useUserStore());
const { viewed: user, isLoading: userLoading, hasError: userError, fetchUser } = useUserService();
const {
  items,
  isEmpty,
  isLoading: listingsLoading,
  hasError: listingsError,
  fetchListings,
} = useListingService();

const isMe = computed(() => props.id === me.value?.id);
const contactUrl = computed(() => (user.value?.telegram ? telegramLink(user.value.telegram) : ''));

function loadListings() {
  fetchListings({
    owner_id: props.id,
    status: 'active',
    list_params: { page: 0, page_size: PROFILE_PAGE_SIZE, sort: ['-created_at'] },
  });
}

function load() {
  fetchUser(props.id);
  loadListings();
}

watch(() => props.id, load, { immediate: true });
</script>
