import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { PRIVATE_API_URL } from '$env/static/private';

export const load = (async () => {
	const apiUrl = PRIVATE_API_URL + '/v1/models/';

	try {
		// Fetch models
		const modelResponse = await fetch(apiUrl);

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
