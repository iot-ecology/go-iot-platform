
import { defineConfig } from 'umi';

export default defineConfig({
  define: {
    APP_ENV: 'prod',
    'process.env.domain' : 'http://localhost:8080'

  },
});
