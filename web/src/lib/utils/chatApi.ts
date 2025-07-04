// lib/services/chatApi.ts
import { env } from '$env/dynamic/public';

export class ChatApiService {
	/**
	 * Delete a chat by ID
	 */
	static async deleteChat(id: string, session_token: string): Promise<void> {
		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/${id}/`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${session_token}`
				}
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
	static async updateChatPinStatus(
		chatId: string,
		isPinned: boolean,
		session_token: string
	): Promise<void> {
		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/${chatId}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${session_token}`
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
	static async updateChatTitle(
		chatId: string,
		newTitle: string,
		session_token: string
	): Promise<void> {
		try {
			const response = await fetch(`${env.PUBLIC_API_URL}/v1/chats/${chatId}/`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${session_token}`
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
