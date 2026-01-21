import adapterAuto from '@sveltejs/adapter-auto';
import adapterStatic from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

// Use adapter-node for production/Docker builds, adapter-auto for development
const isProduction = process.env.NODE_ENV === 'production';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// Use adapter-node for Docker/production deployments
		// Use adapter-auto for development (supports various platforms)
		adapter: isProduction
			? adapterStatic({
					pages: 'build',
					assets: 'build',
					fallback: 'index.html', // This must match your nginx try_files
					precompress: false,
					strict: true
				})
			: adapterAuto()
	}
};

export default config;
