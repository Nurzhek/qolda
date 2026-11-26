import process from 'node:process';
import { fileURLToPath, URL } from 'node:url';
import ui from '@nuxt/ui/vite';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';

const API_TARGET = process.env.VITE_API_TARGET ?? 'http://localhost:8090';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    ui({
      // Single light palette: no dark theme, so there is only one set of tokens to maintain
      colorMode: false,
      ui: {
        colors: {
          primary: 'moss',
          neutral: 'paper',
        },
      },
      autoImport: {
        imports: ['vue', 'vue-router'],
        dts: 'src/auto-imports.d.ts',
      },
      components: {
        dts: 'src/components.d.ts',
        resolvers: [
          // Phosphor icons autoimport.
          (componentName) => {
            if (componentName.startsWith('Ph')) {
              return {
                name: componentName,
                from: '@phosphor-icons/vue',
              };
            }
          },
        ],
      },
    }),
  ],

  resolve: {
    alias: {
      shared: fileURLToPath(new URL('./src/shared', import.meta.url)),
      domains: fileURLToPath(new URL('./src/domains', import.meta.url)),
      pages: fileURLToPath(new URL('./src/pages', import.meta.url)),
      app: fileURLToPath(new URL('./src/app', import.meta.url)),
      i18n: fileURLToPath(new URL('./src/i18n', import.meta.url)),
    },
  },

  server: {
    proxy: {
      '/api': API_TARGET,
      '/files': API_TARGET,
    },
  },
});
