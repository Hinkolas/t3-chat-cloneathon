import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';

export const load = (async ({ cookies }) => {
	const sessionToken = cookies.get('session_token');

	try {
		// Fetch models
		const modelResponse = await fetch(`${env.PRIVATE_API_URL}/v1/auth/logout/`, {
			headers: {
				Authorization: `Bearer ${sessionToken}`
			}
		});

		if (!modelResponse.ok) {
			throw new Error('Something happened during Record');
		}

		// Clear the session cookie
		cookies.delete('session_token', { path: '/' });

		// Redirect to the login page
	} catch (error) {
		console.error(error);
	}
	throw redirect(302, '/auth/login');
}) satisfies PageServerLoad;
