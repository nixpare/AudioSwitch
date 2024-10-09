import { defineConfig } from 'vite';
import { resolve } from 'path';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
	root: "src",
	plugins: [
		solidPlugin()
	],
	server: {
        port: 3000,
    },
	build: {
		target: 'esnext',
		outDir: '../dist',
		emptyOutDir: true,
		rollupOptions: {
			input: {
				main: resolve(__dirname, 'src/index.html'),
				overlay: resolve(__dirname, 'src/overlay.html')
			}
		}
	},
});
