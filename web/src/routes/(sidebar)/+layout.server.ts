import type { LayoutServerLoad } from './$types';
import type { ModelsResponse, ChatHistoryResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load: LayoutServerLoad = async ({ cookies, params, fetch }) => {
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

		// Fetch chat history
		const chatHistoryResponse = await fetch(`${env.PRIVATE_API_URL}/v1/chats/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
		if (!chatHistoryResponse.ok) {
			throw error(500, 'Failed to fetch chats');
		}
		const chats: ChatHistoryResponse = await chatHistoryResponse.json();

		return {
			SESSION_TOKEN: sessionToken,
			models,
			chats
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
};
