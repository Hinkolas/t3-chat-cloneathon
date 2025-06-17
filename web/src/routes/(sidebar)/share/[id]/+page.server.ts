import type { PageServerLoad } from './$types';
import type { ChatResponse } from '$lib/types';
import { error, redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const load = (async ({ params, url, fetch }) => {
    const id = params.id;

    try {
        const chatResponse = await fetch(`${env.PRIVATE_API_URL}/v1/chats/${id}/`);
        if (chatResponse.status === 404) {
            console.error('Server error fetching shared chat:', chatResponse.statusText);
            throw redirect(302, '/');
        }
        const chat: ChatResponse = await chatResponse.json();

        return {
            chat
        };
    } catch (err) {
        // Re-throw SvelteKit redirects and errors
        if (err && typeof err === 'object' && 'status' in err && 'location' in err) {
            throw err; // This is a redirect
        }
        if (err && typeof err === 'object' && 'status' in err && 'body' in err) {
            throw err; // This is an error
        }

        console.error('Load function error:', err);
        throw error(500, 'Failed to load data');
    }
}) satisfies PageServerLoad;
