import { redirect, type Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

interface SessionResponse {
	session?: Session;
}

interface Session {
	id: string;
	user_id: string;
	token: string;
	issued_at: number;
	renewed_at: number;
	time_to_live: number;
	is_verified: boolean;
}

export const handle: Handle = async ({ event, resolve }) => {
	// get UserID from cookie session (null if session invalid)
	const SESSION_ID: string | undefined = event.cookies.get('session_token');
	let data: SessionResponse | null = null;

	// fetch json data from api (including auth headers)
	try {
		const res = await fetch(`${env.PRIVATE_API_URL}/v1/auth/session/`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${SESSION_ID}`
			}
		});

		if (!res.ok) {
			// if not ok, return null
			// console.error("Session fetch failed:", res.status, res.statusText);
			throw Error(`Session fetch failed: ${res.status} ${res.statusText}`);
		}

		data = (await res.json()) as SessionResponse;
	} catch (error) {
		console.error(error);
		// if not, redirect to login
		// throw redirect(302, "/auth/login");
	}

	if (event.route.id !== '/auth/login' && !data) {
		// if not, redirect to login
		throw redirect(302, '/auth/login');
	}

	if (event.route.id === '/auth/login' && data) {
		// check if logged in
		// if not, redirect to login
		throw redirect(302, '/');
	}

	const response = await resolve(event);
	return response;
};
