import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies }) => {
	const apiUrl = env.PRIVATE_API_URL + '/v1/models/';
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const modelResponse = await fetch(apiUrl, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}
		const models: ModelsResponse = await modelResponse.json();

		return {
			SESSION_TOKEN: sessionToken,
			models
		};
	} catch (error) {
		console.error(error);
	}

	return {};
}) satisfies PageServerLoad;
