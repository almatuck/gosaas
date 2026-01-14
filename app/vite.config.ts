import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],

	resolve: {
		alias: {
			$src: 'src'
		}
	},

	server: {
		host: '0.0.0.0',
		port: 4291,
		proxy: {
			// Proxy API requests to Go backend during development
			'/api': {
				target: 'http://api:8847',
				changeOrigin: true
			},
			// Proxy health check
			'/health': {
				target: 'http://api:8847',
				changeOrigin: true
			},
			// Proxy subscription plans (public endpoint)
			'/subscription/plans': {
				target: 'http://api:8847',
				changeOrigin: true
			},
			// Proxy WebSocket connections
			'/ws': {
				target: 'ws://api:8847',
				ws: true,
				changeOrigin: true
			}
		}
	}
});
