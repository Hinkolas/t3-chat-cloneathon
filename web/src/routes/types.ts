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
	provider: string;
	features: ModelFeatures;
	flags: ModelFlags;
}

export interface ModelsResponse {
	[modelId: string]: ModelData;
}