// Tenge is formatted the Kazakhstani way in any UI locale: "4 500 ₸"
const priceFormatter = new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 0 });

export function formatPrice(value: number) {
  return `${priceFormatter.format(value)}\u00A0₸`;
}

// YYYY-MM-DD -> "12 Mar" (localized); no year if it's the current one
export function formatDay(value: string, locale: string) {
  if (!value)
    return '';
  const date = new Date(`${value}T00:00:00`);
  const sameYear = date.getFullYear() === new Date().getFullYear();
  return new Intl.DateTimeFormat(locale, {
    day: 'numeric',
    month: 'short',
    ...(sameYear ? {} : { year: 'numeric' }),
  }).format(date);
}

export function formatMonthYear(value: string, locale: string) {
  if (!value)
    return '';
  return new Intl.DateTimeFormat(locale, { month: 'long', year: 'numeric' }).format(new Date(value));
}

const RELATIVE_STEPS: [Intl.RelativeTimeFormatUnit, number][] = [
  ['minute', 60],
  ['hour', 60],
  ['day', 24],
  ['week', 7],
  ['month', 4.345],
  ['year', 12],
];

// RFC3339 -> "3 days ago" (localized)
export function formatRelative(value: string, locale: string) {
  if (!value)
    return '';
  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
  let amount = (new Date(value).getTime() - Date.now()) / 1000;
  let unit: Intl.RelativeTimeFormatUnit = 'second';
  for (const [nextUnit, size] of RELATIVE_STEPS) {
    if (Math.abs(amount) < size)
      break;
    amount /= size;
    unit = nextUnit;
  }
  if (unit === 'second')
    return rtf.format(0, 'minute');
  return rtf.format(Math.round(amount), unit);
}

// "family_housing" -> "Family housing": fallback when the code has no translation
export function humanizeCode(code: string) {
  const text = code.replace(/[_-]+/g, ' ').trim();
  return text.charAt(0).toUpperCase() + text.slice(1);
}

export function toDateInput(date: Date) {
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}
