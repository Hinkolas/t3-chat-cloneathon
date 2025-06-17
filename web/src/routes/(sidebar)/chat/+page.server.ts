import type { PageServerLoad } from './$types';
import type { ModelsResponse, ChatResponse } from '$lib/types';
import { error, redirect } from '@sveltejs/kit';
import { PRIVATE_HOST_URL } from '$env/static/private';

export const load = (async ({ params, url, fetch }) => {
    throw redirect(302, '/');    
}) satisfies PageServerLoad;
