import { CodeInspectorPlugin } from "code-inspector-plugin";
import path from "path";
import { AntDesignVueResolver } from "unplugin-vue-components/resolvers";
import Components from "unplugin-vue-components/vite";
import { defineConfig, loadEnv } from "vite";
import eslintPlugin from "vite-plugin-eslint";
import { createSvgIconsPlugin } from "vite-plugin-svg-icons";
import vue from "@vitejs/plugin-vue";
const CWD = process.cwd();
// https://vitejs.dev/config/
// export default defineConfig({
//   base: "iot",
//   server: {
//     port: 5958,
//     open: true,
//     host: true,
//     proxy: {},
//   },
//   plugins: [
//     vue(),
//     createSvgIconsPlugin({
//       iconDirs: [path.resolve(__dirname, "./src/icons")],
//       symbolId: "icon-[dir]-[name]",
//     }),
//     Components({
//       resolvers: [
//         AntDesignVueResolver({
//           importStyle: false,
//         }),
//       ],
//     }),
//     CodeInspectorPlugin({
//       bundler: "vite",
//     }),
//   ],
//   resolve: {
//     alias: {
//       "@": path.resolve(__dirname, "./src"),
//     },
//   },
//   css: {
//     preprocessorOptions: {
//       less: {
//         modifyVars: {
//           hack: `true; @import "@/styles/variable.less";`,
//         },
//         javascriptEnabled: true,
//       },
//     },
//   },
//   build: {
//     reportCompressedSize: false,
//     chunkSizeWarningLimit: 1500,
//   },
// });

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  return {
    base: env.VITE_BASE_URL,
    server: {
      port: 5958,
      open: true,
      host: true,
      proxy: {},
    },
    plugins: [
      vue(),
      createSvgIconsPlugin({
        iconDirs: [path.resolve(__dirname, "./src/icons")],
        symbolId: "icon-[dir]-[name]",
      }),
      Components({
        resolvers: [
          AntDesignVueResolver({
            importStyle: false,
          }),
        ],
      }),
      CodeInspectorPlugin({
        bundler: "vite",
      }),
    ],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    css: {
      preprocessorOptions: {
        less: {
          modifyVars: {
            hack: `true; @import "@/styles/variable.less";`,
          },
          javascriptEnabled: true,
        },
      },
    },
    build: {
      reportCompressedSize: false,
      chunkSizeWarningLimit: 1500,
    },
  };
});
