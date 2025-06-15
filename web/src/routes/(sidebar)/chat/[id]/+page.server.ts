import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';

export const load = (async ({ params, url, fetch }) => {
	const apiUrl = 'http://localhost:3141';

	const id = params.id;

	try {
		const chatResponse = await fetch(`${apiUrl}/v1/chats/${id}/`);
		if (!chatResponse.ok) {
			throw error(500, 'Failed to fetch chat');
		}
		const chat: ChatResponse = await chatResponse.json();

		// Fetch models
		const modelResponse = await fetch(`${apiUrl}/v1/models/`);
		if (!modelResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const models: ModelsResponse = await modelResponse.json();

		return {
			chat,
			models
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
