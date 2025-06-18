import type { PageServerLoad } from './$types';
import type { ChatHistoryResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies, params, url, fetch }) => {
	const sessionToken = cookies.get('session_token');

	try {
		const chatHistoryResponse = await fetch(`${env.PRIVATE_API_URL}/v1/chats/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
		if (!chatHistoryResponse.ok) {
			throw error(500, 'Failed to fetch chats');
		}
		const chats: ChatHistoryResponse = await chatHistoryResponse.json();

		return { SESSION_TOKEN: sessionToken, chats };
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
