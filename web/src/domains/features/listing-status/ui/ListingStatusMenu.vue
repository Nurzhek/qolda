<template>
  <UDropdownMenu :items="items" :content="{ align: 'end' }">
    <UButton
      color="neutral"
      variant="ghost"
      square
      :loading="isLoading"
      :aria-label="t('listing.my.actions', { title: listing.title })"
    >
      <PhDotsThree
        :size="20"
        weight="bold"
        aria-hidden="true"
      />
    </UButton>
  </UDropdownMenu>
</template>

<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import type { IListingDTO, ListingStatus } from 'domains/entities/listing';
import { useI18n } from 'vue-i18n';
import { STATUS_ACTIONS, useListingStatus } from '../model/use-listing-status';

const props = defineProps<{
  listing: IListingDTO;
}>();

const emit = defineEmits<{
  edit: [];
}>();

const ACTION_LABEL: Record<ListingStatus, string> = {
  active: 'listing.my.activate',
  paused: 'listing.my.pause',
  archived: 'listing.my.archive',
};

const { t } = useI18n();
const { setStatus, isLoading } = useListingStatus();

const items = computed<DropdownMenuItem[]>(() => [
  { label: t('common.actions.edit'), onSelect: () => emit('edit') },
  ...STATUS_ACTIONS[props.listing.status].map(status => ({
    label: t(ACTION_LABEL[status]),
    onSelect: () => setStatus(props.listing.id, status),
  })),
]);
</script>
