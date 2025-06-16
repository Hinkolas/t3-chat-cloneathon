import type { PageServerLoad } from './$types';
import type { ChatHistoryResponse } from '$lib/types';
import { error } from '@sveltejs/kit';

export const load = (async ({ params, url, fetch }) => {
	const apiUrl = 'http://localhost:3141';

	try {
		const chatHistoryResponse = await fetch(`${apiUrl}/v1/chats/`);
		if (!chatHistoryResponse.ok) {
			throw error(500, 'Failed to fetch chats');
		}
		const chats: ChatHistoryResponse = await chatHistoryResponse.json();

		return { chats };
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
