import type { App } from "vue";

const directives = import.meta.glob("./*.ts");

const install = (app: App): void => {
  for (const [key, value] of Object.entries(directives)) {
    const name = key.slice(key.lastIndexOf("/") + 1, key.lastIndexOf("."));
    value()
      .then((res: any) => {
        app.directive(name, res.default || res[name]);
      })
      .catch((err) => {
        console.log(err);
      });
  }
};

export default install;
