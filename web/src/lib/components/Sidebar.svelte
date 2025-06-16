<script lang="ts">
	let { chats = $bindable() } = $props();

	import { fade } from 'svelte/transition';
	import { Pin, LogOut } from '@lucide/svelte';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import HistoryChat from '$lib/components/HistoryChat.svelte';
	import { showConfirmationPopup, showRenamePopup, popup } from '$lib/store';
	import type { ChatHistoryResponse, ChatHistoryData } from '$lib/types';
	import { toggleSidebar, sidebarState } from '$lib/store';
	import { get } from 'svelte/store';
	import { ChatApiService } from '$lib/utils/chatApi';
	import {
		groupChatsByTime,
		filterChatsBySearchTerm,
		type GroupedChats
	} from '$lib/utils/chatUtils';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let sidebarCollapsed: boolean = $derived($sidebarState.collapsed);
	let chatSearchTerm: string = $state('');
	let activeContextMenuId = $state<string | null>(null);
	let url = $derived(page.url.pathname.split('/').pop());
	const userLoggedIn: boolean = true; // TODO: change this to dynamic with Cookie's

	// Filter and group chats
	let filteredAndGroupedChats = $derived.by(() => {
		let filtered = chats.filter((chat: ChatHistoryData) => chat.is_pinned === false);
		filtered = filterChatsBySearchTerm(filtered, chatSearchTerm);
		return groupChatsByTime(filtered);
	});

	let pinnedChats: ChatHistoryResponse = $derived(
		chats
			.filter((chat: ChatHistoryData) => chat.is_pinned === true)
			.sort((a: ChatHistoryData, b: ChatHistoryData) => b.last_message_at - a.last_message_at)
	);

	function openPopup(id: string, chatTitle: string = 'this chat') {
		showConfirmationPopup({
			title: 'Delete Thread',
			description: `Are you sure you want to delete "${chatTitle}"? This action cannot be undone.`,
			primaryButtonName: 'Delete',
			primaryButtonFunction: () => {
				deleteChat(id);
			}
		});
	}

	function renameChat(chat: ChatHistoryData) {
		showRenamePopup({
			title: 'Rename Chat',
			description: 'Enter a new name for this chat:',
			inputValue: chat.title,
			inputPlaceholder: 'Enter chat name...',
			inputLabel: 'Chat Name',
			primaryButtonName: 'Rename',
			primaryButtonFunction: async () => {
				const currentState = get(popup);

				if (currentState.type === 'rename') {
					const newTitle = currentState.inputValue.trim();

					if (newTitle && newTitle !== chat.title) {
						const oldTitle = chat.title;
						chat.title = newTitle;

						try {
							await ChatApiService.updateChatTitle(chat.id, newTitle);
							chats = [...chats];
						} catch (error) {
							chat.title = oldTitle;
							console.error('Failed to rename chat:', error);
						}
					}
				}
			}
		});
	}

	async function deleteChat(id: string) {
		try {
			await ChatApiService.deleteChat(id);
			const index = chats.findIndex((chat: ChatHistoryData) => chat.id === id);
			if (index > -1) {
				chats.splice(index, 1);
				if (url === id) {
					goto(`/`);
				}
			}
		} catch (error) {
			console.error('Error deleting chat:', error);
		}
	}

	async function patchChat(chat: ChatHistoryData, pin: boolean) {
		try {
			await ChatApiService.updateChatPinStatus(chat.id, pin);
			chat.is_pinned = pin;
		} catch (error) {
			showConfirmationPopup({
				title: `Ups! Something ain't right..`,
				description: 'Some error happened while we tried to sync your data. Try again later.',
				primaryButtonName: 'Confirm',
				primaryButtonFunction: () => {
					// No additional action needed - popup will close automatically
				}
			});
			throw new Error(`${error}`);
		}
	}

	function handleContextMenuOpen(chatId: string) {
		activeContextMenuId = chatId;
	}

	let chatSections = $derived.by(() => {
		const sections = [
			{
				key: 'pinned',
				title: 'Pinned',
				chats: pinnedChats,
				icon: true
			},
			{
				key: 'today',
				title: 'Today',
				chats: filteredAndGroupedChats.today,
				icon: false
			},
			{
				key: 'yesterday',
				title: 'Yesterday',
				chats: filteredAndGroupedChats.yesterday,
				icon: false
			},
			{
				key: 'last7Days',
				title: 'Last 7 Days',
				chats: filteredAndGroupedChats.last7Days,
				icon: false
			},
			{
				key: 'last30Days',
				title: 'Last 30 Days',
				chats: filteredAndGroupedChats.last30Days,
				icon: false
			},
			{
				key: 'older',
				title: 'Older',
				chats: filteredAndGroupedChats.older,
				icon: false
			}
		];

		return sections.filter((section) => section.chats.length > 0);
	});

	let innerWidth = $state(0);

	onMount(() => {
		if (innerWidth <= 1024) {
			toggleSidebar();
		}
	});
</script>

<svelte:window bind:innerWidth />

