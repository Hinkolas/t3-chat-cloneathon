import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies, params, url, fetch }) => {
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const modelResponse = await fetch(`${env.PRIVATE_API_URL}/v1/models/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
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
