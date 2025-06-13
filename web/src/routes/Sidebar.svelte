<script lang="ts">
	let { chats = $bindable(), sidebarCollapsed, newChat, toggleSidebar } = $props();

	import { fade } from 'svelte/transition';
	import { isMobile } from './deviceDetection';
	import SearchInput from './SearchInput.svelte';
	import { X } from '@lucide/svelte';
	import { popupModule } from './store';
	import type { ChatData } from './types';

	let chatSearchTerm = $state('');

	function openPopup(id: string) {
		popupModule.update((currentModule) => {
			return {
				...currentModule, // Keep existing properties
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

			// Direktes Entfernen aus dem reaktiven Array
			const index = chats.findIndex((chat: ChatData) => chat.id === id);
			if (index > -1) {
				chats.splice(index, 1);
			}
		} catch (error) {
			console.error('Error deleting chat:', error);
		}
	}

	function chatSearchFilter() {}
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
			<SearchInput
				bind:value={chatSearchTerm}
				onInputFunction={chatSearchFilter}
				placeholder="Search your threads..."
			/>
		</div>
		<div class="chats-container">
			<div class="day-title">Today</div>
			<div class="chats">
				{#each chats as chat}
					<div class="chat">
						<span>{chat.title}</span>
						<div class="buttons">
							<button
								onclick={() => {
									openPopup(chat.id);
								}}
							>
								<X size="14" />
							</button>
						</div>
					</div>
				{/each}
			</div>
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
		padding-block: 20px;
		background-color: var(--sidebar-background);
		transition: margin 0.15s ease-in-out;

		display: flex;
		flex-direction: column;
		justify-content: space-between;
		align-items: center;
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
		width: 100%;
		display: flex;
		flex-direction: column;

		padding-inline: 16px;
	}

	.day-title {
		font-size: 12px;
		font-weight: 600;
		color: var(--secondary);
		text-shadow: 0 0 4px hsl(var(--primary) / 0.2);
	}

	.chats {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.chat {
		position: relative;
		padding: 8px 8px;
		border-radius: 8px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		overflow: hidden;
	}

	.chat span {
		color: hsl(var(--secondary-foreground));
		font-size: 14px;
		font-weight: 500;
	}

	.chat:hover {
		background-color: var(--sidebar-chat-hover);
	}

	.buttons {
		position: absolute;
		top: 50%;
		left: 100%;
		transform: translateY(-50%);
		width: 100%;

		padding-right: 8px;
		display: flex;
		justify-content: flex-end;
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
