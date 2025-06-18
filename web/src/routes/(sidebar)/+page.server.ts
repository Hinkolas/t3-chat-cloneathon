import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse, ProfileResponse } from '$lib/types';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies }) => {
	const apiUrl = env.PRIVATE_API_URL + '/v1/models/';
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const modelResponse = await fetch(apiUrl, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}
		const models: ModelsResponse = await modelResponse.json();

		// Fetch profile data
		const profileResponse = await fetch(`${env.PRIVATE_API_URL}/v1/profile/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
		if (!profileResponse.ok) {
			throw new Error('Something happened during fetch profile data');
		}
		const profile: ProfileResponse = await profileResponse.json();

		return {
			SESSION_TOKEN: sessionToken,
			models,
			profile
		};
	} catch (error) {
		console.error(error);
	}

	return {};
}) satisfies PageServerLoad;
