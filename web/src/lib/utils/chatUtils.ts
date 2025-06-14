// lib/utils/chatUtils.ts
import type { ChatData } from '$lib/types';

export interface GroupedChats {
	today: ChatData[];
	yesterday: ChatData[];
	last7Days: ChatData[];
	last30Days: ChatData[];
	older: ChatData[];
}

interface DateBoundaries {
	today: Date;
	yesterday: Date;
	last7Days: Date;
	last30Days: Date;
}

/**
 * Get date boundaries for chat grouping
 */
export function getDateBoundaries(): DateBoundaries {
	const now = new Date();
	const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
	const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000);
	const last7Days = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
	const last30Days = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000);

	return {
		today,
		yesterday,
		last7Days,
		last30Days
	};
}

/**
 * Group chats by time periods
 */
export function groupChatsByTime(chats: ChatData[]): GroupedChats {
	const boundaries = getDateBoundaries();
	const grouped: GroupedChats = {
		today: [],
		yesterday: [],
		last7Days: [],
		last30Days: [],
		older: []
	};

	chats.forEach((chat) => {
		const chatTime = new Date(chat.last_message_at);

		if (chatTime >= boundaries.today) {
			grouped.today.push(chat);
		} else if (chatTime >= boundaries.yesterday) {
			grouped.yesterday.push(chat);
		} else if (chatTime >= boundaries.last7Days) {
			grouped.last7Days.push(chat);
		} else if (chatTime >= boundaries.last30Days) {
			grouped.last30Days.push(chat);
		} else {
			grouped.older.push(chat);
		}
	});

	// Sort each group by last_message_at (newest first)
	Object.keys(grouped).forEach((key) => {
		grouped[key as keyof GroupedChats].sort((a, b) => b.last_message_at - a.last_message_at);
	});

	return grouped;
}

/**
 * Filter chats by search term
 */
export function filterChatsBySearchTerm(chats: ChatData[], searchTerm: string): ChatData[] {
	if (searchTerm.trim() === '') {
		return chats;
	}
	
	return chats.filter((chat) =>
		chat.title.toLowerCase().includes(searchTerm.toLowerCase())
	);
}