// lib/services/chatApi.ts
import { PRIVATE_HOST_URL } from '$env/static/private';
import type { ChatData } from '$lib/types';


export class ChatApiService {
	/**
	 * Delete a chat by ID
	 */
	static async deleteChat(id: string): Promise<void> {
		try {
			const response = await fetch(`${PRIVATE_HOST_URL}/v1/chats/${id}/`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error('Something happened during deletion');
			}
		} catch (error) {
			console.error('Error deleting chat:', error);
			throw error;
		}
	}

	/**
	 * Update chat pin status
	 */
	static async updateChatPinStatus(chatId: string, isPinned: boolean): Promise<void> {
		try {
			const response = await fetch(`${PRIVATE_HOST_URL}/v1/chats/${chatId}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					is_pinned: isPinned
				})
			});

			if (!response.ok) {
				throw new Error("Couldn't sync chat pin status");
			}
		} catch (error) {
			console.error('Error updating chat pin status:', error);
			throw error;
		}
	}

	/**
	 * Update chat title
	 */
	static async updateChatTitle(chatId: string, newTitle: string): Promise<void> {
		try {
			const response = await fetch(`${PRIVATE_HOST_URL}/v1/chats/${chatId}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					title: newTitle
				})
			});

			if (!response.ok) {
				throw new Error('Failed to update chat title');
			}
		} catch (error) {
			console.error('Error updating chat title:', error);
			throw error;
		}
	}
}
