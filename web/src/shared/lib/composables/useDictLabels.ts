import { humanizeCode } from 'shared/lib/utils/format';
import { useI18n } from 'vue-i18n';

// Dictionary codes -> labels. Campuses come from backend config: an unknown code is humanized
export function useDictLabels() {
  const { t, te } = useI18n();

  function label(group: 'category' | 'campus', code: string) {
    if (!code)
      return '';
    const key = `dict.${group}.${code}`;
    return te(key) ? t(key) : humanizeCode(code);
  }

  return {
    categoryLabel: (code: string) => label('category', code),
    campusLabel: (code: string) => label('campus', code),
  };
}
