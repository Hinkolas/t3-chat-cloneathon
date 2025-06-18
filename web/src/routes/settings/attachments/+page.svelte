<script lang="ts">
	import { env } from '$env/dynamic/public';
	import type { AttachmentData, AttachmentResponse } from '$lib/types.js';

	interface Props {
		data: {
			attachments: AttachmentData[];
			SESSION_TOKEN: string;
		};
	}

	let { data }: Props = $props();

	import { ExternalLink, FileText, Trash, Check } from '@lucide/svelte';

	let filteredAttachments: AttachmentData[] = $state(data.attachments || []);
	let selectedAttachments = $state(new Set<string>());
	let deleteError = $state<string | null>(null);
	let isDeleting = $state(false);
	let SESSION_TOKEN = $state(data.SESSION_TOKEN || '');

	// Check if all attachments are selected
	let allSelected = $derived(
		filteredAttachments.length > 0 && selectedAttachments.size === filteredAttachments.length
	);

	// Toggle individual attachment selection
	function toggleAttachment(id: string) {
		if (selectedAttachments.has(id)) {
			selectedAttachments.delete(id);
		} else {
			selectedAttachments.add(id);
		}
		selectedAttachments = new Set(selectedAttachments);
	}

	// Toggle all attachments
	function toggleSelectAll() {
		if (allSelected) {
			selectedAttachments.clear();
		} else {
			selectedAttachments = new Set(filteredAttachments.map((att) => att.id));
		}
		selectedAttachments = new Set(selectedAttachments);
	}

	// Delete selected attachments
	async function deleteSelected() {
		if (selectedAttachments.size === 0) return;

		isDeleting = true;
		deleteError = null;

		try {
			const deletePromises = Array.from(selectedAttachments).map(async (id) => {
				const delRes = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/${id}/`, {
					method: 'DELETE',
					headers: {
						Authorization: `Bearer ${SESSION_TOKEN}`
					}
				});
				if (!delRes.ok) throw new Error(`Failed to delete attachment ${id}`);
				return id;
			});

			// Wait for all deletions to complete
			const deletedIds = await Promise.all(deletePromises);

			// Remove deleted attachments from filteredAttachments
			filteredAttachments = filteredAttachments.filter(
				(attachment) => !deletedIds.includes(attachment.id)
			);

			// Clear selection
			selectedAttachments.clear();
			selectedAttachments = new Set(selectedAttachments);
		} catch (error) {
			console.error('Error deleting selected files:', error);
			deleteError = error instanceof Error ? error.message : 'Failed to delete selected files';
		} finally {
			isDeleting = false;
		}
	}

	async function deleteAttachment(id: string) {
		isDeleting = true;
		deleteError = null;

		try {
			// Delete the attachment using the ID
			const delRes = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/${id}/`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${SESSION_TOKEN}`
				}
			});
			if (!delRes.ok) throw new Error('Failed to delete attachment');

			// Remove from local state if successful
			filteredAttachments = filteredAttachments.filter((attachment) => attachment.id !== id);

			// Remove from selection if it was selected
			if (selectedAttachments.has(id)) {
				selectedAttachments.delete(id);
				selectedAttachments = new Set(selectedAttachments);
			}
		} catch (error) {
			console.error('Error removing file:', error);
			deleteError = error instanceof Error ? error.message : 'Failed to remove file';
		} finally {
			isDeleting = false;
		}
	}
</script>

<div class="container">
	<div class="title">Attachments</div>
	<div class="description">
		Manage your uploaded files and attachments. Note that deleting files here will remove them from
		the relevant threads, but not delete the threads. This may lead to unexpected behavior if you
		delete a file that is still being used in a thread.
	</div>

	{#if deleteError}
		<div class="error-message">
			{deleteError}
		</div>
	{/if}

	<div class="header">
		<div class="head">
			<button onclick={toggleSelectAll} disabled={isDeleting}>
				<div class="square" class:selected={allSelected}>
					{#if allSelected}
						<Check size="12" />
					{/if}
				</div>
				Select All
			</button>
		</div>
		<div class="tail">
			<button onclick={deleteSelected} disabled={selectedAttachments.size === 0 || isDeleting}>
				{#if isDeleting}
					Deleting...
				{:else}
					Delete Selected ({selectedAttachments.size})
				{/if}
			</button>
		</div>
	</div>
	<div class="attachments">
		{#if filteredAttachments.length === 0}
			<div class="empty-state">
				<div class="empty-title">No Attachments</div>
				<div class="empty-description">
					We are sorry, but you don't have any Attachments to display yet.
				</div>
			</div>
		{:else}
			{#each filteredAttachments as attachment, index}
				<div
					class="attachment"
					onclick={() => {
						if (!isDeleting) toggleAttachment(attachment.id);
					}}
					role="button"
					aria-label="Select Attachment"
					tabindex="0"
					onkeydown={(e) => {}}
				>
					<div class="square" class:selected={selectedAttachments.has(attachment.id)}>
						{#if selectedAttachments.has(attachment.id)}
							<Check size="12" />
						{/if}
					</div>
					<div class="preview">
						<FileText size="24" />
					</div>
					<div class="infos">
						<a
							onclick={(e) => {
								e.stopPropagation();
							}}
							href={attachment.src}
							target="_blank"
							class="title"
						>
							{attachment.name}
							<ExternalLink size="16" />
						</a>
						<div class="type">{attachment.type}</div>
					</div>
					<div class="buttons">
						<button
							onclick={(e) => {
								e.stopPropagation();
								deleteAttachment(attachment.id);
							}}
							disabled={isDeleting}
						>
							<Trash size="16" />
						</button>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		gap: 20px;
		flex: 1;
		height: 100%;
		padding-bottom: 0;
	}

	.title {
		color: #ffffff;
		font-size: 22px;
		font-weight: 800;
	}

	.description {
		font-size: 14px;
		font-weight: 500;
		color: var(--text);
	}

	.error-message {
		padding: 12px 16px;
		background-color: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.3);
		border-radius: 6px;
		color: #ef4444;
		font-size: 14px;
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
		font-weight: 600;
		letter-spacing: 0.24px;
		padding: 4px 16px;
		border-radius: 5px;
		color: var(--text);
		border: 1px solid #302029;
		background-color: #21141e;
		cursor: pointer;
		transition:
			background-color 0.15s ease-out,
			border 0.15s ease-out;
	}

	.header button:hover:not(:disabled) {
		background-color: #261923;
	}

	.header button:disabled {
		color: hsl(var(--secondary-foreground) / 0.6);
		cursor: not-allowed;
	}

	.square {
		min-width: 15px;
		min-height: 15px;
		border-radius: 4px;
		border: 1px solid #614052;
		display: flex;
		justify-content: center;
		align-items: center;
		transition: all 0.15s ease-out;
	}

	.square.selected {
		background-color: #8b5a6b;
		border-color: #8b5a6b;
		color: white;
	}

	.attachments {
		position: relative;
		height: 100%;
		width: 100%;
		display: flex;
		flex-direction: column;
		flex: 1;
		min-height: 0;
		border: 1px solid #302029;
		border-radius: 8px;
		overflow-y: auto;
	}

	.attachment {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 16px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
	}

	.attachment:hover {
		background-color: rgba(255, 255, 255, 0.05);
	}

	.attachment:not(:last-child) {
		border-bottom: 1px solid #302029;
	}

	.attachment .preview {
		min-width: 48px;
		min-height: 48px;
		border-radius: 4px;
		background-color: #302029;
		display: flex;
		justify-content: center;
		align-items: center;
		color: var(--text);
	}

	.attachment .infos {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
		gap: 2px;
		line-height: 1.3;
	}

	.attachment .infos .title {
		text-decoration: none;
		font-size: 14px;
		font-weight: 500;

		display: flex;
		align-items: center;
		gap: 8px;
		color: var(--text);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		max-width: 100%;
	}

	.attachment .infos .title:hover {
		text-decoration: underline;
	}

	.attachment .infos .type {
		font-size: 12px;
		color: var(--placeholder);
	}

	.attachment .buttons {
		margin-left: auto;
	}

	.attachment .buttons button {
		all: unset;
		border: 1px solid var(--button-border-danger);
		border-radius: 6px;
		background-color: var(--button-background-danger);

		cursor: pointer;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 8px;

		transition: background-color 0.1s ease-out;
	}

	.attachment .buttons button:hover:not(:disabled) {
		background-color: var(--button-hover-danger);
	}

	.attachment .buttons button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 40px 20px;
		text-align: center;
	}

	.empty-title {
		font-size: 16px;
		font-weight: 600;
		color: hsl(var(--foreground));
	}

	.empty-description {
		font-size: 14px;
		color: var(--text);
	}
</style>
