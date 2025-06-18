<script lang="ts">
	import { closeSidebar, sidebarState } from '$lib/store';
	import { X, PinOff, Pin, TextCursor, Share2 } from '@lucide/svelte';
	import { scale } from 'svelte/transition';
	import { showConfirmationPopup } from '$lib/store';
	import { env } from '$env/dynamic/public';

	let {
		chat,
		isCurrent,
		patchChat,
		openPopup,
		renameChat,
		activeContextMenuId,
		SESSION_TOKEN,
		onContextMenuOpen
	} = $props();

	// Context menu state
	let showContextMenu = $derived(activeContextMenuId === chat.id);
	let contextMenuX = $state(0);
	let contextMenuY = $state(0);
	let contextMenuRef = $state<HTMLDivElement>();

	let isStreaming = $derived(false);
	let innerWidth = $state(0);

	// Handle right-click
	function handleContextMenu(event: MouseEvent) {
		event.preventDefault();
		contextMenuX = event.clientX;
		contextMenuY = event.clientY;
		onContextMenuOpen(chat.id);
	}

	// Close context menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (contextMenuRef && !contextMenuRef.contains(event.target as Node)) {
			onContextMenuOpen(null);
		}
	}

	// Context menu actions
	function handlePin(e: Event) {
		e.preventDefault();
		patchChat(chat, !chat.is_pinned);
		onContextMenuOpen(null);
	}

	function handleShare() {
		if (!chat.id) {
			console.error('Chat ID is not available for sharing.');
			return;
		}
		showConfirmationPopup({
			title: 'Share Thread',
			description: `Are you sure you want to share this chat? It will be publicly accessible.`,
			primaryButtonName: 'Confirm',
			primaryButtonFunction: () => {
				shareChat();
				onContextMenuOpen(null);
				navigator.clipboard.writeText(`${env.PUBLIC_HOST_URL}/share/${chat.id}/`);
			}
		});
		onContextMenuOpen(null);
	}

	async function shareChat() {
		try {
			const now = new Date();
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/${chat.id}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${SESSION_TOKEN}`
				},
				body: JSON.stringify({ shared_at: now.getTime() })
			});

			if (!response.ok) {
				throw new Error('Failed to share chat');
			}
		} catch (error) {
			console.error('Error sharing chat:', error);
		}
	}

	function handleDelete(e: Event) {
		e.preventDefault();
		openPopup(chat.id);
		onContextMenuOpen(null);
	}

	function handleRename() {
		renameChat(chat);
		onContextMenuOpen(null);
	}

	$effect(() => {
		if (showContextMenu) {
			document.addEventListener('click', handleClickOutside);
			return () => document.removeEventListener('click', handleClickOutside);
		}

		isStreaming = $sidebarState.chatIds.includes(chat.id);
	});
</script>

<svelte:window on:keydown={(e) => e.key === 'Escape' && onContextMenuOpen(null)} bind:innerWidth />

<a
	href="/chat/{chat.id}"
	class="chat"
	onclick={() => {
		if (innerWidth <= 1024) {
			closeSidebar();
		}
	}}
	oncontextmenu={handleContextMenu}
	class:active={isCurrent}
>
	<span>{chat.title}</span>
	<div class="loading" class:active={isStreaming}></div>
	<div class="buttons">
		<button
			onclick={(e) => {
				handlePin(e);
			}}
		>
			{#if chat.is_pinned}
				<PinOff size="14" />
			{:else}
				<Pin size="14" />
			{/if}
		</button>
		<button
			onclick={(e) => {
				handleDelete(e);
			}}
		>
			<X size="14" />
		</button>
	</div>
</a>

<!-- Context Menu -->
{#if showContextMenu}
	<div
		bind:this={contextMenuRef}
		class="context-menu"
		style="left: {contextMenuX}px; top: {contextMenuY}px;"
		transition:scale={{ duration: 100, start: 0.9 }}
	>
		<button class="context-menu-item" onclick={handlePin}>
			{#if chat.is_pinned}
				<PinOff size="16" />
				Unpin
			{:else}
				<Pin size="16" />
				Pin
			{/if}
		</button>
		<button class="context-menu-item" onclick={handleRename}>
			<TextCursor size="16" />
			Rename
		</button>
		<button class="context-menu-item" onclick={handleShare}>
			<Share2 size="16" />
			Share
		</button>
		<button class="context-menu-item" onclick={handleDelete}>
			<X size="16" />
			Delete
		</button>
	</div>
{/if}

<style>
	.chat {
		text-decoration: none;
		position: relative;
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 8px 8px;
		border-radius: 8px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		overflow: hidden;
		min-width: 0;
	}

	.chat span {
		color: var(--text);
		font-size: 14px;
		font-weight: 500;
		white-space: nowrap;
		text-overflow: ellipsis;
		overflow: hidden;
		flex: 1;
		min-width: 0;
	}

	.loading {
		/* Hidden by default */
		display: none;
		width: 16px;
		height: 16px;
		border: 2px solid #9d989c;
		border-top: 2px solid var(--primary-background-light);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.loading.active {
		display: inline-block;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	.chat:hover {
		background-color: var(--history-chat-hover);
	}

	.chat.active {
		background-color: var(--history-chat-hover);
	}

	.buttons {
		position: absolute;
		top: 50%;
		right: -100%;
		transform: translateY(-50%);
		flex-shrink: 0;

		padding-right: 8px;
		padding-left: 8px;
		display: flex;
		justify-content: flex-end;
		margin-left: auto;
		gap: 4px;
		transition: right 0.15s ease;
		background-color: var(--history-chat-hover);
	}

	.buttons::before {
		content: '';
		position: absolute;
		left: -12px;
		top: 0;
		height: 100%;
		width: 12px;
		background: linear-gradient(to right, transparent, var(--history-chat-hover));
	}

	.buttons button {
		display: flex;
		justify-content: center;
		align-items: center;

		background: none;
		border: none;
		border-radius: 4px;
		padding: 4px;
		cursor: pointer;
		color: var(--text);
		transition: background-color 0.15s ease-out;
	}

	.buttons button:hover {
		background-color: hsl(var(--primary) / 0.5);
	}

	.chat:hover .buttons {
		right: 0;
	}

	/* Context Menu Styles */
	.context-menu {
		position: fixed;
		z-index: 1000;
		background-color: var(--context-menu-background);
		border: 1px solid var(--context-menu-border);
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		padding: 4px;
		min-width: 120px;
	}

	.context-menu-item {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 12px;
		background: none;
		border: none;
		border-radius: 4px;
		text-align: left;
		cursor: pointer;
		font-size: 14px;
		color: var(--text);
		transition: background-color 0.15s ease;
	}

	.context-menu-item:hover {
		background-color: var(--primary-background-hover-light);
	}
</style>
