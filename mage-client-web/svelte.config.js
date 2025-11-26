import adapterAuto from '@sveltejs/adapter-auto';
import adapterNode from '@sveltejs/adapter-node';
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
			? adapterNode({
					// Output directory for the build
					out: 'build',
					// Precompress assets
					precompress: true
				})
			: adapterAuto()
	}
};

export default config;
