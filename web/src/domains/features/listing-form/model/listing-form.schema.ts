import type { ListingType } from 'domains/entities/listing';
import { LISTING_TYPES } from 'domains/entities/listing';
import * as yup from 'yup';

export const MAX_PHOTOS = 6;
const MAX_PRICE = 100_000_000;

type Translate = (key: string, named?: Record<string, unknown>) => string;

// Same rules as the backend: checked before sending, but the server has the final say
export function createListingSchema(t: Translate) {
  const isType = (type: ListingType) => (value: unknown) => value === type;
  const v = (key: string) => t(`listing.form.v.${key}`);

  return yup.object({
    type: yup.mixed<ListingType>().oneOf(LISTING_TYPES).required(v('required')),
    title: yup.string().trim().required(v('required')).min(3, v('titleMin')).max(120, v('titleMax')),
    description: yup.string().trim().max(4000, v('descriptionMax')).default(''),
    category: yup.string().required(v('required')),
    campus: yup.string().default(''),
    photos: yup.array(yup.string().required()).max(MAX_PHOTOS, t('listing.form.v.photosMax', { max: MAX_PHOTOS })).default([]),
    price: yup.number().when('type', {
      is: (type: ListingType) => type !== 'exchange',
      then: s => s.typeError(v('price')).required(v('price')).moreThan(0, v('price')).max(MAX_PRICE, v('price')).integer(v('price')),
      otherwise: s => s.notRequired(),
    }),
    deposit: yup.number().when('type', {
      is: isType('rent'),
      then: s => s.typeError(v('deposit')).required(v('deposit')).min(0, v('deposit')).max(MAX_PRICE, v('deposit')).integer(v('deposit')),
      otherwise: s => s.notRequired(),
    }),
    wanted: yup.string().when('type', {
      is: isType('exchange'),
      then: s => s.trim().required(v('required')).min(3, v('wantedMin')).max(200, v('wantedMax')),
      otherwise: s => s.notRequired(),
    }),
    available_from: yup.string().when('type', {
      is: isType('rent'),
      then: s => s.required(v('required')),
      otherwise: s => s.notRequired(),
    }),
    available_to: yup.string().when('type', {
      is: isType('rent'),
      then: s => s.required(v('required')).test('order', v('datesOrder'), function (value) {
        const from = this.parent.available_from as string | undefined;
        return !from || !value || value >= from;
      }),
      otherwise: s => s.notRequired(),
    }),
  });
}

export type ListingFormValues = yup.InferType<ReturnType<typeof createListingSchema>>;
