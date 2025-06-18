import type { LayoutServerLoad } from './$types';
import type { ProfileResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load: LayoutServerLoad = async ({ cookies, params, fetch }) => {
	const sessionToken = cookies.get('session_token');

	try {
		const profileResponse = await fetch(`${env.PRIVATE_API_URL}/v1/profile/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
		if (!profileResponse.ok) {
			throw error(500, 'Failed to fetch models');
		}
		const profile: ProfileResponse = await profileResponse.json();

		return {
			profile
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
};
