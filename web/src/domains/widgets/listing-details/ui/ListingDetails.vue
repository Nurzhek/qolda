<template>
  <!-- loading -->
  <div v-if="isLoading" class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_22rem]">
    <USkeleton class="aspect-[4/3] w-full rounded-[3px]" />
    <div class="flex flex-col gap-3">
      <USkeleton class="h-5 w-16" />
      <USkeleton class="h-8 w-4/5" />
      <USkeleton class="h-10 w-1/2" />
      <USkeleton class="h-24 w-full" />
    </div>
  </div>

  <!-- error -->
  <StateError
    v-else-if="hasError"
    :message="errorMessage"
    @retry="load"
  />

  <!-- success (a single entity has no empty state: not found is an error) -->
  <article v-else-if="listing" class="grid gap-x-12 gap-y-8 lg:grid-cols-[minmax(0,1fr)_22rem]">
    <div class="flex min-w-0 flex-col gap-8">
      <ListingGallery :photos="listing.photos" :title="listing.title" />

      <section class="hidden lg:block">
        <ListingDescription :text="listing.description" />
      </section>
    </div>

    <div class="flex flex-col gap-6 lg:sticky lg:top-20 lg:self-start">
      <div class="flex flex-col gap-3">
        <div class="flex items-center gap-2">
          <ListingTypeTag :type="listing.type" />
          <span v-if="listing.status !== 'active'" class="text-meta text-graphite">
            {{ t(`listing.status.${listing.status}`) }}
          </span>
        </div>
        <h1 class="text-heading font-semibold tracking-[-0.01em] text-balance">
          {{ listing.title }}
        </h1>
        <ListingPrice :listing="listing" size="lg" />
      </div>

      <dl class="grid grid-cols-[auto_1fr] gap-x-6 gap-y-2 border-y border-rule py-4 text-body">
        <template v-if="listing.type === 'rent'">
          <dt class="text-graphite">
            {{ t('listing.deposit') }}
          </dt>
          <dd class="font-mono tabular">
            {{ listing.deposit > 0 ? formatPrice(listing.deposit) : t('listing.noDeposit') }}
          </dd>
          <dt class="text-graphite">
            {{ t('listing.available') }}
          </dt>
          <dd class="font-mono tabular">
            <time :datetime="listing.available_from">{{ formatDay(listing.available_from, locale) }}</time>
            –
            <time :datetime="listing.available_to">{{ formatDay(listing.available_to, locale) }}</time>
          </dd>
        </template>
        <dt class="text-graphite">
          {{ t('listing.detail.category') }}
        </dt>
        <dd>{{ categoryLabel(listing.category) }}</dd>
        <template v-if="listing.campus">
          <dt class="text-graphite">
            {{ t('listing.detail.campus') }}
          </dt>
          <dd>{{ campusLabel(listing.campus) }}</dd>
        </template>
      </dl>

      <section class="flex flex-col gap-4">
        <h2 class="sr-only">
          {{ t('listing.detail.postedBy') }}
        </h2>
        <UserBadge
          :name="ownerName"
          :meta="ownerMeta"
          :to="{ name: PageName.USER, params: { id: listing.owner_id } }"
        />

        <template v-if="isOwn">
          <p class="text-meta text-graphite">
            {{ t('listing.detail.yours') }}
          </p>
          <div class="flex gap-2">
            <UButton
              color="neutral"
              variant="outline"
              :to="{ name: PageName.LISTING_EDIT, params: { id: listing.id } }"
              :label="t('common.actions.edit')"
            />
            <ListingStatusMenu :listing="listing" @edit="router.push({ name: PageName.LISTING_EDIT, params: { id: listing.id } })" />
          </div>
        </template>
        <template v-else>
          <UButton
            v-if="contactUrl"
            :to="contactUrl"
            target="_blank"
            rel="noopener noreferrer"
            size="lg"
            block
          >
            <PhTelegramLogo :size="18" aria-hidden="true" />
            {{ t('profile.messageOnTelegram') }}
          </UButton>
          <p v-else class="text-body text-graphite">
            {{ t('listing.detail.noContact', { name: ownerName }) }}
          </p>
          <p class="text-meta text-graphite">
            {{ t('listing.detail.later') }}
          </p>
        </template>
      </section>

      <p class="text-meta text-graphite">
        <time :datetime="listing.created_at">{{ t('listing.posted', { when: formatRelative(listing.created_at, locale) }) }}</time>
      </p>

      <section class="lg:hidden">
        <ListingDescription :text="listing.description" />
      </section>
    </div>
  </article>
</template>

<script setup lang="ts">
import { ListingPrice, ListingTypeTag, useListingService } from 'domains/entities/listing';
import { telegramLink, UserBadge, useUserService, useUserStore } from 'domains/entities/user';
import { ListingStatusMenu } from 'domains/features/listing-status';
import { storeToRefs } from 'pinia';
import { getErrorMessage } from 'shared/api/error-handler';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { PageName } from 'shared/lib/types/app/pages';
import { formatDay, formatMonthYear, formatPrice, formatRelative } from 'shared/lib/utils/format';
import { StateError } from 'shared/ui';
import { useI18n } from 'vue-i18n';
import ListingDescription from './ListingDescription.vue';
import ListingGallery from './ListingGallery.vue';

const props = defineProps<{
  id: string;
}>();

const { t, locale } = useI18n();
const router = useRouter();
const { categoryLabel, campusLabel } = useDictLabels();
const { me } = storeToRefs(useUserStore());
const { current: listing, isLoading, hasError, fetchListing } = useListingService();
const { viewed, fetchUser } = useUserService();

const errorCode = shallowRef<string>();
const errorMessage = computed(() => (errorCode.value === 'object_not_found' ? t(getErrorMessage(errorCode.value)) : undefined));

// viewed may be left over from another profile: use it only if it is the owner
const owner = computed(() => (viewed.value?.id === listing.value?.owner_id ? viewed.value : null));
const isOwn = computed(() => !!listing.value && listing.value.owner_id === me.value?.id);
const ownerName = computed(() => owner.value?.name || listing.value?.owner?.name || '');
const ownerMeta = computed(() => {
  const campus = owner.value?.campus || listing.value?.owner?.campus;
  if (campus)
    return campusLabel(campus);
  return owner.value ? t('profile.memberSince', { date: formatMonthYear(owner.value.created_at, locale.value) }) : '';
});
const contactUrl = computed(() => (owner.value?.telegram ? telegramLink(owner.value.telegram) : ''));

async function load() {
  errorCode.value = undefined;
  const result = await fetchListing(props.id);
  if (!result.ok) {
    errorCode.value = result.error.response?.data?.code;
    return;
  }
  document.title = `${result.data.title} — qolda`;
  // The owner's contact lives in the profile; the page is still useful without it
  if (result.data.owner_id !== me.value?.id)
    fetchUser(result.data.owner_id);
}

watch(() => props.id, load, { immediate: true });
</script>