{#if !sidebarCollapsed}
	<div
		class="sidebar-mobile-overlay"
		transition:fade={{ duration: 150 }}
		onclick={toggleSidebar}
		aria-label="Toggle Sidebar"
		role="button"
		tabindex="0"
		onkeydown={(e) => {
			if (e.key === 'Escape') {
				toggleSidebar();
			}
		}}
	></div>
{/if}
<div class="sidebar {sidebarCollapsed ? 'collapsed' : ''}">
	<div class="head">
		<div class="title">Chat</div>
		<div class="newChatButton">
			<a href="/">New Chat</a>
		</div>
		<div class="search-container">
			<SearchInput bind:value={chatSearchTerm} placeholder="Search your threads..." />
		</div>
		<div class="chat-wrapper">
			<div class="chats-container">
				{#each chatSections as section}
					<div class="day-title">
						{#if section.icon}
							<Pin size="14" />
						{/if}
						{section.title}
					</div>
					<div class="chats">
						{#each section.chats as chat}
							<HistoryChat
								{chat}
								isCurrent={url === chat.id}
								{patchChat}
								{openPopup}
								{renameChat}
								{activeContextMenuId}
								onContextMenuOpen={handleContextMenuOpen}
							/>
						{/each}
					</div>
				{/each}
			</div>
		</div>
	</div>
	<div class="foot">
		{#if userLoggedIn}
			<a href="/settings" class="account-button">
				<img src="https://placehold.co/100" alt="Profile" />
				<div class="info">
					<span class="username">Ertu K.</span>
					<span class="subscription">Free</span>
				</div>
			</a>
		{:else}
			<!-- TODO: change href link to login -->
			<a href="#" class="login-button">
				<LogOut size="16" />
				Login
			</a>
		{/if}
	</div>
</div>

<style>
	.sidebar-mobile-overlay {
		z-index: 98;
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: #00000088;
	}

	/* Hide overlay on desktop (when sidebar is not absolute) */
	@media (min-width: 1025px) {
		.sidebar-mobile-overlay {
			display: none;
		}
	}

	.sidebar {
		flex: 0 0 256px;
		max-width: 256px;
		padding-block: 20px;
		background-color: var(--sidebar-background);
		transition: margin 0.15s ease-in-out;
		max-height: 100dvh;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
		gap: 8px;
	}

	@media (max-width: 1024px) {
		.sidebar {
			width: 256px;
			z-index: 99;
			position: absolute;
			left: 0;
			top: 0;
			height: 100%;
			border-right: 1px solid #88888822;
			background-color: var(--sidebar-background);
		}
	}

	.sidebar.collapsed {
		margin-left: -256px;
	}

	.head {
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		flex: 1;
		min-height: 0;
		overflow: hidden;
	}

	.title {
		color: hsl(var(--secondary-foreground));
		font-size: 18px;
		font-weight: 700;
		text-shadow: 0 0 4px hsl(var(--primary) / 0.3);
	}

	.newChatButton {
		width: 100%;
		padding-inline: 16px;
	}

	.newChatButton a {
		all: unset;
		box-sizing: border-box;
		padding: 8px 16px;
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		background-color: hsl(var(--primary) / 0.2);
		box-shadow: 0px 0px 2px hsl(var(--primary));
		border-radius: 8px;
		cursor: pointer;
		font-size: 14px;
		font-weight: 600;
		line-height: 20px;
		color: hsl(var(--secondary-foreground));
		text-shadow: 0 0 4px hsl(var(--primary) / 0.5);
		transition: background-color 0.15s ease;
	}

	.newChatButton:hover a {
		background-color: hsl(var(--primary) / 0.8);
	}

	.search-container {
		width: 100%;
		padding-inline: 2px;
	}

	.chat-wrapper {
		position: relative;
		flex: 1;
		width: 100%;
		min-height: 0;
	}

	.chats-container {
		position: relative;
		height: 100%;
		width: 100%;
		display: flex;
		flex-direction: column;
		gap: 16px;
		overflow-y: auto;
		padding-inline: 16px;
		padding-block: 16px;
	}

	.chat-wrapper::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 20px;
		background: linear-gradient(to top, transparent, #1d131b);
		pointer-events: none;
		z-index: 1;
	}

	.chat-wrapper::after {
		content: '';
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 20px;
		background: linear-gradient(to bottom, transparent, #1d131b);
		pointer-events: none;
		z-index: 1;
	}

	/* Hide scrollbar */
	.chats-container::-webkit-scrollbar {
		width: 0px !important;
		height: 0px !important;
		display: none !important;
	}

	.day-title {
		font-size: 12px;
		font-weight: 600;
		color: var(--secondary);
		text-shadow: 0 0 4px hsl(var(--primary) / 0.2);
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.chats {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.foot {
		display: flex;
		width: 100%;
		padding-inline: 12px;
	}

	.login-button {
		all: unset;
		color: hsl(var(--secondary-foreground));
		font-size: 14px;
		font-weight: 500;
		white-space: nowrap;
		text-shadow: 0px 0px 4px hsl(var(--primary) / 0.8);

		display: flex;
		justify-content: center;
		align-items: center;
		gap: 12px;
		width: 100%;
		padding: 8px 20px;
		cursor: pointer;
		border-radius: 8px;
		background-color: transparent;
		transition: background-color 0.15s ease-out;
	}

	.login-button:hover {
		background-color: hsl(var(--primary) / 0.3);
	}

	.account-button {
		all: unset;
		color: hsl(var(--secondary-foreground));
		font-size: 14px;
		font-weight: 500;
		white-space: nowrap;
		text-shadow: 0px 0px 4px hsl(var(--primary) / 0.8);

		display: flex;
		justify-content: flex-start;
		align-items: center;
		gap: 12px;
		width: 100%;
		padding: 8px;
		cursor: pointer;
		border-radius: 8px;
		background-color: transparent;
		transition: background-color 0.15s ease-out;
	}

	.account-button:hover {
		background-color: hsl(var(--primary) / 0.3);
	}

	.account-button img {
		width: 32px;
		height: 32px;
		object-fit: cover;
		border-radius: 9999px;
	}

	.account-button .info {
		display: flex;
		flex-direction: column;
		gap: 0;
		line-height: 1.3;
		letter-spacing: 0.24px;
	}

	.account-button .username {
		font-size: 14px;
		font-weight: 700;
	}

	.account-button .subscription {
		font-size: 12px;
	}
</style>
