<script lang="ts">
	import type { ChatData, ModelsResponse, ModelData, MessageData } from '$lib/types';

	interface Props {
		data: {
			chat: ChatData;
			models: ModelsResponse;
		};
	}

	let { data }: Props = $props();

	import { ArrowUp, Brain, ChevronDown, FileText, Globe, Paperclip, X } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import ModelRow from '$lib/components/ModelRow.svelte';
	import SearchInput from '$lib/components/SearchInput.svelte';

	import MarkdownIt from 'markdown-it';
	import markdownItHighlightjs from 'markdown-it-highlightjs';
	import 'highlight.js/styles/github-dark.css';
	import { fade } from 'svelte/transition';
	import { addChatId, removeChatId } from '$lib/store';
	import { env } from '$env/dynamic/public';

	const iconSize = 16;

	let textarea: HTMLElement;
	let message = $state('');
	let modelSelectionOpen = $state(false);
	let modelSearchTerm: string = $state('');
	let filteredModels: ModelsResponse = $state(data.models || {});
	let messages: MessageData[] = $state(data.chat.messages);
	let selectedModelKey: string = $state(data.chat.model || Object.keys(data.models)[0] || 'Empty');

	let reasoningStates: Record<string, boolean> = $state({});
	let reasoningEnabled = $state(false);
	let webSearchEnabled = $state(false);
	let isStreaming = $derived(() => messages.some((m) => m.status === 'streaming'));

	let activeStreams = new Set<string>();
	let eventSources = new Map<string, EventSource>();

	$effect(() => {
		// When chat data changes, clean up previous streams
		eventSources.forEach((eventSource, streamId) => {
			eventSource.close();
		});
		eventSources.clear();
		activeStreams.clear();

		// Update the reactive state
		messages = data.chat.messages || [];
		filteredModels = data.models;
		modelSearchTerm = '';
		modelSelectionOpen = false;
		selectedModelKey = data.chat.model || Object.keys(data.models)[0];
	});

	$effect(() => {
		// Monitor messages for streaming status
		messages.forEach((message, index) => {
			if (
				message.status === 'streaming' &&
				message.stream_id &&
				!activeStreams.has(message.stream_id)
			) {
				activeStreams.add(message.stream_id);
				startStreamingForMessage(message.stream_id, index);
			}
		});
	});

	$effect(() => {
		// Clean up when component unmounts
		return () => {
			// Close all active EventSource connections
			eventSources.forEach((eventSource, streamId) => {
				eventSource.close();
			});
			eventSources.clear();
			activeStreams.clear();
		};
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

	interface UploadedFileWithId {
		file: File;
		id: string;
	}

	// svelte-ignore non_reactive_update
	let fileInput: HTMLInputElement;
	let uploadedFiles: UploadedFileWithId[] = $state([]);
	let uploadingFile: File | null = $state(null);
	let uploadError: string | null = $state(null);
	let isDragOver = $state(false);

	async function cancelStreaming(streamId: string) {
		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/streams/${streamId}/`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error(`Cancel failed: ${response.status} ${response.statusText}`);
			}
		} catch (error) {
			console.error('Error canceling streaming:', error);
		}
	}

	function validateFileType(
		file: File,
		selectedModel: ModelData
	): { isValid: boolean; errorMessage?: string } {
		const fileType = file.type.toLowerCase();
		const fileName = file.name.toLowerCase();

		// Always allow text files regardless of model features
		if (fileType.startsWith('text/')) {
			return { isValid: true };
		}

		// Check if model has PDF support
		const hasPdf = selectedModel.features.has_pdf;
		// Check if model has vision support
		const hasVision = selectedModel.features.has_vision;

		// If model has neither PDF nor vision support, reject all non-text files
		if (!hasPdf && !hasVision) {
			return {
				isValid: false,
				errorMessage: `The model "${selectedModel.title}" doesn't support file attachments.`
			};
		}

		// Check PDF files
		if (fileType === 'application/pdf' || fileName.endsWith('.pdf')) {
			if (!hasPdf) {
				return {
					isValid: false,
					errorMessage: `The model "${selectedModel.title}" doesn't support PDF files.`
				};
			}
			return { isValid: true };
		}

		// Check image files
		const imageTypes = ['image/jpeg', 'image/jpg', 'image/png'];
		const imageExtensions = ['.jpg', '.jpeg', '.png'];

		const isImageType = imageTypes.some((type) => fileType === type);
		const isImageExtension = imageExtensions.some((ext) => fileName.endsWith(ext));

		if (isImageType || isImageExtension) {
			if (!hasVision) {
				return {
					isValid: false,
					errorMessage: `The model "${selectedModel.title}" doesn't support image files.`
				};
			}
			return { isValid: true };
		}

		// If file type is not supported by any feature
		const supportedTypes = ['text files'];
		if (hasPdf) supportedTypes.push('PDF');
		if (hasVision) supportedTypes.push('images (JPG, PNG, JPEG)');

		return {
			isValid: false,
			errorMessage: `File type not supported. The model "${selectedModel.title}" only accepts: ${supportedTypes.join(', ')}.`
		};
	}

	// Helper function to set error with timeout
	function setUploadErrorWithTimeout(message: string) {
		uploadError = message;
		setTimeout(() => {
			uploadError = null;
		}, 2000);
	}

	// Updated uploadFile function
	async function uploadFile(file: File): Promise<UploadedFileWithId | null> {
		// Validate file type against selected model
		const selectedModel = data.models[selectedModelKey];
		const validation = validateFileType(file, selectedModel);

		if (!validation.isValid) {
			setUploadErrorWithTimeout(validation.errorMessage || 'File type not supported');
			return null;
		}

		// Set uploading state
		uploadingFile = file;
		uploadError = null;

		try {
			const formData = new FormData();
			formData.append('file', file);
			formData.append('chat_id', data.chat.id);

			const response = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/`, {
				method: 'POST',
				body: formData
			});

			if (!response.ok) {
				throw new Error(`Upload failed: ${response.status} ${response.statusText}`);
			}

			const result = await response.json();

			// Reset uploading state
			uploadingFile = null;

			// Create uploaded file with ID
			const uploaded: UploadedFileWithId = { file, id: result.id };
			uploadedFiles = [...uploadedFiles, uploaded];

			return uploaded;
		} catch (error) {
			console.error('Error uploading file:', error);
			uploadingFile = null;
			setUploadErrorWithTimeout(error instanceof Error ? error.message : 'Upload failed');
			return null;
		}
	}

	// Updated uploadMultipleFiles function
	async function uploadMultipleFiles(files: FileList | File[]): Promise<void> {
		const fileArray = Array.from(files);
		const selectedModel = data.models[selectedModelKey];

		for (const file of fileArray) {
			// Check if file is already uploaded (by name and size)
			const isDuplicate = uploadedFiles.some(
				(uploadedFile) =>
					uploadedFile.file.name === file.name && uploadedFile.file.size === file.size
			);

			if (!isDuplicate) {
				// Validate file type before upload
				const validation = validateFileType(file, selectedModel);

				if (!validation.isValid) {
					// Set error for the first invalid file and stop
					setUploadErrorWithTimeout(validation.errorMessage || 'File type not supported');
					break;
				}

				await uploadFile(file);
			}
		}
	}

	// Updated file select handler - handles multiple files
	async function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const files = target.files;

		if (files && files.length > 0) {
			// Upload all selected files
			await uploadMultipleFiles(files);

			// Reset the input
			target.value = '';
		}
	}

	// Function to trigger file input
	function triggerFileUpload() {
		fileInput?.click();
	}

	// Drag and drop handlers
	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
		isDragOver = true;
	}

	function handleDragLeave(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
		// Only set isDragOver to false if we're leaving the window entirely
		if (!event.relatedTarget) {
			isDragOver = false;
		}
	}

	// Updated drag and drop handler - handles multiple files
	async function handleDrop(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
		isDragOver = false;

		const files = event.dataTransfer?.files;
		if (files && files.length > 0) {
			// Upload all dropped files
			await uploadMultipleFiles(files);
		}
	}

	// Prevent default drag behaviors on the entire window
	function preventDefaults(event: DragEvent) {
		event.preventDefault();
		event.stopPropagation();
	}

	// Updated remove file function with server DELETE
	async function removeFile(index: number) {
		const uploadedFile = uploadedFiles[index];
		if (!uploadedFile) return;

		try {
			// Delete the attachment using the stored ID
			const delRes = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/${uploadedFile.id}/`, {
				method: 'DELETE'
			});
			if (!delRes.ok) throw new Error('Failed to delete attachment');

			// Remove from local state if successful
			uploadedFiles = uploadedFiles.filter((_, i) => i !== index);
			uploadError = null;
		} catch (error) {
			console.error('Error removing file:', error);
			setUploadErrorWithTimeout(error instanceof Error ? error.message : 'Failed to remove file');
		}
	}

	function clearAllFiles() {
		uploadedFiles = [];
		uploadError = null;
	}

	// Helper function to format file size
	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
	}

	// Helper function to get file type icon
	function getFileIcon(file: File): string {
		const type = file.type;
		if (type.startsWith('image/')) return '🖼️';
		if (type === 'application/pdf') return '📄';
		if (type.includes('document') || type.includes('word')) return '📝';
		if (type === 'text/plain') return '📄';
		return '📁';
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
		const tempMessage = message;

		const userChat: MessageData = {
			id: '', //TODO: add id
			chat_id: data.chat.id,
			role: 'user',
			model: selectedModelKey,
			status: 'done',
			stream_id: '',
			content: tempMessage,
			reasoning: '',
			created_at: 0,
			updated_at: 0,
			attachments: uploadedFiles.map((f) => ({
				id: f.id,
				name: f.file.name,
				src: `${env.PUBLIC_API_URL}/v1/attachments/${f.id}/`,
				type: f.file.type,
				created_at: Date.now()
			}))
		};
		messages.push(userChat);

		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/${data.chat.id}/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					model: selectedModelKey,
					content: tempMessage,
					reasoning_effort: reasoningEnabled ? 1024 : 0,
					attachments: uploadedFiles.map((f) => f.id)
				})
			});

			if (!response.ok) {
				console.error('Failed to send message');
				throw new Error('Failed to send message');
			}

			const res = await response.json();
			const streamId = res.stream_id;

			const assistantChat: MessageData = {
				id: '', //TODO: add id
				chat_id: data.chat.id,
				role: 'assistant',
				status: 'streaming',
				model: selectedModelKey,
				stream_id: streamId,
				content: '',
				reasoning: '',
				created_at: 0,
				updated_at: 0
			};

			message = '';
			messages.push(assistantChat);
			clearAllFiles(); // clear up
		} catch (error) {
			console.error('Error sending a Message:', error);
		}
	}

	function startStreamingForMessage(stream_id: string, messageIndex: number) {
		let accumulatedContent = '';
		let accumulatedReasoning = '';

		// Check if EventSource already exists
		if (eventSources.has(stream_id)) {
			console.warn(`EventSource already exists for stream ${stream_id}`);
			return;
		}

		const eventSource = new EventSource(`${env.PUBLIC_API_URL}/v1/streams/${stream_id}/`);

		// Store the EventSource instance
		eventSources.set(stream_id, eventSource);

		eventSource.onopen = () => {
			addChatId(data.chat.id);
		};

		eventSource.addEventListener('message_delta', (event) => {
			try {
				const data = JSON.parse(event.data);

				if (data.content) {
					accumulatedContent += data.content;
				}
				if (data.reasoning) {
					accumulatedReasoning += data.reasoning;
				}

				// Create a new messages array to trigger reactivity
				const newMessages = [...messages];
				newMessages[messageIndex] = {
					...newMessages[messageIndex],
					content: accumulatedContent,
					reasoning: accumulatedReasoning
				};
				messages = newMessages;
			} catch (error) {
				console.error('Error parsing message_delta:', error, 'Raw data:', event.data);
			}
		});

		eventSource.addEventListener('message_end', (event) => {
			try {
				// Create a new messages array and update status
				const newMessages = [...messages];
				newMessages[messageIndex] = {
					...newMessages[messageIndex],
					status: 'done'
				};
				messages = newMessages;

				// Clean up
				activeStreams.delete(stream_id);
				eventSources.delete(stream_id);
				eventSource.close();
				removeChatId(data.chat.id);
			} catch (error) {
				console.error('Error handling message_end:', error);
			}
		});

		eventSource.onerror = (error) => {
			console.error('EventSource error details:', {
				error,
				readyState: eventSource.readyState,
				url: eventSource.url,
				stream_id
			});
		};
	}

	function getAcceptAttribute(): string {
		const selectedModel = data.models[selectedModelKey];
		if (!selectedModel) return '';

		const acceptTypes = [
			'text/*',
			'.txt',
			'text/plain',
			'text/csv',
			'text/html',
			'text/markdown',
			'text/xml',
			'text/css',
			'text/javascript',
			'text/x-python'
		];

		if (selectedModel.features.has_pdf) {
			acceptTypes.push('.pdf', 'application/pdf');
		}

		if (selectedModel.features.has_vision) {
			acceptTypes.push('.jpg', '.jpeg', '.png', 'image/jpeg', 'image/png');
		}

		return acceptTypes.join(',');
	}

	onMount(() => {
		autoResize();

		// Add drag and drop event listeners to the entire window
		window.addEventListener('dragenter', preventDefaults);
		window.addEventListener('dragover', handleDragOver);
		window.addEventListener('dragleave', handleDragLeave);
		window.addEventListener('drop', handleDrop);

		// Cleanup function
		return () => {
			window.removeEventListener('dragenter', preventDefaults);
			window.removeEventListener('dragover', handleDragOver);
			window.removeEventListener('dragleave', handleDragLeave);
			window.removeEventListener('drop', handleDrop);
		};
	});
