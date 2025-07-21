"use strict";

const globals = require("globals");

const { fixupPluginRules, includeIgnoreFile } = require("@eslint/compat");

const _import = require("eslint-plugin-import");
const _refresh = require("eslint-plugin-react-refresh");
const tsParser = require("@typescript-eslint/parser");
const js = require("@eslint/js");

const { FlatCompat } = require("@eslint/eslintrc");
const path = require("path");

const gitignorePath = path.resolve(__dirname, ".gitignore");

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all,
});

module.exports = [
    {
        ignores: ["*.js"],
    },
    includeIgnoreFile(gitignorePath),
    ...compat.extends("eslint:recommended"),
    {
        languageOptions: {
            globals: {
                ...globals.node,
            },

            sourceType: "commonjs",

            parserOptions: {
                project: "./tsconfig.json",
            },
        },

        rules: {
            "no-unused-vars": "error",
            semi: ["error", "always"],
            "comma-spacing": "error",
            "no-extra-semi": "error",

            quotes: [
                "error",
                "double",
                {
                    avoidEscape: true,
                },
            ],

            "no-var": "error",
        },
    },
    ...compat
        .extends(
            "eslint:recommended",
            "plugin:@typescript-eslint/recommended",
            "plugin:@typescript-eslint/recommended-requiring-type-checking",
            "plugin:react-hooks/recommended"
        )
        .map((config) => ({
            ...config,
            files: ["*.ts", "*.tsx", "**/*.ts", "**/*.tsx"],
        })),
    {
        files: ["**/**/*.ts", "**/**/*.tsx", "**/*.ts", "**/*.tsx"],

        plugins: {
            import: fixupPluginRules(_import),
            "react-refresh": fixupPluginRules(_refresh),
        },

        languageOptions: {
            parser: tsParser,
            ecmaVersion: 2023,
            sourceType: "script",

            parserOptions: {
                project: ["./tsconfig.json"],
            },
        },

        rules: {
            "no-var": "off",
            "no-unused-vars": "off",
            "@typescript-eslint/consistent-type-definitions": "error",
            "@typescript-eslint/consistent-type-imports": "error",
            "@typescript-eslint/no-inferrable-types": "off",
            "@typescript-eslint/no-this-alias": "off",
            "@typescript-eslint/no-non-null-assertion": "off",
            "@typescript-eslint/no-unused-vars": "error",
            "@typescript-eslint/prefer-ts-expect-error": "error",
            "@typescript-eslint/restrict-plus-operands": "off",
            "@typescript-eslint/triple-slash-reference": "off",
            "import/consistent-type-specifier-style": ["error", "prefer-top-level"],

            "import/order": [
                "error",
                {
                    groups: ["builtin", "external", "internal", "object", ["parent", "sibling"], "index", "type"],

                    pathGroups: [
                        {
                            pattern: "[!/]",
                            group: "type",
                            position: "before",
                        },
                        {
                            pattern: "../src/*",
                            group: "type",
                            position: "before",
                        },
                        {
                            pattern: "../types/*",
                            group: "type",
                            position: "after",
                        },
                    ],

                    distinctGroup: false,

                    alphabetize: {
                        order: "asc",
                    },

                    "newlines-between": "always-and-inside-groups",
                },
            ],
            "react-refresh/only-export-components": ["warn", { allowConstantExport: true }],
        },
    },
];
