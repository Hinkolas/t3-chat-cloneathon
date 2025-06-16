<script lang="ts">
	import type { ModelsResponse, ModelData, MessageData } from '$lib/types';

	interface Props {
		data: {
			models: ModelsResponse;
		};
	}

	let { data }: Props = $props();

	import {
		ArrowUp,
		ChevronDown,
		Globe,
		Paperclip,
		Sparkles,
		Newspaper,
		Code,
		GraduationCap,
		Icon,
		Brain
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import ModelRow from '$lib/components/ModelRow.svelte';
	import SearchInput from '$lib/components/SearchInput.svelte';
	import { scale } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { refreshChatHistory } from '$lib/store';

	interface ButtonData {
		icon: typeof Icon;
		label: string;
		suggestions: string[];
	}

	const buttonData: Record<string, ButtonData> = {
		create: {
			icon: Sparkles,
			label: 'Create',
			suggestions: [
				'Write a creative story about space exploration',
				'Generate ideas for a mobile app',
				'Create a marketing campaign for a new product'
			]
		},
		explore: {
			icon: Newspaper,
			label: 'Explore',
			suggestions: [
				'What are the latest developments in AI?',
				'Explain quantum computing in simple terms',
				'What are the benefits of renewable energy?'
			]
		},
		code: {
			icon: Code,
			label: 'Code',
			suggestions: [
				'How to create a REST API with Node.js?',
				'Explain React hooks with examples',
				'Write a Python function to sort an array'
			]
		},
		learn: {
			icon: GraduationCap,
			label: 'Learn',
			suggestions: [
				'How does machine learning work?',
				'Explain the basics of blockchain technology',
				'What is the difference between AI and ML?'
			]
		}
	};

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let modelSelectionOpen = $state(false);
	let modelSearchTerm: string = $state('');
	let filteredModels: ModelsResponse = $state(data.models);

	let activeTab: string = $state('create');
	let currentSuggestions: string[] = $derived(buttonData[activeTab]?.suggestions || []);
	let selectedModelKey: string = $state(Object.keys(data.models)[0]);

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

	function changeModel(model: ModelData) {
		// Find the key by comparing model properties instead of object reference
		const modelKey = Object.entries(data.models).find(
			([key, modelData]) => modelData.name === model.name && modelData.title === model.title
		)?.[0];

		if (modelKey) {
			selectedModelKey = modelKey;
		}
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

	function setActiveTab(tab: string) {
		activeTab = tab;
	}

	async function sendMessage(message: string) {
		const tempMessage = message;
		const url = 'http://localhost:3141';

		try {
			const response = await fetch(`${url}/v1/chats/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					model: selectedModelKey,
					content: tempMessage,
					reasoning: 0
				})
			});

			if (!response.ok) {
				console.log('error');
				throw new Error('Failed to send message');
			}

			const res = await response.json();

			console.log(res);

			goto(`/chat/${res.chat_id}/`);
			refreshChatHistory();
		} catch (error) {
			console.log('Error:', error);
		}
	}

	onMount(() => {
		autoResize();
	});

	let reasoningOn: boolean = $state(false);
</script>

<div class="chat">
	{#if message.length == 0}
		<div class="placeholder">
			<div class="title">How can I help you?</div>
			<div class="buttons">
				{#each Object.entries(buttonData) as [key, button]}
					{@const Icon = buttonData[key].icon}
					<button class:active={activeTab === key} onclick={() => setActiveTab(key)}>
						<Icon size="16" />
						{button.label}
					</button>
				{/each}
			</div>
			<div class="suggestions">
				{#each currentSuggestions as suggestion, index}
					{#if index > 0}
						<div class="divider"></div>
					{/if}
					<button
						onclick={() => {
							message = suggestion;
						}}
					>
						{suggestion}
					</button>
				{/each}
			</div>
		</div>
	{/if}
</div>
<div class="input-wrapper">
	<div class="input-container">
		<textarea
			bind:this={textarea}
			bind:value={message}
			oninput={autoResize}
			placeholder="Type your message here..."
			name="message"
			id="Message"
			onkeydown={(e) => {
				if (e.key === 'Enter' && !e.shiftKey) {
					e.preventDefault();
					if (message.length == 0) return;
					sendMessage(message);
					message = ''; // TODO: handle in sendMessage with state
				}
			}}
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
						<span>{data.models[selectedModelKey].title}</span>
						<ChevronDown size={iconSize} />
					</button>
				</div>
				{#if data.models[selectedModelKey].features.has_reasoning}
					<button
						class="reasoning-button-feature"
						class:active={reasoningOn}
						onclick={() => {
							reasoningOn = !reasoningOn;
						}}
					>
						<Brain size={iconSize} />
						Reasoning
					</button>
				{/if}
				{#if data.models[selectedModelKey].features.has_web_search}
					<button
						class="reasoning-button-feature"
						class:active={reasoningOn}
						onclick={() => {
							reasoningOn = !reasoningOn;
						}}
					>
						<Globe size={iconSize} />
						Search
					</button>
				{/if}
				<button>
					<Paperclip size={iconSize} />
					Attach
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
	.chat {
		width: 100%;
		max-width: 768px;
		margin: 0 auto;
		overflow: hidden;
		padding: 40px 16px;
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.placeholder {
		width: 100%;
		padding: 24px;
		padding-top: 100px;
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		gap: 24px;
	}

	.placeholder .title {
		font-size: 30px;
		font-weight: 600;
		color: #ffffff;
	}

	.placeholder .buttons {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 12px;
	}

	.placeholder .buttons button {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 8px 20px;
		font-size: 14px;
		font-weight: 500;
		letter-spacing: 0.24px;
		background-color: #88888811;
		border: 1px solid #88888811;
		border-radius: 999px;
		gap: 12px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		color: hsl(var(--secondary-foreground));
	}
	.placeholder .buttons button.active {
		background-color: hsl(var(--primary) / 0.2);
		box-shadow: 0px 0px 2px hsl(var(--primary) / 0.3);
	}

	.placeholder .buttons button:hover {
		background-color: #88888822;
	}

	.placeholder .buttons button.active:hover {
		background-color: hsl(var(--primary) / 0.4);
	}

	.placeholder .suggestions {
		display: flex;
		flex-direction: column;
		width: 100%;
		gap: 4px;
	}

	.placeholder .suggestions .divider {
		width: 100%;
		height: 1px;
		background-color: #88888811;
	}

	.placeholder .suggestions button {
		all: unset;
		display: flex;
		justify-content: flex-start;
		align-items: center;
		padding: 8px 16px;
		font-size: 16px;
		letter-spacing: 0.24px;
		background-color: transparent;
		border: none;
		border-radius: 12px;
		gap: 12px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		color: hsl(var(--secondary-foreground));
	}

	.placeholder .suggestions button:hover {
		background-color: #88888811;
	}

	@media (max-width: 768px) {
		.placeholder .buttons {
			width: 100%;
			justify-content: space-around;
		}
		.placeholder .buttons button {
			border-radius: 8px;
			flex-direction: column;
			gap: 2px;
			padding: 12px 8px 8px 8px;
		}

		.placeholder .suggestions button:hover {
			background-color: transparent;
		}
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

	.reasoning-button-feature.active {
		background-color: hsl(var(--primary) / 0.5);
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

	#SendButton:hover:not(:disabled) {
		background-color: hsl(var(--primary) / 0.8);
	}

	#SendButton:disabled {
		cursor: not-allowed;
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
