<template>
  <header class="sticky top-0 z-30 border-b border-rule bg-paper">
    <div class="mx-auto flex h-14 w-full max-w-7xl items-center gap-6 px-4 sm:px-6">
      <RouterLink
        :to="{ name: PageName.CATALOGUE }"
        class="shrink-0 rounded-[2px]"
        :aria-label="t('common.nav.browse')"
      >
        <AppWordmark />
      </RouterLink>

      <nav class="hidden items-center gap-5 md:flex">
        <RouterLink
          v-for="link in links"
          :key="link.name"
          :to="{ name: link.name }"
          class="text-body text-graphite hover:text-ink aria-[current=page]:font-medium aria-[current=page]:text-ink"
        >
          {{ link.label }}
        </RouterLink>
      </nav>

      <div class="ml-auto flex items-center gap-2">
        <UButton :to="{ name: PageName.LISTING_NEW }" :aria-label="t('common.nav.post')">
          <PhPlus
            :size="16"
            weight="bold"
            aria-hidden="true"
          />
          <span class="hidden sm:inline">{{ t('common.nav.post') }}</span>
        </UButton>

        <UDropdownMenu
          :items="menu"
          :content="{ align: 'end' }"
          :ui="{ content: 'w-56' }"
        >
          <UButton
            color="neutral"
            variant="ghost"
            :aria-label="t('common.nav.account')"
            class="gap-2 px-1.5"
          >
            <span
              class="grid size-7 place-items-center rounded-[3px] bg-paper-deep text-meta font-semibold"
              aria-hidden="true"
            >
              {{ initial }}
            </span>
            <span class="hidden max-w-32 truncate lg:inline">{{ me?.name }}</span>
            <PhCaretDown :size="12" aria-hidden="true" />
          </UButton>
        </UDropdownMenu>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { useUserStore } from 'domains/entities/user';
import { useAuth } from 'domains/features/auth';
import { storeToRefs } from 'pinia';
import { LOCALE_LIST } from 'shared/config/locale-config';
import { PageName } from 'shared/lib/types/app/pages';
import { useLocaleStore } from 'shared/stores/locale.store';
import { AppWordmark } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const router = useRouter();
const { me } = storeToRefs(useUserStore());
const localeStore = useLocaleStore();
const { signOut } = useAuth();

const initial = computed(() => me.value?.name.trim().charAt(0).toUpperCase() || '?');

const links = computed(() => [
  { name: PageName.CATALOGUE, label: t('common.nav.browse') },
  { name: PageName.MY_LISTINGS, label: t('common.nav.myListings') },
]);

const menu = computed<DropdownMenuItem[][]>(() => [
  [
    { type: 'label', label: me.value?.name ?? '' },
    { label: t('common.nav.profile'), onSelect: () => router.push({ name: PageName.PROFILE }) },
    { label: t('common.nav.myListings'), onSelect: () => router.push({ name: PageName.MY_LISTINGS }) },
  ],
  // the language switcher appears once there is more than one locale
  ...(LOCALE_LIST.length > 1
    ? [[
        {
          label: t('common.nav.language'),
          children: LOCALE_LIST.map(item => ({
            label: item.title,
            type: 'checkbox' as const,
            checked: localeStore.locale === item.key,
            onSelect: () => localeStore.setLocale(item.key),
          })),
        },
      ]]
    : []),
  [
    { label: t('common.nav.signOut'), onSelect: signOut },
  ],
]);
</script>
