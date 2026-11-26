import antfu from '@antfu/eslint-config';

export default antfu({
  ignores: ['dist', 'src/auto-imports.d.ts', 'src/components.d.ts'],
  vue: {
    overrides: {
      'vue/block-order': ['error', {
        order: ['template', 'script', 'style'],
      }],
      'vue/max-attributes-per-line': ['error', {
        singleline: 2,
      }],
    },
  },
  javascript: {
    overrides: {
      'no-restricted-syntax': [
        'error',
        'TSExportAssignment',
      ],
    },
  },
  typescript: {
    overrides: {
      'ts/ban-ts-comment': 'off',
    },
  },
  stylistic: {
    indent: 2,
    quotes: 'single',
    semi: true,
    overrides: {
      'style/brace-style': ['error', '1tbs', { allowSingleLine: true }],
    },
  },
});
