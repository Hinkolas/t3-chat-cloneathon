<script lang="ts">
	import type { ChatData, ModelsResponse, ModelData, MessageData } from '$lib/types';

	interface Props {
		data: {
			chat: ChatData;
			models: ModelsResponse;
		};
	}

	let { data }: Props = $props();

	import { ArrowUp, ChevronDown, Globe, Paperclip } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import ModelRow from '$lib/components/ModelRow.svelte';
	import SearchInput from '$lib/components/SearchInput.svelte';

	import MarkdownIt from 'markdown-it';
	import markdownItHighlightjs from 'markdown-it-highlightjs';
	import 'highlight.js/styles/github-dark.css';
	import { ChatApiService } from '$lib/utils/chatApi';

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let modelSelectionOpen = $state(false);
	let modelSearchTerm: string = $state('');
	let filteredModels: ModelsResponse = $state(data.models);
	let messages: MessageData[] = $state(data.chat.messages);
	let selectedModelKey: string = $state(Object.keys(data.models)[0]);

	$effect(() => {
		messages = data.chat.messages;
		filteredModels = data.models;
		modelSearchTerm = '';
		modelSelectionOpen = false;
	});

	// Initialize markdown-it
	const md = new MarkdownIt({
		html: false,
		xhtmlOut: false,
		breaks: true,
		linkify: true,
		typographer: true
	}).use(markdownItHighlightjs);

	// Custom render for code topbar
	const defaultFenceRenderer =
		md.renderer.rules.fence ||
		function (
			tokens: any,
			idx: any,
			options: any,
			env: any,
			renderer: { renderToken: (arg0: any, arg1: any, arg2: any) => any }
		) {
			return renderer.renderToken(tokens, idx, options);
		};

	md.renderer.rules.fence = function (
		tokens: { [x: string]: any },
		idx: string | number,
		options: any,
		env: any,
		renderer: any
	) {
		const token = tokens[idx];
		const info = token.info ? token.info.trim() : '';
		const langName = info ? info.split(/\s+/g)[0] : '';
		const displayLang = langName || 'text';

		// Get the original rendered code
		const originalCode = defaultFenceRenderer(tokens, idx, options, env, renderer);

		// Extract just the inner content without the <pre><code> wrapper
		const codeContent = token.content;

		// Create our custom structure
		return `<div class="code-block-wrapper">
			<div class="code-block-header">
				<span class="code-lang">${displayLang.toUpperCase()}</span>
				<button class="copy-code-btn" data-code="${encodeURIComponent(codeContent)}" type="button">
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
						<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
					</svg>
				</button>
			</div>
			${originalCode}
		</div>`;
	};

	function renderMarkdown(content: string): string {
		return md.render(content);
	}

	function handleCopyClick(event: Event) {
		const target = event.target as HTMLElement;
		const button = target.closest('.copy-code-btn') as HTMLButtonElement;

		if (button) {
			const encodedCode = button.getAttribute('data-code');
			if (encodedCode) {
				const code = decodeURIComponent(encodedCode);

				navigator.clipboard
					.writeText(code)
					.then(() => {
						// Visual feedback
						const svg = button.querySelector('svg');
						if (svg) {
							const originalSVG = svg.outerHTML;
							svg.outerHTML = `<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<polyline points="20,6 9,17 4,12"/>
						</svg>`;

							button.style.color = '#22c55e';

							setTimeout(() => {
								const newSvg = button.querySelector('svg');
								if (newSvg) {
									newSvg.outerHTML = originalSVG;
								}
								button.style.color = '';
							}, 1500);
						}
					})
					.catch((err) => {
						console.error('Failed to copy code:', err);
					});
			}
		}
	}

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

	async function sendMessage(message: string) {
		const url = 'http://localhost:3141';

		const userChat: MessageData = {
			id: '', //TODO: add id
			chat_id: data.chat.id,
			role: 'user',
			model: selectedModelKey,
			content: message,
			reasoning: '',
			created_at: 0,
			updated_at: 0
		};
		messages.push(userChat);

		let accumulatedContent = '';

		try {
			const response = await fetch(`${url}/v1/chats/${data.chat.id}/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					model: selectedModelKey,
					content: message,
					reasoning: 0
				})
			});

			if (!response.ok) {
				console.log('error');
				throw new Error('Failed to send message');
			}

			const assistantChat: MessageData = {
				id: '', //TODO: add id
				chat_id: data.chat.id,
				role: 'assistant',
				model: selectedModelKey,
				content: '',
				reasoning: '',
				created_at: 0,
				updated_at: 0
			};

			messages.push(assistantChat);

			const reader = response.body!.getReader();
			const decoder = new TextDecoder();
			let buffer = '';
			let currentEvent = '';

			const assistantChatIndex = messages.length - 1; // Store the index

			while (true) {
				const { done, value } = await reader.read();
				if (done) break;

				buffer += decoder.decode(value, { stream: true });

				const lines = buffer.split('\n');
				buffer = lines.pop() || '';

				for (const line of lines) {
					if (line.startsWith('event: ')) {
						currentEvent = line.slice(7).trim();
					} else if (line.startsWith('data: ')) {
						try {
							const jsonData = line.slice(6).trim();
							if (!jsonData) continue;

							const parsedData = JSON.parse(jsonData);

							if (currentEvent === 'message_delta' && parsedData.content) {
								accumulatedContent += parsedData.content;

								messages[assistantChatIndex] = {
									...messages[assistantChatIndex],
									content: accumulatedContent
								};
								messages = [...messages];
							} else if (currentEvent === 'message_end') {
								// TODO: handle Message End
							}
						} catch (parseError) {
							console.warn('Failed to parse JSON:', parseError, 'Line:', line);
						}
					}
				}
			}
		} catch (error) {
			console.log('Error:', error);
		}
	}

	onMount(() => {
		autoResize();
	});
</script>

{#if messages}
	<div class="chat-wrapper">
		<div class="chat">
			{#if !(messages.length == 0)}
				{#each messages as message}
					<!-- svelte-ignore a11y_click_events_have_key_events -->
					<div
						class="single-chat {message.role == 'user' ? 'user' : 'assistant'}"
						aria-label="Copy Codeblock"
						role="button"
						tabindex="0"
						onclick={handleCopyClick}
					>
						{@html renderMarkdown(message.content)}
					</div>
				{/each}
			{/if}
			<div class="chat-spacer"></div>
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
						<span>{data.models[selectedModelKey].title}</span>
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
		/* Remove margin-bottom */
	}

	.chat {
		height: auto; /* Change from height: 100% */
		width: 100%;
		max-width: 768px;
		margin: 0 auto;
		padding: 40px 16px;
		padding-bottom: 200px; /* Add this instead of margin-bottom */
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.chat-spacer {
		min-height: 130px;
		width: 100%;
	}

	.single-chat {
		display: block;
		width: 100%;
		max-width: 100%;
		word-wrap: break-word;
		overflow-wrap: break-word;
		color: hsl(var(--secondary-foreground));
		line-height: 1.7;
	}

	.single-chat.user {
		background-color: #2b2430;
		box-shadow: 0 0 2px #88888866;
		border-radius: 10px;
		margin-left: auto;
		padding: 16px;
		width: fit-content;
		max-width: 80%;
	}

	.single-chat.assistant {
		width: 100%;
		max-width: 100%;
		padding: 8px 0;
	}

	/* Markdown styling */
	.single-chat :global(*) {
		max-width: 100%;
		word-wrap: break-word;
		overflow-wrap: break-word;
	}

	.single-chat :global(h1),
	.single-chat :global(h2),
	.single-chat :global(h3),
	.single-chat :global(h4),
	.single-chat :global(h5),
	.single-chat :global(h6) {
		margin: 1.2em 0 0.6em 0;
		font-weight: 600;
		line-height: 1.3;
	}

	.single-chat :global(h1) {
		font-size: 1.5em;
	}
	.single-chat :global(h2) {
		font-size: 1.3em;
	}
	.single-chat :global(h3) {
		font-size: 1.1em;
	}

	.single-chat :global(p) {
		margin: 0.8em 0;
		font-size: 1rem;
		line-height: 1.6;
	}

	.single-chat :global(p:first-child) {
		margin-top: 0;
	}

	.single-chat :global(p:last-child) {
		margin-bottom: 0;
	}

	/* Code block styling with topbar */
	.single-chat :global(.code-block-wrapper) {
		margin: 1em 0;
		border-radius: 8px;
		overflow: hidden;
		border: 1px solid #333;
		background-color: #1a1a1a;
	}

	.single-chat :global(.code-block-header) {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background-color: #262626;
		padding: 8px 12px;
		border-bottom: 1px solid #333;
		font-size: 0.75rem;
	}

	.single-chat :global(.code-lang) {
		color: #a1a1aa;
		font-weight: 500;
		font-size: 0.7rem;
		letter-spacing: 0.05em;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
	}

	.single-chat :global(.copy-code-btn) {
		all: unset;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 6px;
		border-radius: 4px;
		cursor: pointer;
		color: #a1a1aa;
		transition: all 0.2s ease;
		background: transparent;
		border: 1px solid transparent;
	}

	.single-chat :global(.copy-code-btn:hover) {
		background-color: #404040;
		color: #e5e5e5;
		border-color: #555;
	}

	.single-chat :global(.copy-code-btn:active) {
		transform: scale(0.95);
	}

	.single-chat :global(.code-block-wrapper pre) {
		background-color: transparent;
		border-radius: 0;
		padding: 16px;
		overflow-x: auto;
		margin: 0;
		border: none;
		max-width: 100%;
		font-size: 0.875rem;
	}

	.single-chat :global(pre) {
		background-color: #1a1a1a;
		border-radius: 8px;
		padding: 12px;
		overflow-x: auto;
		margin: 1em 0;
		border: 1px solid #333;
		max-width: 100%;
		font-size: 0.875rem;
	}

	.single-chat :global(code) {
		background-color: #2a2a2a;
		padding: 2px 6px;
		border-radius: 4px;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-size: 0.875rem;
		word-break: break-all;
	}

	.single-chat :global(pre code) {
		background-color: transparent;
		padding: 0;
		word-break: normal;
	}

	.single-chat :global(.code-block-wrapper pre code) {
		background-color: transparent;
		padding: 0;
		word-break: normal;
	}

	.single-chat :global(blockquote) {
		border-left: 3px solid #666;
		padding-left: 12px;
		margin: 1em 0;
		font-style: italic;
		color: #ccc;
		background-color: rgba(255, 255, 255, 0.02);
		padding: 8px 12px;
		border-radius: 4px;
	}

	.single-chat :global(ul),
	.single-chat :global(ol) {
		padding-left: 20px;
		margin: 0.8em 0;
	}

	.single-chat :global(li) {
		margin: 0.3em 0;
		line-height: 1.5;
	}

	.single-chat :global(a) {
		color: #7c3aed;
		text-decoration: none;
		word-break: break-all;
	}

	.single-chat :global(a:hover) {
		text-decoration: underline;
	}

	.single-chat :global(table) {
		border-collapse: collapse;
		width: 100%;
		margin: 1em 0;
		font-size: 0.875rem;
		overflow-x: auto;
		display: block;
		white-space: nowrap;
	}

	.single-chat :global(thead),
	.single-chat :global(tbody),
	.single-chat :global(tr) {
		display: table;
		width: 100%;
		table-layout: fixed;
	}

	.single-chat :global(th),
	.single-chat :global(td) {
		border: 1px solid #444;
		padding: 6px 8px;
		text-align: left;
		word-wrap: break-word;
	}

	.single-chat :global(th) {
		background-color: #333;
		font-weight: 600;
	}

	.single-chat :global(hr) {
		border: none;
		border-top: 1px solid #444;
		margin: 1.5em 0;
	}

	.single-chat :global(strong) {
		font-weight: 600;
	}

	.single-chat :global(em) {
		font-style: italic;
	}

	/* Input styling */
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
		font-size: 0.875rem;
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
