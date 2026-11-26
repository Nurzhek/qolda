<template>
  <p
    v-if="listing.type === 'exchange'"
    class="flex min-w-0 items-baseline gap-1.5"
    :class="size === 'lg' ? 'text-lead' : 'text-body'"
  >
    <span class="shrink-0 text-graphite">{{ t('listing.swapFor') }}</span>
    <span class="truncate font-semibold">{{ listing.wanted }}</span>
  </p>
  <p v-else class="flex items-baseline gap-1.5">
    <span
      class="font-mono font-medium tracking-[-0.02em] tabular"
      :class="size === 'lg' ? 'text-display' : 'text-title'"
    >
      {{ formatPrice(listing.price) }}
    </span>
    <span
      v-if="listing.type === 'rent'"
      class="text-graphite"
      :class="size === 'lg' ? 'text-lead' : 'text-meta'"
    >
      {{ t('listing.perDay') }}
    </span>
  </p>
</template>

<script setup lang="ts">
import type { IListingDTO } from '../model/listing.dto';
import { formatPrice } from 'shared/lib/utils/format';
import { useI18n } from 'vue-i18n';

withDefaults(defineProps<{
  listing: Pick<IListingDTO, 'type' | 'price' | 'wanted'>;
  size?: 'md' | 'lg';
}>(), {
  size: 'md',
});

const { t } = useI18n();
</script>
