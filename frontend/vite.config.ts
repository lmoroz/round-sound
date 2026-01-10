import vue from '@vitejs/plugin-vue';
import autoprefixer from 'autoprefixer';
import * as path from 'path';
import tailwind from 'tailwindcss';
import {
  defineConfig,
} from 'vite';
import checker from 'vite-plugin-checker';
import eslint from 'vite-plugin-eslint';
import removeConsole from 'vite-plugin-remove-console';

const PATH_SRC = './src';
const PATH_COMPONENTS = `${PATH_SRC}/mvc/view/components`;
const PATH_PAGES = `${PATH_SRC}/mvc/view/pages`;
const PATH_UI = `${PATH_SRC}/mvc/view/ui`;
const PATH_SHARED = `${PATH_SRC}/_shared`;
const PATH_ASSETS = 'assets';

console.log('> Vite -> config: APPLICATION_PORT =', process.env.APPLICATION_PORT);

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    minify: true,
    rollupOptions: {
      external: ['src/**/story/**', 'src/**/*.stories.ts', 'tests/**'],
    },
    target: 'esnext',
  },
  css: {
    postcss: {
      plugins: [tailwind(), autoprefixer()],
    },
  },
  define: {
    'import.meta.env.APP_VERSION': JSON.stringify(process.env.npm_package_version),
  },
  esbuild: {
    supported: {
      'top-level-await': true, //browsers can handle top-level-await features
    },
  },
  plugins: [
    vue(),
    eslint(),
    checker({
      vueTsc: true,
    }),
    removeConsole(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, PATH_SRC),
      '@ui': path.resolve(__dirname, PATH_UI),
      '~': path.resolve(__dirname, PATH_SRC),
      '~assets': path.resolve(__dirname, PATH_ASSETS),
      '~components': path.resolve(__dirname, PATH_COMPONENTS),
      '~pages': path.resolve(__dirname, PATH_PAGES),
      '~shared': path.resolve(__dirname, PATH_SHARED),
    },
    preserveSymlinks: true,
  },
  server: {
    port: parseInt(process.env.APPLICATION_PORT || '5000'),
  },
});
