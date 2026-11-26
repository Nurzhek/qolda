<template>
  <div class="flex flex-col gap-2">
    <ListingCover
      :src="photos[active]"
      :title="title"
      class="aspect-[4/3] w-full [&>img]:object-contain"
    >
      <slot />
    </ListingCover>

    <ul v-if="photos.length > 1" class="flex gap-2 overflow-x-auto pb-1">
      <li
        v-for="(url, index) in photos"
        :key="url"
        class="shrink-0"
      >
        <button
          type="button"
          class="block size-16 overflow-hidden rounded-[2px] border-2 transition-colors"
          :class="index === active ? 'border-moss-500' : 'border-transparent opacity-70 hover:opacity-100'"
          :aria-label="t('listing.detail.showPhoto', { n: index + 1 })"
          :aria-current="index === active"
          @click="active = index"
        >
          <img
            :src="url"
            alt=""
            class="size-full object-cover"
            loading="lazy"
          >
        </button>
      </li>
    </ul>
    <p
      v-if="photos.length > 1"
      class="sr-only"
      aria-live="polite"
    >
      {{ t('listing.detail.photo', { n: active + 1, total: photos.length }) }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ListingCover } from 'domains/entities/listing';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  photos: string[];
  title: string;
}>();

const { t } = useI18n();
const active = ref(0);

watch(() => props.photos, () => {
  active.value = 0;
});
</script>
