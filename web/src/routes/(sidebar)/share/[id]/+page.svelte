<script lang="ts">
	import type { ChatData, MessageData } from '$lib/types';

	interface Props {
		data: {
			chat: ChatData;
		};
	}

	let { data }: Props = $props();

	import { ChevronDown, CircleX, FileText } from '@lucide/svelte';
	import MarkdownIt from 'markdown-it';
	import texmath from 'markdown-it-texmath';
	import katex from 'katex';
	import markdownItHighlightjs from 'markdown-it-highlightjs';
	import 'highlight.js/styles/github-dark.css';
	import { fade } from 'svelte/transition';

	let messages: MessageData[] = $state(data.chat.messages || []);
	let reasoningStates: Record<string, boolean> = $state({});

	$effect(() => {
		messages = data.chat.messages || [];
	});

	// Initialize markdown-it
	const md = new MarkdownIt({
		html: false,
		xhtmlOut: false,
		breaks: true,
		linkify: true,
		typographer: true
	})
		.use(markdownItHighlightjs)
		.use(texmath, {
			engine: katex,
			delimiters: 'dollars',
			katexOptions: { macros: { '\\RR': '\\mathbb{R}' } }
		});

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
</script>

{#if messages}
	<div class="chat-wrapper">
		<div class="chat">
			{#if !(messages.length == 0)}
				{#each messages as message}
					<div class="single-chat-container">
						<div
							class="single-chat {message.role == 'user' ? 'user' : 'assistant'}"
							aria-label="Copy Codeblock"
							role="button"
							tabindex="0"
							onclick={handleCopyClick}
							onkeydown={() => {}}
						>
							{#if message.reasoning}
								<div class="reasoning-box">
									<button
										onclick={() => {
											const messageKey = message.id || messages.indexOf(message).toString();
											reasoningStates[messageKey] = !reasoningStates[messageKey];
											reasoningStates = { ...reasoningStates };
										}}
										class="reasoning-button"
										><div
											class="chevron-icon"
											class:rotated={!reasoningStates[
												message.id || messages.indexOf(message).toString()
											]}
										>
											<ChevronDown size="14" />
										</div>
										Reasoning</button
									>
									{#if reasoningStates[message.id || messages.indexOf(message).toString()]}
										<div transition:fade={{ duration: 100 }} class="reasoning-text">
											{@html renderMarkdown(message.reasoning)}
										</div>
									{/if}
								</div>
							{/if}
							{#if message.status === 'done' || message.status === 'error'}
								{@html renderMarkdown(message.content)}
								{#if message.attachments && message.attachments.length > 0}
									<div class="attachments">
										{#each message.attachments as attachment}
											<a
												href={attachment.src}
												target="_blank"
												rel="noopener noreferrer"
												class="attachment-link"
											>
												{#if attachment.type.startsWith('image/')}
													<img src={attachment.src} alt={attachment.name} />
												{:else}
													<FileText size="24" />
													<div class="extension">{attachment.name}</div>
												{/if}
											</a>
										{/each}
									</div>
								{/if}
								{#if message.status === 'error'}
									<div class="error-message">
										<span class="error-icon"><CircleX size="16" /></span>
										<span class="error-text">An error occurred while processing this message.</span>
									</div>
								{/if}
							{:else if message.status === 'streaming'}
								{@html renderMarkdown(message.content)}
							{/if}
						</div>
					</div>
				{/each}
			{/if}
			<div class="chat-spacer"></div>
		</div>
	</div>
{/if}

<style>
	.chat-wrapper {
		width: 100%;
		height: 100%;
		display: flex;
		justify-content: center;
		overflow-y: auto;
	}

	.chat {
		height: auto;
		width: 100%;
		max-width: 768px;
		margin: 0 auto;
		padding: 40px 16px;
		display: flex;
		flex-direction: column;
		gap: 48px;
	}

	.chat-spacer {
		min-height: 130px;
		width: 100%;
	}

	.single-chat-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.single-chat {
		display: block;
		width: 100%;
		max-width: 100%;
		word-wrap: break-word;
		overflow-wrap: break-word;
		color: var(--text);
		line-height: 1.7;
	}

	.single-chat.user {
		background-color: #2b2430;
		box-shadow: 0 0 2px var(--border);
		border-radius: 10px;
		margin-left: auto;
		padding: 16px;
		width: fit-content;
	}

	.single-chat.assistant {
		width: 100%;
		max-width: 100%;
		padding: 8px 0;
	}

	.attachments {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: flex-end;
		width: max-content;
		flex-wrap: wrap;
		gap: 16px;
		max-width: 100%;
		color: var(--text);
		line-height: 1.7;
		margin-top: 8px;
		margin-left: auto;
		padding: 4px 6px;
		border-radius: 10px;
	}

	.attachment-link {
		position: relative;
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
		width: 100%;
		color: var(--text-light) !important;
		text-decoration: none;
		font-size: 14px;
		padding-top: 2px;
		border-radius: 8px;
		transition: background-color 0.1s ease;
	}

	.attachment-link:not(:has(img)) {
		padding: 16px;
		border: 1px solid var(--primary-border);
		background-color: var(--primary-disabled);
	}

	.attachment-link img {
		max-width: 700px;
		max-height: 400px;
		border-radius: 8px;
		object-fit: cover;
	}

	.extension {
		font-size: 12px;
		font-weight: 600;
		padding: 2px 4px;
		border-radius: 8px;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		color: var(--text);
		transition: color 0.1s ease;
	}

	.attachment-link:hover:not(:has(img)) {
		background-color: var(--primary-background-light);
	}

	.reasoning-box {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.reasoning-button {
		all: unset;
		font-size: 14px;
		display: flex;
		align-items: center;
		gap: 4px;
		padding: 6px 10px;
		padding-right: 12px;
		width: max-content;
		cursor: pointer;
		border-radius: 4px;
		transition: background-color 0.1s ease;
	}

	.reasoning-button:hover {
		background-color: var(--primary-disabled);
	}

	.reasoning-text {
		padding: 16px;
		border-radius: 8px;
		color: #c7c3cf;
		background-color: #1a1720;
	}

	.chevron-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		transition: transform 0.2s ease;
	}

	.chevron-icon.rotated {
		transform: rotate(-90deg);
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
</style>
