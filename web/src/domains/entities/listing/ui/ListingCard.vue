<template>
  <article
    class="group relative flex min-w-0 flex-col gap-2 rounded-[3px] outline-offset-4 has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-moss-500"
  >
    <ListingCover
      :src="listing.photos[0]"
      :title="listing.title"
      class="aspect-[4/3]"
    >
      <ListingTypeTag :type="listing.type" />
    </ListingCover>

    <div class="flex min-w-0 flex-col gap-0.5">
      <ListingPrice :listing="listing" />
      <h3 class="line-clamp-2 text-body font-medium">
        <RouterLink
          :to="to"
          class="underline-offset-2 group-hover:underline after:absolute after:inset-0 focus-visible:outline-none"
        >
          {{ listing.title }}
        </RouterLink>
      </h3>
      <p class="flex flex-wrap gap-x-3 text-meta text-graphite">
        <span v-if="listing.campus">{{ campusLabel(listing.campus) }}</span>
        <time :datetime="listing.created_at">{{ formatRelative(listing.created_at, locale) }}</time>
      </p>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { RouteLocationRaw } from 'vue-router';
import type { IListingDTO } from '../model/listing.dto';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { formatRelative } from 'shared/lib/utils/format';
import { useI18n } from 'vue-i18n';
import ListingCover from './ListingCover.vue';
import ListingPrice from './ListingPrice.vue';
import ListingTypeTag from './ListingTypeTag.vue';

defineProps<{
  listing: IListingDTO;
  to: RouteLocationRaw;
}>();

const { locale } = useI18n();
const { campusLabel } = useDictLabels();
</script>
