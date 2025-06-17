import type { PageServerLoad } from './$types';
import type { ChatHistoryResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { PRIVATE_HOST_URL } from '$env/static/private';

export const load = (async ({ params, url, fetch }) => {
	try {
		const chatHistoryResponse = await fetch(`${PRIVATE_HOST_URL}/v1/chats/`);
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
