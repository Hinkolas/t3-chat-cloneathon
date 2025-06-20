<script lang="ts">
	import { env } from '$env/dynamic/public';
	import type { ModelsResponse, ModelData, MessageData, ProfileResponse } from '$lib/types';

	interface Props {
		data: {
			models: ModelsResponse;
			SESSION_TOKEN: string;
			profile: ProfileResponse;
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
	import { fade, scale } from 'svelte/transition';
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
	let filteredModels: ModelsResponse = $state(data.models || {});
	let profile = $state(data.profile || {});

	let activeTab: string = $state('create');
	let currentSuggestions: string[] = $derived(buttonData[activeTab]?.suggestions || []);
	let selectedModelKey: string = $state(Object.keys(data.models).at(-1) || 'Empty');
	let showPlaceholder: boolean = $state(true);

	let reasoningOn: boolean = $state(false);
	let webSearchEnabled = $state(false);

	// File upload related state
	let fileInput: HTMLInputElement;
	let uploadingFile: File | null = $state(null);
	let uploadError: string | null = $state(null);
	let isDragOver = $state(false);

	interface UploadedFileWithId {
		file: File;
		id: string;
	}

	let uploadedFiles: UploadedFileWithId[] = $state([]);

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

			const response = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${data.SESSION_TOKEN}`
				},
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

	// Updated remove file function
	async function removeFile(index: number) {
		const uploadedFile = uploadedFiles[index];
		if (!uploadedFile) return;

		try {
			// Delete the attachment using the stored ID
			const delRes = await fetch(`${env.PUBLIC_API_URL}/v1/attachments/${uploadedFile.id}/`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${data.SESSION_TOKEN}`
				}
			});
			if (!delRes.ok) throw new Error('Failed to delete attachment');

			// Remove from local state if successful
			uploadedFiles = uploadedFiles.filter((_, i) => i !== index);
			uploadError = null;
		} catch (error) {
			console.error('Error removing file:', error);
			uploadError = error instanceof Error ? error.message : 'Failed to remove file';
		}
	}

	function clearAllFiles() {
		uploadedFiles = [];
		uploadError = null;
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
		// Check if model is compatible with uploaded files
		if (!isModelCompatibleWithFiles(model, uploadedFiles)) {
			return; // Don't switch to incompatible model
		}

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

	async function sendMessage(messageText: string) {
		const tempMessage = messageText;
		showPlaceholder = false;

		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${data.SESSION_TOKEN}`
				},
				body: JSON.stringify({
					model: selectedModelKey,
					content: tempMessage,
					reasoning_effort: reasoningOn ? 1024 : 0,
					attachments: uploadedFiles.map((f) => f.id)
				})
			});

			if (!response.ok) {
				console.error('Failed to send message');
				throw new Error('Failed to send message');
			}

			const res = await response.json();

			message = '';
			goto(`/chat/${res.chat_id}/`);
			refreshChatHistory();
			clearAllFiles();
		} catch (error) {
			console.error('Error sending a Message:', error);
			showPlaceholder = true;
		}
	}

	function isModelCompatibleWithFiles(model: ModelData, files: UploadedFileWithId[]): boolean {
		if (model.provider == 'anthropic' && model.flags.is_key_required && profile.anthropic_api_key === '') {
			return false;
		}
		if (model.provider == 'gemini' && model.flags.is_key_required && profile.gemini_api_key === '') {
			return false;
		}
		if (model.provider == 'ollama' && model.flags.is_key_required && profile.ollama_base_url === '') {
			return false;
		}
		if (files.length === 0) return true;

		for (const uploadedFile of files) {
			const validation = validateFileType(uploadedFile.file, model);
			if (!validation.isValid) {
				return false;
			}
		}
		return true;
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

<div class="chat">
	{#if showPlaceholder && message.length === 0}
		<div class="placeholder" transition:fade={{ duration: 100 }}>
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
		<!-- Show uploading file -->
		{#if uploadingFile}
			<div class="upload-container">
				<div class="uploading-file">
					<div class="file-item uploading">
						<span class="file-icon">{getFileIcon(uploadingFile)}</span>
						<div class="file-info">
							<span class="file-name">{uploadingFile.name}</span>
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
							{#each Object.entries(filteredModels).reverse() as [modelId, model]}
								<ModelRow
									{model}
									{changeModel}
									disabled={!isModelCompatibleWithFiles(model, uploadedFiles)}
								/>
							{/each}
						</div>
					</div>
					<button
						disabled={Object.keys(data.models).length === 0}
						onclick={toggleModelSelection}
						class="selection-button non-selectable {modelSelectionOpen ? 'active' : ''}"
					>
						{#if Object.keys(data.models).length > 0}
							<span class="model-button">{data.models[selectedModelKey].title}</span>
							<ChevronDown size={iconSize} />
						{:else}
							<span class="model-button">No Models</span>
						{/if}
					</button>
				</div>
				{#if Object.keys(data.models).length > 0}
					{#if data.models[selectedModelKey].features.has_reasoning}
						<button
							class="reasoning-button-feature"
							class:active={reasoningOn}
							onclick={() => {
								reasoningOn = !reasoningOn;
							}}
						>
							<Brain size={iconSize} />
							<span>Reasoning</span>
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
							<span>Search</span>
						</button>
					{/if}
					<button
						onclick={triggerFileUpload}
						class:has-file={uploadedFiles.length > 0}
						disabled={!!uploadingFile}
					>
						<Paperclip size={iconSize} />
						<span>
							{#if uploadedFiles.length > 0}
								Attach ({uploadedFiles.length})
							{:else}
								Attach
							{/if}
						</span>
						{#if uploadingFile}
							<div class="button-spinner"></div>
						{/if}
					</button>
				{/if}
			</div>
			<div class="button-group">
				<button
					class={message.length == 0 ? '' : 'active'}
					onclick={() => {
						if (message.length == 0) return;
						sendMessage(message);
					}}
					disabled={message.length == 0 || Object.keys(data.models).length === 0}
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
		color: var(--white);
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
		background-color: var(--button-gray-background);
		border: 1px solid var(--button-gray-border);
		border-radius: 999px;
		gap: 12px;
		cursor: pointer;
		transition: background-color 0.15s ease-out;
		color: var(--text);
	}
	.placeholder .buttons button.active {
		background-color: var(--primary-background-light);
		border: 1px solid var(--primary-border);
	}

	.placeholder .buttons button:hover {
		background-color: var(--button-gray-hover);
	}

	.placeholder .buttons button.active:hover {
		background-color: var(--primary-background-light);
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
		background-color: var(--button-gray-background);
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
		color: var(--text);
	}

	.placeholder .suggestions button:hover {
		background-color: var(--button-gray-background);
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
		border: 1px solid var(--input-border);
		border-bottom: none;
		background-color: var(--input-background-secondary);
	}

	.input-container {
		display: flex;
		flex-direction: column;
		gap: 14px;

		border-top-left-radius: 12px;
		border-top-right-radius: 12px;
		padding-block: 12px;
		border: 1px solid var(--input-border);

		background-color: var(--input-background-pirmary);
	}

	textarea {
		all: unset;
		font-size: 15px;
		min-height: 1.5rem;
		max-height: 200px;
		overflow-y: auto;
		resize: none;
		padding-inline: 12px;
		color: var(--text);
	}

	textarea::-webkit-scrollbar {
		display: none;
	}

	textarea::placeholder {
		color: var(--placeholder);
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
		border: 1px solid var(--border);
		color: var(--text);
		transition: background-color 0.15s ease;
	}

	button:hover {
		background-color: var(--button-hover);
	}

	.reasoning-button-feature.active {
		background-color: var(--primary-background-light);
	}

	@media (max-width: 768px) {
		button span {
			display: none;
		}

		.model-button {
			display: block;
		}
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
		background-color: var(--model-selection-box-background);
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
		letter-spacing: 0.4px;
	}

	#SendButton {
		all: unset;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 8px;
		background-color: var(--primary-disabled);
		box-shadow: 0px 0px 2px var(--primary-background);
		border-radius: 8px;
		cursor: pointer;
		color: var(--text);
		transition: background-color 0.15s ease;
	}

	#SendButton.active {
		background-color: var(--primary-background-light);
	}

	#SendButton:hover:not(:disabled) {
		background-color: var(--primary-background);
	}

	#SendButton:disabled {
		cursor: not-allowed;
	}

	@media (hover: none) and (pointer: coarse) {
		button:hover {
			background-color: transparent;
		}

		#SendButton:hover {
			background-color: var(--primary-background-light);
		}

		#SendButton.active:hover {
			background-color: var(--primary-background-light);
		}
	}

	@media (max-width: 768px) {
		button:hover {
			background-color: transparent;
		}

		.selection-button.active {
			background-color: var(--button-gray-background);
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
		background-color: var(--primary-disabled);
		border: 1px solid var(--primary-border);
		border-radius: 8px;
		padding: 6px 12px;
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
		color: var(--text);
		font-weight: 500;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		max-width: 120px;
	}

	.remove-file {
		all: unset;
		cursor: pointer;
		color: var(--placeholder);
		font-size: 18px;
		line-height: 1;
		padding: 0px 1px 2px 3px; /* TODO: Fix this shit*/
		border-radius: 4px;
		transition: all 0.15s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.remove-file:hover {
		color: var(--button-hover-danger);
		background-color: var(--button-background-danger);
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
		border: 2px solid var(--border);
		border-top: 2px solid var(--primary-background);
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
		background-color: var(--primary-background-light);
		border-color: var(--primary-border);
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
		border: 2px dashed var(--primary-border);
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
		color: var(--text);
	}

	.drag-content small {
		display: block;
		margin-top: 8px;
		font-size: 14px;
		color: var(--placeholder);
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
