<script lang="ts">
	import type { ChatData, ModelsResponse, ModelData, MessageData } from '$lib/types';

	interface Props {
		data: {
			chat: ChatData;
			models: ModelsResponse;
		};
		sendMessage: (message: string) => void;
	}

	let { data, sendMessage }: Props = $props();

	import { ArrowUp, ChevronDown, Globe, Paperclip } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import ModelRow from '$lib/components/ModelRow.svelte';
	import SearchInput from '$lib/components/SearchInput.svelte';

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let modelSelectionOpen = $state(false);
	let modelSearchTerm: string = $state('');
	let filteredModels: ModelsResponse = $state(data.models);
	let messages: MessageData[] = $state(data.chat.messages);

	$effect(() => {
		messages = data.chat.messages;
		filteredModels = data.models;
		modelSearchTerm = '';
		modelSelectionOpen = false;
	});

	function toggleModelSelection() {
		if (modelSelectionOpen) {
			closeModelSelection();
		} else {
			modelSelectionOpen = true;
		}
	}

	function closeModelSelection() {
		modelSelectionOpen = false;

		setTimeout(() => {
			modelSearchTerm = '';
			filteredModels = data.models;
		}, 150);
	}

	function modelSearchFilter() {
		const filteredEntries = Object.entries(data.models).filter(
			([modelId, model]: [string, ModelData]) =>
				model.title.toLowerCase().includes(modelSearchTerm.toLowerCase())
		);

		filteredModels = Object.fromEntries(filteredEntries);
	}

	onMount(() => {
		autoResize();
	});

	function changeModel() {
		closeModelSelection();
	}

	function autoResize() {
		if (textarea) {
			textarea.style.height = 'auto';
			textarea.style.height = textarea.scrollHeight + 'px';
		}
	}

	function clickOutside(node: Element) {
		const handleClick = (event: Event) => {
			if (!node.contains(<Node>event.target)) {
				node.dispatchEvent(new CustomEvent('outsideclick'));
			}
		};

		document.addEventListener('click', handleClick, true);

		return {
			destroy() {
				document.removeEventListener('click', handleClick, true);
			}
		};
	}
</script>

