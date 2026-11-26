<template>
  <UserProfile :id="id" />
</template>

<script setup lang="ts">
import { UserProfile } from 'domains/widgets/user-profile';

defineProps<{
  id: string;
}>();
</script>
