import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { error, redirect } from '@sveltejs/kit';

export const load = (async ({ params, url, fetch }) => {
	const apiUrl = 'http://localhost:3141';

	const id = params.id;

	try {
		const chatResponse = await fetch(`${apiUrl}/v1/chats/${id}/`);
		if (chatResponse.status === 404) {
			console.error('Server error fetching chat:', chatResponse.statusText);
			throw redirect(302, '/');
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
		// Re-throw SvelteKit redirects and errors
		if (err && typeof err === 'object' && 'status' in err && 'location' in err) {
			throw err; // This is a redirect
		}
		if (err && typeof err === 'object' && 'status' in err && 'body' in err) {
			throw err; // This is an error
		}

		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
