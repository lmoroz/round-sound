import js from '@eslint/js';
import stylistic from '@stylistic/eslint-plugin';
import vueTS from '@vue/eslint-config-typescript';
import jsonc from 'eslint-plugin-jsonc';
import perfectionist from 'eslint-plugin-perfectionist';
import playwright from 'eslint-plugin-playwright';
import storybook from 'eslint-plugin-storybook';
import tailwindcss from 'eslint-plugin-tailwindcss';
import vue from 'eslint-plugin-vue';

const ignores = ['**/node_modules/', '.git/', 'dist/', '.output/',  '.nitro', 'dist', '.output/', '.vite', 'tests/playwright.state.json', 'tests/server/e2e.json', 'test-results', 'scripts/'];

const customGroups = [
  {
    elementNamePattern: '~assets/*', groupName: 'assets',
  },
  {
    elementNamePattern: '~components/*', groupName: 'components',
  },
  {
    elementNamePattern: '^model/*', groupName: 'model',
  },
  {
    elementNamePattern: '~pages/*', groupName: 'pages',
  },
  {
    elementNamePattern: '^reka-ui', groupName: 'reka',
  },
  {
    elementNamePattern: '~shared/*', groupName: 'shared',
  },
  {
    elementNamePattern: '~tests/*', groupName: 'tests',
  },
  {
    elementNamePattern: '@ui/*', groupName: 'ui',
  },
  {
    elementNamePattern: ['~/utils/*', '~shared/utils/*'], groupName: 'utils',
  },
  {
    elementNamePattern: ['^vue', '^vue-*'], groupName: 'vue',
  },
  {
    elementNamePattern: ['^wire', '^wire-*'], groupName: 'wire',
  },
];
// @ts-check
export default [
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  ...vueTS({
    extends: ['strict'],
  }),
  ...jsonc.configs['flat/base'],
  {
    ignores,
  },
  {
    ...playwright.configs['flat/recommended'],
    files: ['tests/**'],
    rules: {
      ...playwright.configs['flat/recommended'].rules,
    },
  },
  {
    files: ['**/*.stories.@(ts|tsx|js|jsx|mjs|cjs)'],
    languageOptions: {
      ecmaVersion: 'latest',
      parserOptions: {
        parser: '@typescript-eslint/parser',
      },
    },
    plugins: {
      storybook,
    },
    rules: {
      'storybook/default-exports': 'error',
      'storybook/hierarchy-separator': 0,
    },
  },
  {
    plugins: {
      tailwindcss,
    },
    rules: {
      'tailwindcss/classnames-order': [
        'error',
        {
          removeDuplicates: true,
        },
      ],
      'tailwindcss/no-contradicting-classname': 2,
    },
  },
  {
    files: ['**/*.{js,ts,vue,json}'],
    plugins: {
      '@stylistic': stylistic,
      perfectionist,
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
      'class-methods-use-this': 0,
      'comma-dangle': ['error', 'always-multiline'],
      'comma-spacing': ['error', {
        'after': true, 'before': false,
      }],
      curly: ['error', 'all'],
      'eol-last': ['error', 'always'],
      'func-call-spacing': ['error', 'never'],
      'func-names': 0,
      'function-paren-newline': 0,
      'implicit-arrow-linebreak': 0,
      'indent': ['error', 2, {
        'ignoreComments': false,
        'offsetTernaryExpressions': true,
      }],
      'jsonc/indent': ['error', 2],
      'jsonc/key-name-casing': 0,
      'jsonc/sort-keys': ['error',
        {
          'order': [
            'name',
            'version',
            'private',
            'publishConfig',
          ],
          'pathPattern': '^$', // Hits the root properties
        },
        {
          'order': {
            'type': 'asc',
          },
          'pathPattern': '^(?:dev|peer|optional|bundled)?[Dd]ependencies$',
        },
      ],
      'key-spacing': ['error', {
        'multiLine': {
          'afterColon': true,
          'beforeColon': false,
        },
        'singleLine': {
          'afterColon': true,
          'beforeColon': false,
        },
      }],
      'keyword-spacing': ['error', {
        'after': true, 'before': true,
      }],
      'linebreak-style': 0,
      'lines-between-class-members': ['error', 'always', {
        'exceptAfterSingleLine': true,
      }],
      'no-console': 0,
      'no-duplicate-imports': 0,
      'no-else-return': 0,
      'no-extra-semi': 'error',
      'no-irregular-whitespace': 'error',
      'no-multiple-empty-lines': ['error', {
        'max': 1, 'maxBOF': 0, 'maxEOF': 0,
      }],
      'no-param-reassign': [2, {
        'props': false,
      }],
      'no-plusplus': 0,
      'no-restricted-imports': [
        'error',
        {
          'patterns': [
            {
              'group': ['../'],
              'message': 'Relative imports are not allowed.',
            },
          ],
        },
      ],
      'no-return-assign': 0,
      'no-sequences': 0,
      'no-trailing-spaces': 'error',
      'no-undef': 0,
      'no-underscore-dangle': 0,
      'no-unused-expressions': 0,
      'no-unused-vars': 0,
      'object-curly-newline': ['error', {
        'ExportDeclaration': {
          'minProperties': 2,
        },
        'ImportDeclaration': 'always',
        'ObjectExpression': 'always',
        'ObjectPattern': 'always',
      }],
      'object-curly-spacing': ['error', 'always', {
        'arraysInObjects': false,
        'objectsInObjects': false,
      }],
      'perfectionist/sort-classes': ['error', {
        groups: [
          'index-signature',
          'static-property',
          'static-method',
          'private-property',
          'property',
          'constructor',
          'private-method',
          'method',
        ],
        order: 'asc',
        type: 'natural',
      }],
      'perfectionist/sort-enums': 'off',
      'perfectionist/sort-exports': ['error', {
        order: 'asc',
        type: 'natural',
      }],
      'perfectionist/sort-imports': [
        'error',
        {
          customGroups: customGroups,
          environment: 'bun',
          groups: [
            'vue',
            'reka',
            'wire',
            ['builtin', 'external'],
            'type',
            'shared',
            'model',
            'internal-type',
            'internal',
            ['parent-type', 'sibling-type', 'index-type'],
            ['parent', 'sibling', 'index'],
            'object',
            'unknown',
            'components',
            'pages',
            'ui',
            'tests',
            'utils',
            'assets',
          ],
          ignoreCase: true,
          maxLineLength: undefined,
          newlinesBetween: 'always',
          order: 'asc',
          type: 'alphabetical',
        },
      ],
      'perfectionist/sort-interfaces': ['error', {
        order: 'asc',
        type: 'natural',
      }],
      'perfectionist/sort-objects': ['error', {
        customGroups: {
          id: 'id',
        },
        groups: ['id', 'unknown'],
        ignoreCase: true,
        ignorePattern: [
          'spacing',
          'container',
          'fontSize',
          'fontWeight',
          'borderWidth',
          'borderRadius',
          'lineHeight',
          'boxShadow',
        ],
        order: 'asc',
        type: 'natural',
      }],
      'perfectionist/sort-vue-attributes': 0,
      quotes: ['error', 'single'],
      semi: ['error', 'always'],
      'sort-imports': 0,
      'space-before-function-paren': ['error', 'never'],
      'space-in-parens': ['error', 'never'],
      'vue/attributes-order': ['error', {
        'alphabetical': true,
        'order': [
          'DEFINITION',
          'CONDITIONALS',
          'LIST_RENDERING',
          'OTHER_DIRECTIVES',
          'ATTR_DYNAMIC',
          ['UNIQUE', 'SLOT'],
          'EVENTS',
          'ATTR_STATIC',
          'RENDER_MODIFIERS',
          'GLOBAL',
          'TWO_WAY_BINDING',
          'CONTENT',
        ],
      }],
      'vue/comma-dangle': ['error', {
        'arrays': 'always-multiline',
        'exports': 'always-multiline',
        'functions': 'always-multiline',
        'imports': 'always-multiline',
        'objects': 'always-multiline',
      }],
      'vue/html-indent': ['error', 2],
      'vue/max-attributes-per-line': [1, {
        multiline: 1,
        singleline: 3,
      }],
      'vue/max-len': ['error', {
        code: 140,
        ignoreComments: true,
        ignoreHTMLAttributeValues: true,
        ignoreHTMLTextContents: true,
        ignoreStrings: true,
        ignoreTemplateLiterals: true,
        template: 120,
      }],
      'vue/multi-word-component-names': 0,
      'vue/multiline-html-element-content-newline': ['error'],
      'vue/no-empty-component-block': ['error'],
      'vue/no-multiple-template-root': 0,
      'vue/no-undef-components': ['error'],
      'vue/object-curly-newline': ['error', {
        'ExportDeclaration': {
          'consistent': true, 'minProperties': 1, 'multiline': true,
        },
        'ImportDeclaration': {
          'consistent': true, 'minProperties': 1, 'multiline': true,
        },
        'ObjectExpression': {
          'consistent': true, 'minProperties': 1, 'multiline': true,
        },
        'ObjectPattern': {
          'consistent': true, 'minProperties': 1, 'multiline': true,
        },
      }],
      'vue/require-default-prop': 0,
      'vue/require-v-for-key': 2,
      'vue/singleline-html-element-content-newline': ['error'],
    },
  },
];
