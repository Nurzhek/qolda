<template>
  <form
    class="flex flex-col gap-4"
    novalidate
    @submit="onSubmit"
  >
    <div>
      <h2 class="text-body font-semibold">
        {{ t('auth.dev.title') }}
      </h2>
      <p class="text-meta text-graphite">
        {{ t('auth.dev.hint') }}
      </p>
    </div>

    <UFormField
      :label="t('auth.dev.email')"
      :error="errors.email"
      name="email"
    >
      <UInput
        v-model.trim="email"
        v-bind="emailAttrs"
        type="email"
        autocomplete="email"
        :placeholder="domainHint"
        class="w-full"
      />
    </UFormField>

    <UFormField
      :label="t('auth.dev.name')"
      :error="errors.name"
      name="name"
    >
      <UInput
        v-model.trim="name"
        v-bind="nameAttrs"
        autocomplete="name"
        class="w-full"
      />
    </UFormField>

    <UButton
      type="submit"
      color="neutral"
      variant="outline"
      block
      :loading="isLoading"
      :label="t('auth.dev.submit')"
    />
  </form>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/yup';
import { getErrorMessage, getFieldErrorMessage } from 'shared/api/error-handler';
import { useForm } from 'vee-validate';
import { useI18n } from 'vue-i18n';
import * as yup from 'yup';
import { useAuth } from '../model/use-auth';

const props = defineProps<{
  redirect: string;
  domains: string[];
}>();

const { t } = useI18n();
const { signInDev, isLoading } = useAuth();

const domainHint = computed(() => (props.domains[0] ? `name@${props.domains[0]}` : ''));

const { defineField, errors, handleSubmit, setErrors } = useForm({
  validationSchema: toTypedSchema(yup.object({
    email: yup.string().required(t('errors.field.required')).email(t('errors.field.invalid')),
    name: yup.string().required(t('errors.field.required')).max(80, t('errors.field.too_long')),
  })),
});

const [email, emailAttrs] = defineField('email');
const [name, nameAttrs] = defineField('name');

const onSubmit = handleSubmit(async (values) => {
  const error = await signInDev(values, props.redirect);
  const data = error?.response?.data;
  if (!data)
    return;
  // Domain/blocked errors belong to the e-mail field: show them there, not only as a toast
  if (data.fields)
    setErrors(Object.fromEntries(Object.entries(data.fields).map(([k, v]) => [k, t(getFieldErrorMessage(v))])));
  else if (['domain_not_allowed', 'account_blocked'].includes(data.code))
    setErrors({ email: t(getErrorMessage(data.code)) });
});
</script>
