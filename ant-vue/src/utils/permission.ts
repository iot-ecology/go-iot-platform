import { usePermissionStore } from "@/stores/permission";

// 判断是否有权限
export const hasPermission = (id: string): boolean => {
  const permissionStore = usePermissionStore();
  const btnPermissionList = permissionStore.btnPermissionList;
  if (!permissionStore.openValidator) {
    return true;
  } else if (btnPermissionList.includes(id)) {
    return true;
  }
  return false;
};
