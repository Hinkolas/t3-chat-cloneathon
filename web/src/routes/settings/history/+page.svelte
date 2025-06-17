<script lang="ts">
	import type { ChatHistoryData } from '$lib/types.js';
	import { ChatApiService } from '$lib/utils/chatApi.js';
	import { Upload, Trash2, Download, Pin } from '@lucide/svelte';

	let { data } = $props();

	let chats = $state(data.chats || []);
	let selectedChats = $state(new Set<string>());
	let selectAll = $state(false);

	// Reactive computed values
	$effect(() => {
		// Update selectAll state based on selection
		selectAll = chats.length > 0 && selectedChats.size === chats.length;
	});

	function toggleSelectAll() {
		if (selectAll) {
			selectedChats.clear();
		} else {
			selectedChats = new Set(chats.map((chat) => chat.id));
		}
	}

	function toggleChatSelection(chatId: string) {
		if (selectedChats.has(chatId)) {
			selectedChats.delete(chatId);
		} else {
			selectedChats.add(chatId);
		}
		// Trigger reactivity
		selectedChats = new Set(selectedChats);
	}

	async function deleteSelectedChats() {
		if (selectedChats.size === 0) return;

		try {
			// Delete all selected chats
			const deletePromises = Array.from(selectedChats).map((id) => ChatApiService.deleteChat(id));

			await Promise.all(deletePromises);

			// Remove deleted chats from the local state
			chats = chats.filter((chat) => !selectedChats.has(chat.id));

			// Clear selection
			selectedChats.clear();
		} catch (error) {
			console.error('Error deleting chats:', error);
		}
	}

	async function deleteAllChats() {
		if (chats.length === 0) return;

		try {
			// Delete all chats
			const deletePromises = chats.map((chat) => ChatApiService.deleteChat(chat.id));

			await Promise.all(deletePromises);

			// Clear all chats
			chats = [];
			selectedChats.clear();
		} catch (error) {
			console.error('Error deleting all chats:', error);
		}
	}

	function exportSelectedChats() {
		if (selectedChats.size === 0) return;

		const selectedChatData = chats.filter((chat) => selectedChats.has(chat.id));
		const dataStr = JSON.stringify(selectedChatData, null, 2);
		const dataBlob = new Blob([dataStr], { type: 'application/json' });

		const url = URL.createObjectURL(dataBlob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `chat-history-${new Date().toISOString().split('T')[0]}.json`;
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
		URL.revokeObjectURL(url);
	}

	function formatTimestamp(timestamp: number): string {
		const date = new Date(timestamp);

		const month = date.getMonth() + 1; // Months are 0-indexed
		const day = date.getDate();
		const year = date.getFullYear();

		let hours = date.getHours();
		const minutes = date.getMinutes();
		const seconds = date.getSeconds();
		const ampm = hours >= 12 ? 'PM' : 'AM';

		// Convert to 12-hour format
		hours = hours % 12;
		hours = hours ? hours : 12; // If 0, set to 12

		const formattedDate = `${month}/${day}/${year}, ${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')} ${ampm}`;
		return formattedDate;
	}
</script>

<div class="container">
	<div class="head">
		<div class="title">Message History</div>
		<div class="description">
			Save your history as JSON, or import someone else's. Importing will NOT delete existing
			messages
		</div>
		<div class="history-container">
			<div class="header">
				<div class="head">
					<button onclick={toggleSelectAll}>
						<div class="square" class:checked={selectAll}></div>
						<span>Select All</span>
					</button>
				</div>
				<div class="tail">
					<button disabled={selectedChats.size === 0} onclick={exportSelectedChats}>
						<Upload size="14" />
						<span>Export</span>
					</button>
					<button
						class="delete-button"
						disabled={selectedChats.size === 0}
						onclick={deleteSelectedChats}
					>
						<Trash2 size="14" />
						<span>Delete</span>
					</button>
					<button>
						<Download size="14" />
						<span>Import</span>
					</button>
				</div>
			</div>
			<div class="history-box">
				<div class="histories">
					{#each chats as chat (chat.id)}
						<button class="history" onclick={() => toggleChatSelection(chat.id)}>
							<div class="head">
								<div class="square" class:checked={selectedChats.has(chat.id)}></div>
								<div class="title">{chat.title}</div>
							</div>
							<div class="tail">
								{#if chat.is_pinned}
									<Pin color="#9E174D" size="16" />
								{/if}
								<div class="date">{formatTimestamp(chat.created_at)}</div>
							</div>
						</button>
					{/each}
				</div>
			</div>
		</div>
	</div>
	<div class="tail">
		<div class="title">Danger Zone</div>
		<div class="description">Permanently delete your account and all associated data.</div>
		<button class="main-delete-button" onclick={deleteAllChats}>Delete Chat History</button>
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.container > .head,
	.container > .tail {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.title {
		color: #ffffff;
		font-size: 22px;
		font-weight: 800;
	}

	.description {
		font-size: 14px;
		font-weight: 500;
		color: hsl(var(--secondary-foreground));
	}

	.history-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.header {
		display: flex;
		justify-content: space-between;
	}

	.header .head,
	.header .tail {
		display: flex;
		gap: 8px;
	}

	.header button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 8px;
		font-size: 12px;
		font-weight: 500;
		padding: 4px 16px;
		border-radius: 5px;
		color: hsl(var(--secondary-foreground));
		border: 1px solid #302029;
		background-color: #21141e;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.header button:hover:not(:disabled) {
		background-color: #261923;
	}

	.header button:disabled {
		color: hsl(var(--secondary-foreground) / 0.6);
		cursor: not-allowed;
	}

	.header button.delete-button {
		border: 1px solid #9e174d;
		background-color: #9e174d;
	}

	.header button.delete-button:disabled {
		border: 1px solid #621734;
		background-color: #621734;
	}

	@media (max-width: 768px) {
		.header button {
			padding: 10px 12px;
			font-size: 12px;
		}
		.header button span {
			display: none;
		}
	}

	.history-box {
		display: flex;
	}

	.histories {
		width: 100%;
		display: flex;
		flex-direction: column;
		border: 1px solid #26242d;
		border-radius: 5px;
		max-height: 240px;
		overflow-y: auto;
	}

	.history {
		all: unset;
		display: flex;
		justify-content: space-between;
		padding: 8px 16px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.history:hover {
		background-color: #2e272f;
	}

	.history:not(:last-child) {
		border-bottom: 1px solid #26242d;
	}

	.history .head {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.square {
		width: 15px;
		height: 15px;
		border-radius: 4px;
		border: 1px solid #614052;
		position: relative;
		transition: all 0.15s ease-out;
	}

	.square.checked {
		background-color: #9e174d;
		border-color: #9e174d;
	}

	.square.checked::after {
		content: '✓';
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		color: white;
		font-size: 10px;
		font-weight: bold;
	}

	.history .head .title {
		font-size: 14px;
		font-weight: 500;
		color: #f7f5f8;
	}

	.history .tail {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.history .tail .date {
		font-size: 12px;
		color: hsl(var(--secondary-foreground));
	}

	.main-delete-button {
		all: unset;
		width: max-content;
		font-size: 14px;
		font-weight: 600;
		color: #ffffff;
		padding: 8px 16px;
		cursor: pointer;
		border-radius: 8px;
		border: 1px solid var(--button-border-danger);
		background-color: var(--button-background-danger);
		transition: background-color 0.15s ease-out;
	}

	.main-delete-button:hover {
		background-color: var(--button-hover-danger);
	}
</style>
