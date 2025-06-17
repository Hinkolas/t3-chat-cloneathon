<script lang="ts">
	import type { ModelsResponse, ChatHistoryResponse } from '$lib//types';
	interface Props {
		data: {
			models: ModelsResponse;
			chats: ChatHistoryResponse;
		};
	}
	let { data }: Props = $props();

	import { PanelLeft } from '@lucide/svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import Popup from '$lib/components/Popup.svelte';
	import Error from '$lib/components/Error.svelte';
	import { toggleSidebar } from '$lib/store';
	import { sidebarState } from '$lib/store';
	import { PUBLIC_HOST_URL } from '$env/static/public';

	let chats = $state(data.chats);
	let models = $state(data.models || {});

	$effect(() => {
		if ($sidebarState.refresh) {
			refreshChats();
		}
	});

	async function refreshChats() {

		try {
			// Fetch chat history
			const chatHistoryResponse = await fetch(`${PUBLIC_HOST_URL}/v1/chats/`);
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
	{#if !models }
		<Error />
	{:else}
		<Popup />
		<button onclick={toggleSidebar} class="sidebar-button">
			<PanelLeft size="16" />
		</button>
		<Sidebar {chats} />
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
		border: 1px solid #88888822;
		border-radius: 8px;
		padding: 8px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease-out;
	}

	@media (max-width: 1024px) {
		.sidebar-button {
			z-index: 100;
		}
	}

	.sidebar-button:hover {
		background-color: #88888822;
	}

	.content {
		position: relative;
		width: 100%;
		height: 100%;
		display: flex;
	}
</style>
