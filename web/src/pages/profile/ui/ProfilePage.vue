<template>
  <section v-if="me" class="mx-auto flex w-full max-w-lg flex-col gap-8">
    <div class="flex flex-col gap-1">
      <h1 class="text-heading font-semibold">
        {{ t('profile.ownTitle') }}
      </h1>
      <p class="text-body text-graphite">
        {{ t('profile.memberSince', { date: formatMonthYear(me.created_at, locale) }) }}
      </p>
    </div>

    <dl class="flex flex-col gap-1 border-y border-rule py-4">
      <dt class="text-meta text-graphite">
        {{ t('profile.email') }}
      </dt>
      <dd class="font-medium break-all">
        {{ me.email }}
      </dd>
      <dd v-if="me.role !== 'student'" class="text-meta text-graphite">
        {{ t(`profile.roles.${me.role}`) }}
      </dd>
    </dl>

    <ProfileForm
      :user="me"
      :campuses="campuses"
      @saved="onSaved"
    />

    <p class="flex flex-wrap gap-x-6 gap-y-2 text-body">
      <RouterLink :to="{ name: PageName.MY_LISTINGS }" class="underline underline-offset-2 hover:text-moss-600">
        {{ t('profile.myListingsLink') }}
      </RouterLink>
      <RouterLink :to="{ name: PageName.USER, params: { id: me.id } }" class="underline underline-offset-2 hover:text-moss-600">
        {{ t('profile.publicListings', { name: me.name }) }}
      </RouterLink>
    </p>
  </section>
</template>

<script setup lang="ts">
import { useDictStore } from 'domains/entities/dict';
import { useUserStore } from 'domains/entities/user';
import { ProfileForm } from 'domains/features/profile-edit';
import { storeToRefs } from 'pinia';
import { PageName } from 'shared/lib/types/app/pages';
import { formatMonthYear } from 'shared/lib/utils/format';
import { useI18n } from 'vue-i18n';

const { t, locale } = useI18n();
const toast = useToast();
const { me } = storeToRefs(useUserStore());
const { campuses } = storeToRefs(useDictStore());

function onSaved() {
  toast.add({ title: t('profile.saved'), color: 'success' });
}
</script>
