import { useRouter } from "vue-router";

interface RouteJumpParams {
  path: string;
  query?: Record<string, any>;
  newWindow?: boolean;
}
export const useRouteJump = () => {
  const router = useRouter();
  const routeJump = ({ path = "", query = {}, newWindow = false }: RouteJumpParams) => {
    if (newWindow) {
      const route = router.resolve({ path, query });
      window.open(route.href, "_blank");
    } else {
      router.push({ path, query });
    }
  };
  return { routeJump };
};