</script>

{#if messages}
	<input
		bind:this={fileInput}
		type="file"
		multiple
		style="display: none;"
		onchange={handleFileSelect}
		accept={getAcceptAttribute()}
	/>

	{#if isDragOver}
		<div class="drag-overlay" transition:fade={{ duration: 150 }}>
			<div class="drag-content">
				<div class="drag-icon">📁</div>
				<p>Drop files here to upload</p>
				<small>Multiple files supported</small>
			</div>
		</div>
	{/if}

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
							{#if message.status === 'done'}
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
<div class="input-wrapper">
	<div class="input-container">
		<!-- Show uploading file -->
		{#if uploadingFile}
			<div class="upload-container">
				<div class="uploading-file">
					<div class="file-item uploading">
						<span class="file-icon">{getFileIcon(uploadingFile)}</span>
						<div class="file-info">
							<span class="file-name">{uploadingFile.name}</span>
							<span class="file-size">{formatFileSize(uploadingFile.size)}</span>
						</div>
						<div class="upload-progress">
							<div class="spinner"></div>
						</div>
					</div>
				</div>
			</div>
		{/if}

		<!-- Show uploaded files -->
		{#if uploadedFiles.length > 0 && !uploadingFile}
			<div class="upload-container">
				<div class="uploaded-files">
					{#each uploadedFiles as uploaded, index}
						<div class="file-item uploaded">
							<span class="file-icon">{getFileIcon(uploaded.file)}</span>
							<div class="file-info">
								<span class="file-name">{uploaded.file.name}</span>
							</div>
							<button
								class="remove-file"
								onclick={() => removeFile(index)}
								aria-label="Remove file"
							>
								×
							</button>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Show upload error -->
		{#if uploadError}
			<div class="upload-container">
				<div class="upload-error">
					<span class="error-icon">⚠️</span>
					<span class="error-message">{uploadError}</span>
					<button class="dismiss-error" onclick={() => (uploadError = null)}>×</button>
				</div>
			</div>
		{/if}

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
					autoResize();
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
						class:active={reasoningEnabled}
						onclick={() => {
							reasoningEnabled = !reasoningEnabled;
						}}
					>
						<Brain size={iconSize} />
						Reasoning
					</button>
				{/if}
				{#if data.models[selectedModelKey].features.has_web_search}
					<button
						class="reasoning-button-feature"
						class:active={webSearchEnabled}
						onclick={() => {
							webSearchEnabled = !webSearchEnabled;
						}}
					>
						<Globe size={iconSize} />
						Search
					</button>
				{/if}
				<button
					onclick={triggerFileUpload}
					class:has-file={uploadedFiles.length > 0}
					disabled={!!uploadingFile}
				>
					<Paperclip size={iconSize} />
					{#if uploadedFiles.length > 0}
						Attach ({uploadedFiles.length})
					{:else}
						Attach
					{/if}
					{#if uploadingFile}
						<div class="button-spinner"></div>
					{/if}
				</button>
			</div>
			<div class="button-group">
				{#if isStreaming()}
					<button
						class="active"
						onclick={() => {
							const streamingMessage = messages.find(
								(m) => m.status === 'streaming' && m.stream_id
							);
							if (streamingMessage && streamingMessage.stream_id) {
								cancelStreaming(streamingMessage.stream_id);
								messages = messages.map((m) =>
									m.id === streamingMessage.id ? { ...m, status: 'done' } : m
								);
							}
						}}
						id="SendButton"
						aria-label="Cancel streaming"
					>
						<X size="20" />
					</button>
				{:else}
					<button
						class={message.length == 0 ? '' : 'active'}
						onclick={() => {
							sendMessage(message);
							autoResize();
						}}
						disabled={message.length == 0}
						id="SendButton"
						aria-label="Send message"
					>
						<ArrowUp size="20" />
					</button>
				{/if}
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
		color: hsl(var(--secondary-foreground));
		line-height: 1.7;
		margin-top: 8px;
		margin-left: auto;
		padding: 4px 6px;
		background-color: #2b2430;
		border-radius: 10px;
	}

	.attachment-link {
		position: relative;
		display: flex;
		flex-direction: row;
		justify-content: center;
		align-items: center;
		width: 100%;
		color: hsl(var(--secondary-foreground) / 0.8) !important;
		text-decoration: none;
		font-size: 14px;
		padding-top: 2px;
		border-radius: 8px;
		transition: background-color 0.1s ease;
	}

	.attachment-link:not(:has(img)) {
		padding: 16px;
		background-color: hsl(var(--primary) / 0.5);
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
		color: hsl(var(--secondary-foreground));
		transition: color 0.1s ease;
	}

	.attachment-link:hover:not(:has(img)) {
		background-color: hsl(var(--primary) / 0.6);
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
		background-color: hsl(var(--primary) / 0.2);
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

	.reasoning-button-feature.active {
		border: 1px solid hsl(var(--primary) / 0.3);
		background-color: hsl(var(--primary) / 0.3);
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

	/* File upload styles */
	.upload-container {
		border-radius: 8px;
	}

	.uploaded-files,
	.uploading-file {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		padding: 8px 12px;
	}

	.file-item {
		display: flex;
		align-items: center;
		gap: 12px;
		background-color: hsl(var(--primary) / 0.1);
		border: 1px solid hsl(var(--primary) / 0.3);
		border-radius: 8px;
		padding: 8px 12px;
		font-size: 14px;
		transition: all 0.2s ease;
		width: max-content;
	}

	.file-item.uploaded {
		background-color: hsl(var(--primary) / 0.15);
		border-color: hsl(var(--primary) / 0.4);
	}

	.file-item.uploading {
		background-color: rgba(255, 165, 0, 0.1);
		border-color: rgba(255, 165, 0, 0.3);
		animation: pulse 2s infinite;
	}

	.file-icon {
		font-size: 14px;
		min-width: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.file-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
		flex: 1;
		min-width: 0;
	}

	.file-name {
		color: hsl(var(--secondary-foreground));
		font-weight: 500;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		max-width: 120px;
	}

	.file-size {
		color: #888;
		font-size: 12px;
	}

	.remove-file {
		all: unset;
		cursor: pointer;
		color: #888;
		font-size: 18px;
		line-height: 1;
		padding: 4px;
		border-radius: 4px;
		transition: all 0.15s ease;
		min-width: 16px;
		height: 16px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.remove-file:hover {
		color: #ff6b6b;
		background-color: rgba(255, 107, 107, 0.1);
	}

	.upload-progress {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 24px;
	}

	.spinner,
	.button-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid #333;
		border-top: 2px solid hsl(var(--primary));
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.button-spinner {
		width: 12px;
		height: 12px;
		border-width: 1.5px;
		margin-left: 4px;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.7;
		}
	}

	.upload-error {
		display: flex;
		align-items: center;
		gap: 8px;
		background-color: rgba(255, 107, 107, 0.1);
		border: 1px solid rgba(255, 107, 107, 0.3);
		border-radius: 8px;
		padding: 12px;
		color: #ff6b6b;
		font-size: 14px;
	}

	.error-icon {
		font-size: 16px;
		min-width: 16px;
	}

	.error-message {
		flex: 1;
	}

	.dismiss-error {
		all: unset;
		cursor: pointer;
		color: #ff6b6b;
		font-size: 16px;
		line-height: 1;
		padding: 2px 4px;
		border-radius: 4px;
		transition: background-color 0.15s ease;
	}

	.dismiss-error:hover {
		background-color: rgba(255, 107, 107, 0.2);
	}

	button.has-file {
		background-color: hsl(var(--primary) / 0.2);
		border-color: hsl(var(--primary) / 0.4);
	}

	button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	button:disabled:hover {
		background-color: transparent;
	}

	.drag-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		pointer-events: none;
	}

	.drag-content {
		text-align: center;
		color: white;
		padding: 40px;
		border: 2px dashed hsl(var(--primary));
		border-radius: 12px;
		background-color: rgba(0, 0, 0, 0.5);
	}

	.drag-icon {
		font-size: 48px;
		margin-bottom: 16px;
	}

	.drag-content p {
		font-size: 18px;
		margin: 0;
		color: hsl(var(--secondary-foreground));
	}

	.drag-content small {
		display: block;
		margin-top: 8px;
		font-size: 14px;
		color: #888;
	}

	.uploaded-files {
		max-height: 200px;
		overflow-y: auto;
	}

	.uploaded-files::-webkit-scrollbar {
		width: 4px;
	}

	.uploaded-files::-webkit-scrollbar-track {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 2px;
	}

	.uploaded-files::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.3);
		border-radius: 2px;
	}

	.uploaded-files::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.5);
	}
</style>
