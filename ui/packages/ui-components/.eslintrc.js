const path = require('path');

const commonJavaScriptRules = {
    'prettier/prettier': 'error',

    // Do not require implicit return value.
    'arrow-body-style': 'off', // override eslint-config-airbnb-base

    // Require braces even when block has one statement.
    curly: ['error', 'all'],

    // Forbid use of console in favor of proper error capturing.
    'no-console': 'error',

    'import/no-extraneous-dependencies': [
        'error',
        {
            devDependencies: [
                path.join(__dirname, '**/*.test.ts'),
                path.join(__dirname, '**/*.test.tsx'),
                path.join(__dirname, '**/*.stories.tsx'),
                path.join(__dirname, '.storybook/**/*'),
                path.join(__dirname, '.prettierrc.js'),
                path.join(__dirname, '.postcssrc.js'),
                path.join(__dirname, 'tailwind.config.js'),
                path.join(__dirname, 'jest.config.js'),
            ],
            optionalDependencies: false,
        },
    ],
};

const commonTypeScriptReactRules = {
    // Neither require nor forbid destructuring assignment for props, state, context.
    'react/destructuring-assignment': ['off'],

    // Forbid 'data-test-id' instead use 'data-testid' attribute name.
    'react/forbid-dom-props': [
        'error',
        {
            forbid: ['data-test-id'],
        },
    ],
};

// Cannot easily factor out JavaScript extensions because the order matters.
const commonTypeScriptReactExtensions = [
    'plugin:react/recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:@typescript-eslint/recommended-requiring-type-checking',
    'plugin:eslint-comments/recommended',
    'airbnb-typescript',
    'prettier',
    'prettier/@typescript-eslint',
    'prettier/react',
];

module.exports = {
    plugins: ['@typescript-eslint', 'prettier', 'jest', 'jest-dom', 'testing-library'],
    parser: '@typescript-eslint/parser',
    parserOptions: {
        project: './tsconfig.eslint.json',
        tsconfigRootDir: __dirname,
    },

    overrides: [
        {
            files: ['*.js'],
            env: {
                node: true,
            },
            extends: ['eslint:recommended', 'plugin:eslint-comments/recommended', 'prettier'],
            rules: {
                ...commonJavaScriptRules,
            },
        },
        {
            files: ['*.ts', '*.tsx'],
            env: {
                browser: true,
            },
            extends: [...commonTypeScriptReactExtensions],
            rules: {
                ...commonJavaScriptRules,
                ...commonTypeScriptReactRules,

                'jsx-a11y/label-has-associated-control': [
                    2,
                    {
                        labelAttributes: ['label'],
                        controlComponents: ['Field'],
                        depth: 3,
                    },
                ],

                // Provide ECMAScript default values instead of defaultProps.
                'react/require-default-props': 'off',
            },
        },
        {
            files: ['*.test.ts', '*.test.tsx'],
            env: {
                browser: true,
                jest: true,
            },
            extends: [
                ...commonTypeScriptReactExtensions,
                'plugin:jest/recommended',
                'plugin:jest-dom/recommended',
                'plugin:testing-library/react',
            ],
            rules: {
                ...commonJavaScriptRules,
                ...commonTypeScriptReactRules,

                'jest/no-focused-tests': 'error',
            },
        },
        {
            files: ['*.stories.tsx'],
            env: {
                browser: true,
            },
            extends: [...commonTypeScriptReactExtensions],
            rules: {
                ...commonJavaScriptRules,
                ...commonTypeScriptReactRules,

                // Do not require type checking for story components.
                'react/prop-types': [
                    'error',
                    {
                        skipUndeclared: true,
                    },
                ],

                /*
                // Do not ban <{}> because it means what it says for React components, see:
                // https://github.com/typescript-eslint/typescript-eslint/issues/2063#issuecomment-675156492
                '@typescript-eslint/ban-types': [
                    'error',
                    {
                        extendDefaults: true,
                        types: {
                            '{}': false,
                        },
                    },
                ],
                */
            },
        },
    ],
};
