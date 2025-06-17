import type { LayoutServerLoad } from './$types';
import type { ModelsResponse, ChatHistoryResponse, ChatResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { PRIVATE_HOST_URL } from '$env/static/private';

export const load: LayoutServerLoad = async ({ params, fetch }) => {

	try {
		// Fetch models
		const modelResponse = await fetch(`${PRIVATE_HOST_URL}/v1/models/`);
		if (!modelResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const models: ModelsResponse = await modelResponse.json();

		// Fetch chat history
		const chatHistoryResponse = await fetch(`${PRIVATE_HOST_URL}/v1/chats/`);
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
