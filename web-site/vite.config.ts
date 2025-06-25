import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { visualizer } from 'rollup-plugin-visualizer';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [react(), visualizer()],
  resolve: {
    alias: {
      '@': '/src', // 配置别名，@ 指向 src 目录
    },
  },
  build: {
    outDir: path.resolve(__dirname, '..', 'server-site/static'),
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端服务地址
        changeOrigin: true, // 允许跨域
      }
    }
  }
});
