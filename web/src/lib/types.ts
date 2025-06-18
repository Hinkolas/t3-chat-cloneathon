export interface ModelFeatures {
	has_vision?: boolean;
	has_pdf?: boolean;
	has_reasoning?: boolean;
	has_effort_control?: boolean;
	has_web_search?: boolean;
	has_fast?: boolean;
	has_image_generation?: boolean;
}

export interface ModelFlags {
	is_premium?: boolean;
	is_experimental?: boolean;
	is_key_required?: boolean;
	is_free?: boolean;
	is_new?: boolean;
	is_recommended?: boolean;
	is_open_source?: boolean;
}

export interface ModelData {
	title: string;
	description: string;
	icon: string;
	name: string;
	provider: string;
	features: ModelFeatures;
	flags: ModelFlags;
}

export interface ModelsResponse {
	[modelId: string]: ModelData;
}

export type ChatHistoryResponse = ChatHistoryData[];

export interface ChatHistoryData {
	id: string;
	title: string;
	is_pinned: boolean;
	last_message_at: number;
	created_at: number;
}

export type ChatResponse = ChatData;

export type AttachmentResponse = AttachmentData[];

export interface AttachmentData {
	id: string;
	user_id?: string;
	message_id?: string;
	name: string;
	src: string;
	type: string;
	created_at: number;
}

export interface MessageData {
	id: string;
	chat_id?: string;
	role: string;
	model: string;
	status: string;
	stream_id: string;
	content: string;
	reasoning?: string;
	created_at: number;
	updated_at: number;
	attachments?: AttachmentData[];
}

export interface ChatData {
	id: string;
	user_id: string;
	title: string;
	model: string;
	status: string;
	is_pinned: boolean;
	last_message_at: number;
	created_at: number;
	updated_at: number;
	shared_at: number;
	messages: MessageData[];
}

export type ProfileResponse = ProfileData;

export interface ProfileData {
	user_id: string;
	limit_standard: number;
	limit_premium: number;
	usage_standard: number;
	usage_premium: number;
	anthropic_api_key: string;
	openai_api_key: string;
	gemini_api_key: string;
	ollama_base_url: string;
	custom_user_name: string;
	custom_user_profession: string;
	custom_assistant_trait: string;
	custom_context: string;
}
