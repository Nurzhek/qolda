<template>
  <ErrorScreen
    :code="code"
    :title="t(copy.title)"
    :text="t(copy.text)"
  >
    <UButton
      v-if="code === '500'"
      :label="t('errors.page.reload')"
      @click="reload"
    />
    <UButton
      v-else
      :to="{ name: PageName.CATALOGUE }"
      :label="t('common.actions.goHome')"
    />
  </ErrorScreen>
</template>

<script setup lang="ts">
import { PageName } from 'shared/lib/types/app/pages';
import { ErrorScreen } from 'shared/ui';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  code: '403' | '404' | '500';
}>();

const COPY = {
  403: { title: 'errors.page.forbidden', text: 'errors.page.forbiddenText' },
  404: { title: 'errors.page.notFound', text: 'errors.page.notFoundText' },
  500: { title: 'errors.page.server', text: 'errors.page.serverText' },
} as const;

const { t } = useI18n();
const router = useRouter();
const copy = computed(() => COPY[props.code]);

// After a failure try to go back where we came from; otherwise to the catalogue
function reload() {
  const back = router.options.history.state.back;
  window.location.assign(typeof back === 'string' ? back : '/');
}
</script>
