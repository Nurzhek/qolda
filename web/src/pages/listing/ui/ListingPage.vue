<template>
  <ListingDetails :id="id" />
</template>

<script setup lang="ts">
import { ListingDetails } from 'domains/widgets/listing-details';

defineProps<{
  id: string;
}>();
</script>
