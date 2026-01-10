import js from '@eslint/js'
import stylistic from '@stylistic/eslint-plugin'
import vueTS from '@vue/eslint-config-typescript'
import vue from 'eslint-plugin-vue'

const ignores = ['**/node_modules/', '.git/', 'dist/', '.output/', '.vite', 'wailsjs/']

// @ts-check
export default [
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  ...vueTS({
    extends: ['strict'],
  }),
  {
    ignores,
  },
  {
    files: ['**/*.{js,ts,vue}'],
    plugins: {
      '@stylistic': stylistic,
    },
    rules: {
      'vue/valid-v-for': 1,
      '@stylistic/member-delimiter-style': [
        'error',
        {
          multiline: {
            delimiter: 'semi',
            requireLast: true,
          },
          singleline: {
            delimiter: 'semi',
            requireLast: false,
          },
        },
      ],
      '@typescript-eslint/no-empty-object-type': 0,
      '@typescript-eslint/no-explicit-any': 0,
      '@typescript-eslint/no-extraneous-class': 0,
      '@typescript-eslint/no-non-null-assertion': 2,
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          'argsIgnorePattern': '^_',
          'caughtErrorsIgnorePattern': '^_[^_].*$|^_$',
          'destructuredArrayIgnorePattern': '^_',
          'varsIgnorePattern': '^_',
        },
      ],
      '@typescript-eslint/triple-slash-reference': 0,
      'arrow-body-style': 'off',
      'block-spacing': ['error', 'always'],
      'brace-style': ['error', 'stroustrup', {
        'allowSingleLine': true,
      }],
      'comma-dangle': ['error', 'always-multiline'],
      'comma-spacing': ['error', {
        'after': true, 'before': false,
      }],
      curly: ['error', 'multi-line'],
      'eol-last': ['error', 'always'],
      'indent': ['error', 2, {
        'ignoreComments': false,
        'offsetTernaryExpressions': true,
      }],
      'no-console': 0,
      'no-duplicate-imports': 0,
      'no-multiple-empty-lines': ['error', {
        'max': 1, 'maxBOF': 0, 'maxEOF': 0,
      }],
      'no-trailing-spaces': 'error',
      'no-undef': 0,
      'no-unused-vars': 0,
      'object-curly-spacing': ['error', 'always'],
      quotes: ['error', 'single'],
      semi: ['error', 'never'],
      'vue/attributes-order': ['error', {
        'alphabetical': true,
      }],
      'vue/html-indent': ['error', 2],
      'vue/max-attributes-per-line': [1, {
        multiline: 1,
        singleline: 3,
      }],
      'vue/multi-word-component-names': 0,
      'vue/no-undef-components': ['error'],
      'vue/require-default-prop': 0,
    },
  },
]
