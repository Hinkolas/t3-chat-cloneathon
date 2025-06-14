import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';

export const load = (async () => {
	const url = 'http://localhost:3141';

	try {
		// Fetch models
		const modelResponse = await fetch(`${url}/v1/models/`);

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}
		const models: ModelsResponse = await modelResponse.json();

		return { models };
	} catch (error) {
		console.log(error);
	}

	return {};
}) satisfies PageServerLoad;
