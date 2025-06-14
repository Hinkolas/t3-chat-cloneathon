<script lang="ts">
	let { chats = $bindable(), sidebarCollapsed, newChat, toggleSidebar } = $props();

	import { fade } from 'svelte/transition';
	import { isMobile } from '$lib/deviceDetection';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import { X, Pin, PinOff } from '@lucide/svelte';
	import { popupModule } from '$lib/store';
	import type { ChatResponse, ChatData } from '$lib/types';

	console.log(chats);

	// Interface for grouped chats
	interface GroupedChats {
		today: ChatData[];
		yesterday: ChatData[];
		last7Days: ChatData[];
		last30Days: ChatData[];
		older: ChatData[];
	}

	let chatSearchTerm: string = $state('');

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

	function openPopup(id: string) {
		popupModule.update((currentModule) => {
			return {
				...currentModule,
				show: true,
				title: 'Delete Thread',
				description:
					'Are you sure you want to delete "Greeting Title"? This action cannot be undone.',
				primaryButtonName: 'Delete',
				primaryButtonFunction: () => {
					deleteChat(id);
					$popupModule.show = false;
				}
			};
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
				popupModule.update((currentModule) => {
					return {
						...currentModule,
						show: true,
						title: `Ups! Something ain't right..`,
						description: 'Some error happend while we tried to sync your data. Try again later.',
						primaryButtonName: 'Confirm',
						primaryButtonFunction: () => {
							$popupModule.show = false;
						}
					};
				});
				throw new Error("Couldn't sync Chats History");
			}

			chat.is_pinned = pin;
		} catch (error) {
			throw new Error(`${error}`);
		}
	}
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
		<div class="chats-container">
			<!-- Pinned Chats -->
			{#if pinnedChats.length > 0}
				<div class="day-title">
					<Pin size="14" />
					Pinned
				</div>
				<div class="chats">
					{#each pinnedChats as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, false)}>
									<PinOff size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Today -->
			{#if filteredAndGroupedChats.today.length > 0}
				<div class="day-title">Today</div>
				<div class="chats">
					{#each filteredAndGroupedChats.today as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, true)}>
									<Pin size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Yesterday -->
			{#if filteredAndGroupedChats.yesterday.length > 0}
				<div class="day-title">Yesterday</div>
				<div class="chats">
					{#each filteredAndGroupedChats.yesterday as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, true)}>
									<Pin size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Last 7 Days -->
			{#if filteredAndGroupedChats.last7Days.length > 0}
				<div class="day-title">Last 7 Days</div>
				<div class="chats">
					{#each filteredAndGroupedChats.last7Days as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, true)}>
									<Pin size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Last 30 Days -->
			{#if filteredAndGroupedChats.last30Days.length > 0}
				<div class="day-title">Last 30 Days</div>
				<div class="chats">
					{#each filteredAndGroupedChats.last30Days as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, true)}>
									<Pin size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}

			<!-- Older -->
			{#if filteredAndGroupedChats.older.length > 0}
				<div class="day-title">Older</div>
				<div class="chats">
					{#each filteredAndGroupedChats.older as chat}
						<div class="chat">
							<span>{chat.title}</span>
							<div class="buttons">
								<button onclick={() => patchChat(chat, true)}>
									<Pin size="14" />
								</button>
								<button onclick={() => openPopup(chat.id)}>
									<X size="14" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
	<div class="foot">Login</div>
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

	.chats-container {
		flex: 1;
		width: 100%;
		min-height: 0;
		display: flex;
		flex-direction: column;
		gap: 16px;
		overflow-y: auto;
		padding-inline: 16px;
	}

	/* WebKit scrollbar styling (Chrome, Edge) */
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

	.chat {
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
		color: hsl(var(--secondary-foreground));
		font-size: 14px;
		font-weight: 500;
		white-space: nowrap;
		text-overflow: ellipsis;
		overflow: hidden;
		flex: 1;
		min-width: 0;
	}

	.chat:hover {
		background-color: var(--sidebar-chat-hover);
	}

	.buttons {
		position: relative;
		top: 50%;
		left: 100%;
		transform: translateY(-50%);
		flex-shrink: 0;

		padding-right: 8px;
		display: flex;
		justify-content: flex-end;
		gap: 4px;
		transition: left 0.15s ease;
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
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease-out;
	}

	.buttons button:hover {
		background-color: hsl(var(--primary) / 0.5);
	}

	.chat:hover .buttons {
		left: 0;
	}
</style>
