<script lang="ts">
	import type { ModelsResponse, ChatResponse } from './types';
	interface Props {
		data: {
			models: ModelsResponse;
			chats: ChatResponse;
			// add other properties of data here if needed
		};
	}
	let { data }: Props = $props();

	import { PanelLeft } from '@lucide/svelte';
	import Sidebar from './Sidebar.svelte';
	import ChatView from './ChatView.svelte';
	import Popup from './Popup.svelte';
	import Error from './Error.svelte';

	let chats = $state(data.chats);
	let sidebarCollapsed = $state(true);
	function sendMessage(message: string) {
		console.log(message);
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
	}

	function newChat() {}
</script>

<div class="container">
	{#if !data.models || Object.keys(data.models).length === 0}
		<Error />
	{:else}
		<Popup />
		<button onclick={toggleSidebar} class="sidebar-button">
			<PanelLeft size="16" />
		</button>
		<Sidebar {chats} {sidebarCollapsed} {newChat} {toggleSidebar} />
		<div class="content">
			<ChatView {data} {sendMessage} />
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

	@media(max-width: 1024px) {
		.sidebar-button{
			z-index: 100;
		}
	}

	.sidebar-button:hover {
		background-color: #88888822;
	}

	.content {
		position: relative;
		width: 100%;
		display: flex;
	}
</style>
