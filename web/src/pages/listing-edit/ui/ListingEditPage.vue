<template>
  <ListingEditor :id="id" />
</template>

<script setup lang="ts">
import { ListingEditor } from 'domains/widgets/listing-editor';

defineProps<{
  id: string;
}>();
</script>
