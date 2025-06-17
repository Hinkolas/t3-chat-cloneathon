import type { PageServerLoad } from './$types';
import type { ProfileResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies, params, url, fetch }) => {
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const profileResponse = await fetch(`${env.PRIVATE_API_URL}/v1/profile/`, {
			headers: {
				'Authorization': `Bearer ${sessionToken}`,
			}
		});
		if (!profileResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const profile: ProfileResponse = await profileResponse.json();

		return { 
			SESSION_TOKEN: sessionToken,
			profile 
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
