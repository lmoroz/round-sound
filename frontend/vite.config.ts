import vue from '@vitejs/plugin-vue';
import * as path from 'path';
import {
  defineConfig,
} from 'vite';
import checker from 'vite-plugin-checker';

const PATH_SRC = './src';
const PATH_ASSETS = 'assets';

console.log('> Vite -> config: APPLICATION_PORT =', process.env.APPLICATION_PORT);

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    minify: true,
    target: 'esnext',
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
    checker({
      vueTsc: true,
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, PATH_SRC),
      '~': path.resolve(__dirname, PATH_SRC),
      '~assets': path.resolve(__dirname, PATH_ASSETS),
    },
  },
  server: {
    port: parseInt(process.env.APPLICATION_PORT || '5000'),
  },
});
