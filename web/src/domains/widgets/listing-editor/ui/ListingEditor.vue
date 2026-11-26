<template>
  <section class="mx-auto flex w-full max-w-2xl flex-col gap-8">
    <h1 class="text-heading font-semibold">
      {{ id ? t('listing.form.titleEdit') : t('listing.form.titleNew') }}
    </h1>

    <template v-if="id">
      <!-- loading -->
      <div v-if="isLoading" class="flex flex-col gap-4">
        <USkeleton class="h-16 w-full" />
        <USkeleton class="h-10 w-full" />
        <USkeleton class="h-32 w-full" />
      </div>

      <!-- error -->
      <StateError v-else-if="hasError" @retry="load" />

      <!-- success -->
      <ListingForm
        v-else-if="listing"
        :key="listing.id"
        :listing="listing"
        :categories="categories"
        :campuses="campuses"
        @saved="openListing"
        @cancel="router.back()"
      />
    </template>

    <ListingForm
      v-else
      :categories="categories"
      :campuses="campuses"
      :default-campus="me?.campus"
      @saved="openListing"
      @cancel="router.back()"
    />
  </section>
</template>

<script setup lang="ts">
import { useDictStore } from 'domains/entities/dict';
import { useListingService } from 'domains/entities/listing';
import { useUserStore } from 'domains/entities/user';
import { ListingForm } from 'domains/features/listing-form';
import { storeToRefs } from 'pinia';
import { PageName } from 'shared/lib/types/app/pages';
import { StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  id?: string;
}>();

const { t } = useI18n();
const router = useRouter();
const { categories, campuses } = storeToRefs(useDictStore());
const { me } = storeToRefs(useUserStore());
const { current: listing, isLoading, hasError, fetchListing } = useListingService();

async function load() {
  if (!props.id)
    return;
  const result = await fetchListing(props.id);
  // Can't edit someone else's listing: the backend refuses too, but don't let them fill the form for nothing
  if (result.ok && result.data.owner_id !== me.value?.id)
    router.replace({ name: PageName.ERROR_403 });
}

function openListing(id: string) {
  router.push({ name: PageName.LISTING, params: { id } });
}

watch(() => props.id, load, { immediate: true });
</script>
