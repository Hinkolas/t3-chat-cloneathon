import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { PRIVATE_HOST_URL } from '$env/static/private';

export const load = (async ({ params, url, fetch }) => {

	try {
		// Fetch models
		const modelResponse = await fetch(`${PRIVATE_HOST_URL}/v1/models/`);
		if (!modelResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const models: ModelsResponse = await modelResponse.json();

		return { models };
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
