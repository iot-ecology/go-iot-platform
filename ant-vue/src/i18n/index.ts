import { createI18n, I18n } from 'vue-i18n';
import en from './en';
import zhCHS from './zh-CHS.ts';

const messages = {
    en,
    zhCHS
};

const i18n: I18n = createI18n({
    legacy: false, // 使用 Composition API
    locale: 'zhCHS', // 默认语言
    fallbackLocale: 'zhCHS', // 回退语言
    messages
});

export default i18n;
