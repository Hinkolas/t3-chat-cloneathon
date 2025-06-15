import type { LayoutServerLoad } from './$types';
import type { ModelsResponse, ChatHistoryResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ params, fetch }) => {
	const apiUrl = 'http://localhost:3141';

	try {
		// Fetch models
		const modelResponse = await fetch(`${apiUrl}/v1/models/`);
		if (!modelResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const models: ModelsResponse = await modelResponse.json();

		// Fetch chat history
		const chatHistoryResponse = await fetch(`${apiUrl}/v1/chats/`);
		if (!chatHistoryResponse.ok) {
			throw error(500, 'Failed to fetch chats');
		}
		const chats: ChatHistoryResponse = await chatHistoryResponse.json();

		return {
			models,
			chats
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
};
