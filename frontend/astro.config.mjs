import { defineConfig } from 'astro/config';
import node from '@astrojs/node';

export default defineConfig({
  output: 'static',
  adapter: node({
    mode: 'standalone'
  }),
  build: {
    assets: 'assets',
    format: 'directory'
  },
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:8090'
    }
  }
});