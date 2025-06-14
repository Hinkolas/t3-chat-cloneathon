<script lang="ts">
	let { chats = $bindable(), sidebarCollapsed, newChat, toggleSidebar } = $props();

	import { fade } from 'svelte/transition';
	import { isMobile } from '$lib/deviceDetection';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import { Pin, LogOut } from '@lucide/svelte';

	// Updated import for new popup system
	import { showConfirmationPopup, showRenamePopup, popup } from '$lib/store';
	import { get } from 'svelte/store';

	import type { ChatResponse, ChatData } from '$lib/types';
	import HistoryChat from './HistoryChat.svelte';

	// Interface for grouped chats
	interface GroupedChats {
		today: ChatData[];
		yesterday: ChatData[];
		last7Days: ChatData[];
		last30Days: ChatData[];
		older: ChatData[];
	}

	let chatSearchTerm: string = $state('');
	const userLoggedIn: boolean = true; // TODO: change this to dynamic with Cookie's

	// Helper function to get date boundaries
	function getDateBoundaries() {
		const now = new Date();
		const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
		const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000);
		const last7Days = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
		const last30Days = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000);

		return {
			today: today,
			yesterday: yesterday,
			last7Days: last7Days,
			last30Days: last30Days
		};
	}

	// Function to group chats by time periods
	function groupChatsByTime(chats: ChatData[]): GroupedChats {
		const boundaries = getDateBoundaries();
		const grouped: GroupedChats = {
			today: [],
			yesterday: [],
			last7Days: [],
			last30Days: [],
			older: []
		};

		chats.forEach((chat) => {
			const chatTime = new Date(chat.last_message_at);

			if (chatTime >= boundaries.today) {
				grouped.today.push(chat);
			} else if (chatTime >= boundaries.yesterday) {
				grouped.yesterday.push(chat);
			} else if (chatTime >= boundaries.last7Days) {
				grouped.last7Days.push(chat);
			} else if (chatTime >= boundaries.last30Days) {
				grouped.last30Days.push(chat);
			} else {
				grouped.older.push(chat);
			}
		});

		// Sort each group by last_message_at (newest first)
		Object.keys(grouped).forEach((key) => {
			grouped[key as keyof GroupedChats].sort((a, b) => b.last_message_at - a.last_message_at);
		});

		return grouped;
	}

	// Filter and group chats
	let filteredAndGroupedChats = $derived.by(() => {
		let filtered = chats.filter((chat: ChatData) => chat.is_pinned === false);

		// Apply search filter if search term exists
		if (chatSearchTerm.trim() !== '') {
			filtered = filtered.filter((chat: ChatData) =>
				chat.title.toLowerCase().includes(chatSearchTerm.toLowerCase())
			);
		}

		return groupChatsByTime(filtered);
	});

	let pinnedChats: ChatResponse = $derived(
		chats
			.filter((chat: ChatData) => chat.is_pinned === true)
			.sort((a: ChatData, b: ChatData) => b.last_message_at - a.last_message_at)
	);

	// Updated to use new popup system
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

	async function deleteChat(id: string) {
		const url = 'http://localhost:3141';

		try {
			const modelResponse = await fetch(`${url}/v1/chats/${id}/`, {
				method: 'DELETE'
			});

			if (!modelResponse.ok) {
				throw new Error('Something happened during Record');
			}

			const index = chats.findIndex((chat: ChatData) => chat.id === id);
			if (index > -1) {
				chats.splice(index, 1);
			}
		} catch (error) {
			console.error('Error deleting chat:', error);
		}
	}

	async function patchChat(chat: ChatData, pin: boolean) {
		const url = 'http://localhost:3141';

		try {
			const response = await fetch(`${url}/v1/chats/${chat.id}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					is_pinned: pin
				})
			});

			if (!response.ok) {
				// Updated to use new popup system
				showConfirmationPopup({
					title: `Ups! Something ain't right..`,
					description: 'Some error happened while we tried to sync your data. Try again later.',
					primaryButtonName: 'Confirm',
					primaryButtonFunction: () => {
						// No additional action needed - popup will close automatically
					}
				});
				throw new Error("Couldn't sync Chats History");
			}

			chat.is_pinned = pin;
		} catch (error) {
			throw new Error(`${error}`);
		}
	}

	let activeContextMenuId = $state<string | null>(null);

	// Simple function to handle context menu opening
	function handleContextMenuOpen(chatId: string) {
		activeContextMenuId = chatId;
	}

	async function updateChatTitle(chatId: string, newTitle: string) {
		const url = 'http://localhost:3141';

		try {
			const response = await fetch(`${url}/v1/chats/${chatId}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					title: newTitle
				})
			});

			if (!response.ok) {
				throw new Error('Failed to update chat title');
			}

			console.log('Chat title updated successfully');
		} catch (error) {
			console.error('Error updating chat title:', error);
		}
	}

	// Complete example for your component
	function renameChat(chat: ChatData) {
		showRenamePopup({
			title: 'Rename Chat',
			description: 'Enter a new name for this chat:',
			inputValue: chat.title, // Pre-fill with current title
			inputPlaceholder: 'Enter chat name...',
			inputLabel: 'Chat Name',
			primaryButtonName: 'Rename',
			primaryButtonFunction: async () => {
				const currentState = get(popup);

				if (currentState.type === 'rename') {
					const newTitle = currentState.inputValue.trim();

					if (newTitle && newTitle !== chat.title) {
						// Update locally first for immediate UI feedback
						const oldTitle = chat.title;
						chat.title = newTitle;

						try {
							// Update on server
							await updateChatTitle(chat.id, newTitle);
							console.log(`Chat renamed from "${oldTitle}" to "${newTitle}"`);

							// Trigger reactivity if needed
							chats = [...chats];
						} catch (error) {
							// Revert on error
							chat.title = oldTitle;
							console.error('Failed to rename chat:', error);
						}
					}
				}
			}
		});
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

		// Only return sections that have chats
		return sections.filter((section) => section.chats.length > 0);
	});
</script>

{#if $isMobile && !sidebarCollapsed}
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
<div class="sidebar {sidebarCollapsed ? 'collapsed' : ''} {$isMobile ? 'isMobile' : ''}">
	<div class="head">
		<div class="title">Chat</div>
		<div class="newChatButton">
			<button onclick={newChat}>New Chat</button>
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
			<!-- TODO: change href link to account-settings -->
			<a href="#" class="account-button">
				<img src="https://placehold.co/100" alt="Profile Image" />
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
			background-color: var(--chat-background);
		}
	}

	.sidebar.isMobile {
		width: 256px;
		z-index: 99;
		position: absolute;
		left: 0;
		top: 0;
		height: 100%;
		border-right: 1px solid #88888822;
		background-color: var(--chat-background);
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

	.newChatButton button {
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

	.newChatButton:hover button {
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
