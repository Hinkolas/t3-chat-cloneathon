<script lang="ts">
	import type { ModelsResponse, ChatHistoryResponse, ProfileResponse } from '$lib//types';
	interface Props {
		data: {
			SESSION_TOKEN: string;
			models: ModelsResponse;
			chats: ChatHistoryResponse;
			profile: ProfileResponse;
		};
	}
	let { data }: Props = $props();

	import { PanelLeft } from '@lucide/svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Popup from '$lib/components/Popup.svelte';
	import Error from '$lib/components/Error.svelte';
	import { toggleSidebar } from '$lib/store';
	import { sidebarState } from '$lib/store';
	import { env } from '$env/dynamic/public';

	let chats = $state(data.chats || {});
	let models = $state(data.models || {});
	let profile = $state(data.profile || {});
	let SESSION_TOKEN = $state(data.SESSION_TOKEN || '');

	$effect(() => {
		if ($sidebarState.refresh) {
			refreshChats();
		}
	});

	async function refreshChats() {
		try {
			// Fetch chat history
			const chatHistoryResponse = await fetch(`${env.PUBLIC_API_URL}/v1/chats/`, {
				headers: {
					Authorization: `Bearer ${SESSION_TOKEN}`
				}
			});
			if (!chatHistoryResponse.ok) {
				console.error(500, 'Failed to fetch chats');
			}
			const newChats: ChatHistoryResponse = await chatHistoryResponse.json();

			chats = newChats;
		} catch (err) {
			console.error('Load function error:', err);
		}
	}
</script>

<div class="container">
	{#if !models}
		<Error />
	{:else}
		<Popup />
		<button onclick={toggleSidebar} class="sidebar-button">
			<PanelLeft size="16" />
		</button>
		<Sidebar {chats} {SESSION_TOKEN} {profile} />
		<div class="content">
			<slot />
		</div>
	{/if}
</div>

<style>
	.container {
		position: relative;
		width: 100%;
		height: 100dvh;
		display: flex;
		overflow: hidden;
	}

	.sidebar-button {
		all: unset;
		position: absolute;
		top: 16px;
		left: 16px;

		z-index: 1;
		display: flex;
		justify-self: center;
		align-items: center;
		border: 1px solid var(--sidebar-toggle-button-border);
		border-radius: 8px;
		padding: 8px;
		cursor: pointer;
		color: var(--text);
		transition: background-color 0.15s ease-out;
	}

	@media (max-width: 1024px) {
		.sidebar-button {
			z-index: 100;
		}
	}

	.sidebar-button:hover {
		background-color: var(--sidebar-toggle-button-hover);
	}

	.content {
		position: relative;
		width: 100%;
		height: 100%;
		display: flex;
	}
</style>
