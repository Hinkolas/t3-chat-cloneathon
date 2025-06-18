import { fail, redirect, type Actions } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

interface LoginResponse {
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

export const actions: Actions = {
	default: async ({ cookies, request }) => {
		try {
			const formData = await request.formData();
			const username = formData.get('username') as string | null;
			const password = formData.get('password') as string | null;

			// Validate input
			if (!username || !password) {
				return fail(400, {
					error: 'Username and password are required',
					username
				});
			}

			const response = await fetch(`${env.PRIVATE_API_URL}/v1/auth/login/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			if (!response.ok) {
				// Handle different HTTP error codes
				if (response.status === 401) {
					return fail(401, {
						error: 'Invalid username or password',
						username
					});
				} else if (response.status === 429) {
					return fail(429, {
						error: 'Too many login attempts. Please try again later.',
						username
					});
				} else {
					return fail(response.status, {
						error: 'Login failed. Please try again.',
						username
					});
				}
			}

			const data: LoginResponse = await response.json();

			// Validate response data
			if (!data.session?.id) {
				return fail(500, {
					error: 'Invalid response from server',
					username
				});
			}

			cookies.set('session_token', data.session.token, {
				path: '/',
				sameSite: 'strict',
				secure: false, // TODO: Set to true in Production!!!
				httpOnly: true, // Recommended for security
				maxAge: 60 * 60 * 24 * 7 // 7 days
			});
		} catch (error) {
			// Handle network errors, JSON parsing errors, etc.
			if (error instanceof Error && error.message.includes('redirect')) {
				// Re-throw redirect errors (SvelteKit redirects are thrown as errors)
				throw error;
			}

			console.error('Login error:', error);
			return fail(500, {
				error: 'An unexpected error occurred. Please try again.'
			});
		}

		throw redirect(302, '/');
	}
};
