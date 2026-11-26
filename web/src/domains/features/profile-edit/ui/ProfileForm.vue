<template>
  <form
    class="flex flex-col gap-5"
    novalidate
    @submit="onSubmit"
  >
    <UFormField
      :label="t('profile.name')"
      :error="errors.name"
      name="name"
      required
    >
      <UInput
        v-model.trim="name"
        v-bind="nameAttrs"
        autocomplete="name"
        class="w-full"
      />
    </UFormField>

    <UFormField
      :label="t('profile.campus')"
      :error="errors.campus"
      name="campus"
      required
    >
      <USelect
        v-model="campus"
        :items="campusItems"
        :placeholder="t('profile.campusPlaceholder')"
        class="w-full"
        @update:open="(open: boolean) => !open && campusAttrs.onBlur?.()"
      />
    </UFormField>

    <UFormField
      :label="t('profile.telegram')"
      :help="t('profile.telegramHint')"
      :error="errors.telegram"
      name="telegram"
    >
      <UInput
        v-model.trim="telegram"
        v-bind="telegramAttrs"
        autocomplete="off"
        autocapitalize="off"
        spellcheck="false"
        class="w-full"
      >
        <template #leading>
          <span class="text-graphite">@</span>
        </template>
      </UInput>
    </UFormField>

    <div>
      <UButton
        type="submit"
        size="lg"
        :loading="isLoading"
        :label="submitLabel ?? t('common.actions.save')"
      />
    </div>
  </form>
</template>

<script setup lang="ts">
import type { IUserDTO } from 'domains/entities/user';
import { toTypedSchema } from '@vee-validate/yup';
import { getFieldErrorMessage } from 'shared/api/error-handler';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { useForm } from 'vee-validate';
import { useI18n } from 'vue-i18n';
import * as yup from 'yup';
import { useProfileEdit } from '../model/use-profile-edit';

const props = defineProps<{
  user: IUserDTO;
  campuses: string[];
  submitLabel?: string;
}>();

const emit = defineEmits<{
  saved: [];
}>();

const TELEGRAM_RE = /^\w{5,32}$/;

const { t } = useI18n();
const { campusLabel } = useDictLabels();
const { saveProfile, isLoading } = useProfileEdit();

const campusItems = computed(() => props.campuses.map(code => ({ label: campusLabel(code), value: code })));

const { defineField, errors, handleSubmit, setErrors } = useForm({
  validationSchema: toTypedSchema(yup.object({
    name: yup.string().required(t('errors.field.required')).max(80, t('errors.field.too_long')),
    campus: yup.string().required(t('errors.field.required')),
    telegram: yup.string().default('').test('nick', t('errors.field.invalid'), v => !v || TELEGRAM_RE.test(v)),
  })),
  initialValues: {
    name: props.user.name,
    campus: props.user.campus || undefined,
    telegram: props.user.telegram.replace(/^@/, ''),
  },
});

const [name, nameAttrs] = defineField('name');
const [campus, campusAttrs] = defineField('campus');
const [telegram, telegramAttrs] = defineField('telegram');

const onSubmit = handleSubmit(async (values) => {
  const error = await saveProfile({ name: values.name, campus: values.campus, telegram: values.telegram ?? '' });
  const fields = error?.response?.data?.fields;
  if (fields) {
    setErrors(Object.fromEntries(Object.entries(fields).map(([k, v]) => [k, t(getFieldErrorMessage(v))])));
    return;
  }
  if (!error)
    emit('saved');
});
</script>
