<script lang="ts">
	let { data }: PageProps = $props();

	import { PanelLeft } from '@lucide/svelte';
	import Sidebar from './Sidebar.svelte';
	import type { PageProps } from './$types';
	import ChatView from './ChatView.svelte';

	let sidebarCollapsed = $state(false);
	function sendMessage(message: string) {
		console.log(message);
	}

	function toggleSidebar() {
		sidebarCollapsed = !sidebarCollapsed;
	}

	function newChat() {}

</script>

<div class="container">
	<button onclick={toggleSidebar} class="sidebar-button">
		<PanelLeft size=16 />
	</button>
	<Sidebar {sidebarCollapsed} {newChat} />
	<div class="content">
		<ChatView {sendMessage}/>
	</div>
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

	.content {
		position: relative;
		flex: 1;
		background-color: var(--chat-background);
	}

	button {
		all: unset;

		display: flex;
		justify-content: center;
		align-items: center;
		gap: 4px;
		padding: 8px;
		border-radius: 9999px;
		white-space: nowrap;
		font-size: 12px;
		line-height: 1rem;
		cursor: pointer;
		border: 1px solid #88888833;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease;
	}

	button:hover {
		background-color: var(--button-hover);
	}
</style>
