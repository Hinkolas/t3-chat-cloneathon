import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { env } from '$env/dynamic/private';

export const load = (async () => {
	const apiUrl = env.PRIVATE_API_URL + '/v1/models/';

	try {
		// Fetch models
		const modelResponse = await fetch(apiUrl);

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}
		const models: ModelsResponse = await modelResponse.json();

		return { models };
	} catch (error) {
		console.error(error);
	}

	return {};
}) satisfies PageServerLoad;
