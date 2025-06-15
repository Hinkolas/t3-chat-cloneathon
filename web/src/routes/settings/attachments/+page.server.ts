import type { PageServerLoad } from './$types';
import type { AttachmentResponse } from '$lib/types';
import { error } from '@sveltejs/kit';

export const load = (async ({ params, url, fetch }) => {
	const apiUrl = 'http://localhost:3141';

	try {
		// Fetch models
		const attachmentsResponse = await fetch(`${apiUrl}/v1/attachments/`);
		if (!attachmentsResponse.ok) {
			throw error(500, 'Failed to fetch attachments');
		}

		const attachments: AttachmentResponse = await attachmentsResponse.json();
		return { attachments };
	} catch (err) {
		console.error('Load function error:', err);
		throw error(500, 'Failed to load data');
	}
}) satisfies PageServerLoad;