{#if messages}
	<div class="chat-wrapper">
		<div class="chat">
			{#if !(messages.length == 0)}
				{#each messages as message}
					<div class="single-chat {message.role == 'user' ? 'user' : 'assistant'}">
						{message.content}
					</div>
				{/each}
			{/if}
		</div>
	</div>
{/if}
<div class="input-wrapper">
	<div class="input-container">
		<textarea
			bind:this={textarea}
			bind:value={message}
			oninput={autoResize}
			placeholder="Type your message here..."
			name="message"
			id="Message"
		></textarea>
		<div class="buttons">
			<div class="button-group">
				<div class="selection-container" use:clickOutside onoutsideclick={closeModelSelection}>
					<div class="selection-box {modelSelectionOpen ? 'visible' : ''}">
						<SearchInput
							bind:value={modelSearchTerm}
							onInputFunction={modelSearchFilter}
							placeholder="Search Models..."
						/>
						<div class="model-container">
							{#each Object.entries(filteredModels) as [modelId, model]}
								<ModelRow {model} {changeModel} />
							{/each}
						</div>
					</div>
					<button
						onclick={toggleModelSelection}
						class="selection-button non-selectable {modelSelectionOpen ? 'active' : ''}"
					>
						<span>Gemini 2.5 Flash</span>
						<ChevronDown size={iconSize} />
					</button>
				</div>
				<button>
					<Globe size={iconSize} />
					<span>Search</span>
				</button>
				<button>
					<Paperclip size={iconSize} />
				</button>
			</div>
			<div class="button-group">
				<button
					class={message.length == 0 ? '' : 'active'}
					onclick={() => {
						sendMessage(message);
						message = '';
					}}
					disabled={message.length == 0}
					id="SendButton"
				>
					<ArrowUp size="20" />
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	.chat-wrapper {
		width: 100%;
		height: 100%;
		display: flex;
		justify-content: center;
		overflow-y: auto;
	}

	.chat {
		height: 100%;
		width: 100%;
		max-width: 768px;
		margin: 0 auto;
		padding: 40px 16px;
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.single-chat {
		display: flex;
		width: max-content;
		max-width: 100%;
		flex-wrap: wrap;
		justify-content: flex-end;
		justify-self: flex-end;
		color: hsl(var(--secondary-foreground));
		line-height: 1.7;
	}

	.single-chat.user {
		background-color: #2b2430;
		box-shadow: 0 0 2px #88888866;
		border-radius: 10px;
		margin-left: auto;
		padding: 16px;
	}

	.input-wrapper {
		position: absolute;
		bottom: 0;
		left: 50%;
		transform: translateX(-50%);

		width: 100%;
		max-width: 768px;
		padding: 8px;
		padding-bottom: 0px;
		border-top-left-radius: 20px;
		border-top-right-radius: 20px;
		border: 1px solid #88888811;
		border-bottom: none;
		background-color: hsl(var(--chat-input-gradient));
	}

	.input-container {
		display: flex;
		flex-direction: column;
		gap: 14px;

		border-top-left-radius: 12px;
		border-top-right-radius: 12px;
		padding-block: 12px;
		border: 1px solid #88888811;

		background-color: var(--chat-input-background);
	}

	textarea {
		all: unset;
		font-size: 15px;
		min-height: 1.5rem;
		max-height: 200px;
		overflow-y: auto;
		resize: none;
		padding-inline: 12px;
		color: hsl(var(--secondary-foreground));
	}

	textarea::-webkit-scrollbar {
		display: none;
	}

	textarea::placeholder {
		color: #888888;
	}

	.input-container .buttons {
		width: 100%;
		display: flex;
		justify-content: space-between;
		padding-left: 6px;
		padding-right: 12px;
	}

	.button-group {
		display: flex;
		gap: 8px;
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

	.selection-container {
		position: relative;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.selection-box {
		position: absolute;
		bottom: 110%;
		left: 0;
		background-color: #0f0a0e;
		border-radius: 8px;
		display: flex;
		flex-direction: column;
		/* TODO: max-width adjusment -> overflow on mobile */
		min-width: 360px;

		opacity: 0;
		transform: translateY(10px) scale(0.95);
		pointer-events: none;
		transition:
			opacity 0.2s ease,
			transform 0.2s ease;
	}

	.selection-box.visible {
		opacity: 1;
		transform: translateY(0) scale(1);
		pointer-events: auto;
	}

	.model-container {
		display: flex;
		flex-direction: column;
		padding-inline: 8px;
		padding-bottom: 16px;
	}

	.selection-button {
		border: none;
		border-radius: 8px;
		font-size: 13px;
		font-weight: 600;
	}

	#SendButton {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 8px;
		background-color: hsl(var(--primary) / 0.2);
		box-shadow: 0px 0px 2px hsl(var(--primary));
		border-radius: 8px;
		cursor: pointer;
		color: hsl(var(--secondary-foreground));
		transition: background-color 0.15s ease;
	}

	#SendButton.active {
		background-color: hsl(var(--primary) / 0.4);
	}

	#SendButton:hover {
		background-color: hsl(var(--primary) / 0.8);
	}

	@media (hover: none) and (pointer: coarse) {
		button:hover {
			background-color: transparent;
		}

		#SendButton:hover {
			background-color: hsl(var(--primary) / 0.2);
		}

		#SendButton.active:hover {
			background-color: hsl(var(--primary) / 0.4);
		}
	}

	@media (max-width: 768px) {
		button:hover {
			background-color: transparent;
		}

		.selection-button.active {
			background-color: #88888811;
		}
	}
</style>
