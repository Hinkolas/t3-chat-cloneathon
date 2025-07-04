import type { PageServerLoad } from './$types';
import type { AttachmentResponse } from '$lib/types';
import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load = (async ({ cookies, params, url, fetch }) => {
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const attachmentsResponse = await fetch(`${env.PRIVATE_API_URL}/v1/attachments/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});
		if (!attachmentsResponse.ok) {
			throw error(500, 'Failed to fetch attachments');
		}

		const attachments: AttachmentResponse = await attachmentsResponse.json();
		return {
			SESSION_TOKEN: sessionToken,
			attachments
		};
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
