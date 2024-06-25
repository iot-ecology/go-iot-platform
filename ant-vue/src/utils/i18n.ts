// utils/i18n.ts
import { useI18n } from 'vue-i18n';

export function getMetaTitle(key: string) {
    const { t } = useI18n();
    return t(key);
}