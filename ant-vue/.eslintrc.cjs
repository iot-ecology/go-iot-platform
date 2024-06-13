module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: ["standard-with-typescript", "plugin:vue/vue3-recommended", "plugin:prettier/recommended"],
  parser: "vue-eslint-parser",
  parserOptions: {
    parser: "@typescript-eslint/parser",
    ecmaVersion: "latest",
    sourceType: "module",
    project: ["./tsconfig.json"],
    extraFileExtensions: [".vue"],
  },
  plugins: ["vue", "simple-import-sort"],
  rules: {
    // typescript
    "@typescript-eslint/strict-boolean-expressions": "off",
    "@typescript-eslint/no-floating-promises": "off",
    "@typescript-eslint/explicit-function-return-type": "off",

    // vue
    "vue/multi-word-component-names": [
      "error",
      {
        ignores: ["index"],
      },
    ],

    // 导入排序
    "simple-import-sort/imports": [
      "error",
      {
        groups: [["^vue", "^vue-router", "^\\w", "^"], ["^@/"], ["^\\.\\.(/.*)?$"], ["^\\./"], ["^@/types", ".*/types"]],
      },
    ],
    "simple-import-sort/exports": "error",
  },
};
