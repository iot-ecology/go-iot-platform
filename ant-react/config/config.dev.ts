
import { defineConfig } from 'umi';

export default defineConfig({
  define: {
    APP_ENV: 'dev',
    'process.env.domain' : 'http://localhost:8080'
  },
});
