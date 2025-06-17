import type { PageServerLoad } from './$types';
import type { ProfileResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { PRIVATE_HOST_URL } from '$env/static/private';

export const load = (async ({ params, url, fetch }) => {
	try {
		// Fetch models
		const profileResponse = await fetch(`${PRIVATE_HOST_URL}/v1/profile/`);
		if (!profileResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const profile: ProfileResponse = await profileResponse.json();

		return { profile };
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
