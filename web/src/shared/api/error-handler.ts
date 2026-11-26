import { i18n } from 'i18n/index';

// Backend error code -> translation key. Unknown code -> generic wording, never the framework text
export function getErrorMessage(code?: string) {
  const key = `errors.${code}`;
  return code && i18n.te(key) ? key : 'errors.unknown';
}

// Field error from the backend (a value of the fields map) -> translation key
export function getFieldErrorMessage(value: string) {
  const key = `errors.field.${value}`;
  return i18n.te(key) ? key : 'errors.field.invalid';
}
