<template>
  <form
    class="flex flex-col gap-7"
    novalidate
    @submit="onSubmit"
  >
    <ListingTypePicker
      :model-value="type"
      :disabled="!!listing"
      @update:model-value="(v) => (type = v)"
    />

    <section class="flex flex-col gap-5">
      <UFormField
        :label="t('listing.form.title')"
        :error="errors.title"
        name="title"
        required
      >
        <UInput
          v-model="title"
          v-bind="titleAttrs"
          :placeholder="t('listing.form.titlePlaceholder')"
          maxlength="120"
          class="w-full"
        />
      </UFormField>

      <div class="grid gap-5 sm:grid-cols-2">
        <UFormField
          :label="t('listing.form.category')"
          :error="errors.category"
          name="category"
          required
        >
          <USelect
            v-model="category"
            :items="categoryItems"
            :placeholder="t('listing.form.categoryPlaceholder')"
            class="w-full"
            @update:open="(open: boolean) => !open && categoryAttrs.onBlur?.()"
          />
        </UFormField>

        <UFormField
          :label="t('listing.form.campus')"
          :error="errors.campus"
          name="campus"
        >
          <USelect
            v-model="campus"
            :items="campusItems"
            class="w-full"
          />
        </UFormField>
      </div>

      <UFormField
        :label="t('listing.form.description')"
        :error="errors.description"
        name="description"
      >
        <UTextarea
          v-model="description"
          v-bind="descriptionAttrs"
          :placeholder="t('listing.form.descriptionPlaceholder')"
          :rows="5"
          autoresize
          class="w-full"
        />
      </UFormField>
    </section>

    <!-- Fields that depend on the listing type -->
    <section v-if="type" class="flex flex-col gap-5 border-t border-rule pt-6">
      <template v-if="type === 'rent' || type === 'sale'">
        <div class="grid gap-5 sm:grid-cols-2">
          <UFormField
            :label="type === 'rent' ? t('listing.form.pricePerDay') : t('listing.form.price')"
            :error="errors.price"
            name="price"
            required
          >
            <UInput
              v-model.number="price"
              v-bind="priceAttrs"
              type="number"
              inputmode="numeric"
              min="1"
              step="100"
              class="w-full"
              :ui="{ base: 'font-mono tabular' }"
            />
          </UFormField>

          <UFormField
            v-if="type === 'rent'"
            :label="t('listing.form.deposit')"
            :help="t('listing.form.depositHint')"
            :error="errors.deposit"
            name="deposit"
            required
          >
            <UInput
              v-model.number="deposit"
              v-bind="depositAttrs"
              type="number"
              inputmode="numeric"
              min="0"
              step="500"
              class="w-full"
              :ui="{ base: 'font-mono tabular' }"
            />
          </UFormField>
        </div>
      </template>

      <div v-if="type === 'rent'" class="grid gap-5 sm:grid-cols-2">
        <UFormField
          :label="t('listing.form.availableFrom')"
          :error="errors.available_from"
          name="available_from"
          required
        >
          <UInput
            v-model="availableFrom"
            v-bind="availableFromAttrs"
            type="date"
            class="w-full"
          />
        </UFormField>
        <UFormField
          :label="t('listing.form.availableTo')"
          :error="errors.available_to"
          name="available_to"
          required
        >
          <UInput
            v-model="availableTo"
            v-bind="availableToAttrs"
            type="date"
            :min="availableFrom"
            class="w-full"
          />
        </UFormField>
      </div>

      <UFormField
        v-if="type === 'exchange'"
        :label="t('listing.form.wanted')"
        :error="errors.wanted"
        name="wanted"
        required
      >
        <UInput
          v-model="wanted"
          v-bind="wantedAttrs"
          :placeholder="t('listing.form.wantedPlaceholder')"
          maxlength="200"
          class="w-full"
        />
      </UFormField>
    </section>

    <section class="flex flex-col gap-2 border-t border-rule pt-6">
      <h2 class="text-body font-medium">
        {{ t('listing.form.photos') }}
      </h2>
      <p class="text-meta text-graphite">
        {{ t('listing.form.photosHint', { max: MAX_PHOTOS }) }}
      </p>
      <PhotoUploader v-model="photos" :max="MAX_PHOTOS" />
      <p v-if="errors.photos" class="text-meta text-error">
        {{ errors.photos }}
      </p>
    </section>

    <div class="flex flex-wrap items-center gap-3 border-t border-rule pt-6">
      <UButton
        type="submit"
        size="lg"
        :loading="isSaving"
        :label="listing ? t('listing.form.submitEdit') : t('listing.form.submitNew')"
      />
      <UButton
        color="neutral"
        variant="ghost"
        size="lg"
        :label="t('common.actions.cancel')"
        @click="emit('cancel')"
      />
    </div>
  </form>
</template>

<script setup lang="ts">
import type { IListingDTO } from 'domains/entities/listing';
import { toTypedSchema } from '@vee-validate/yup';
import { getFieldErrorMessage } from 'shared/api/error-handler';
import { useDictLabels } from 'shared/lib/composables/useDictLabels';
import { useForm } from 'vee-validate';
import { useI18n } from 'vue-i18n';
import { createListingSchema, MAX_PHOTOS } from '../model/listing-form.schema';
import { toFormValues, useListingSave } from '../model/use-listing-form';
import ListingTypePicker from './ListingTypePicker.vue';
import PhotoUploader from './PhotoUploader.vue';

const props = defineProps<{
  listing?: IListingDTO;
  categories: string[];
  campuses: string[];
  defaultCampus?: string;
}>();

const emit = defineEmits<{
  saved: [id: string];
  cancel: [];
}>();

const { t } = useI18n();
const toast = useToast();
const { categoryLabel, campusLabel } = useDictLabels();
const { save, isSaving } = useListingSave();

const categoryItems = computed(() => props.categories.map(value => ({ value, label: categoryLabel(value) })));
const campusItems = computed(() => props.campuses.map(value => ({ value, label: campusLabel(value) })));

const { defineField, errors, handleSubmit, setErrors } = useForm({
  validationSchema: toTypedSchema(createListingSchema(t)),
  initialValues: props.listing
    ? toFormValues(props.listing)
    : { type: 'rent', campus: props.defaultCampus || undefined, photos: [] },
});

const [type] = defineField('type');
const [title, titleAttrs] = defineField('title');
const [description, descriptionAttrs] = defineField('description');
const [category, categoryAttrs] = defineField('category');
const [campus] = defineField('campus');
const [photos] = defineField('photos');
const [price, priceAttrs] = defineField('price');
const [deposit, depositAttrs] = defineField('deposit');
const [wanted, wantedAttrs] = defineField('wanted');
const [availableFrom, availableFromAttrs] = defineField('available_from');
const [availableTo, availableToAttrs] = defineField('available_to');

const onSubmit = handleSubmit(async (values) => {
  const result = await save(values, props.listing?.id);

  if (!result.ok) {
    const fields = result.error.response?.data?.fields;
    if (fields)
      setErrors(Object.fromEntries(Object.entries(fields).map(([k, v]) => [k, t(getFieldErrorMessage(v))])));
    return;
  }

  toast.add({ title: props.listing ? t('listing.form.saved') : t('listing.form.published'), color: 'success' });
  emit('saved', result.id);
});
</script>
