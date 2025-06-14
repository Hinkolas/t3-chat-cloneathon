import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';

export const load = (async () => {
	const url = 'http://localhost:3141';

	try {
		// Fetch models
		const modelResponse = await fetch(`${url}/v1/models/`);

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}
		const models: ModelsResponse = await modelResponse.json();

		// Fetch chat history
		const chatResponse = await fetch(`${url}/v1/chats/`);

		if (!chatResponse.ok) {
			throw new Error('Something happened during Record');
		}

		const chats: ChatResponse = await chatResponse.json();
		return { models, chats };
	} catch (error) {
		console.log(error);
	}

	return {};
}) satisfies PageServerLoad;
