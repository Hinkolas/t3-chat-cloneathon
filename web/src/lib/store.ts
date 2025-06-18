import { writable } from 'svelte/store';

export type PopupType = 'confirmation' | 'rename' | 'input';

export interface BasePopup {
	show: boolean;
	type: PopupType;
	title: string;
	description: string;
	primaryButtonName: string;
	primaryButtonFunction: () => void;
	onConfirmTitle: string;
	onConfirmDescription: string;
	secondaryButtonName?: string;
	secondaryButtonFunction?: () => void;
	inputValue?: string;
	inputPlaceholder?: string;
	inputLabel?: string;
	inputType?: 'text' | 'email' | 'password' | 'number';
}

export interface ConfirmationPopup extends BasePopup {
	type: 'confirmation';
}

export interface RenamePopup extends BasePopup {
	type: 'rename';
	inputValue: string;
}

export interface InputPopup extends BasePopup {
	type: 'input';
	inputValue: string;
}

export type PopupData = ConfirmationPopup | RenamePopup | InputPopup;

const initialState: PopupData = {
	show: false,
	type: 'confirmation',
	title: '',
	description: '',
	primaryButtonName: '',
	primaryButtonFunction: () => {},
	inputValue: '',
	onConfirmTitle: '',
	onConfirmDescription: ''
};

export const popup = writable<PopupData>(initialState);

// Helper functions to show different types of popups
export const showConfirmationPopup = (data: Omit<ConfirmationPopup, 'show' | 'type'>) => {
	popup.set({
		...data,
		show: true,
		type: 'confirmation'
	});
};

export const showRenamePopup = (data: Omit<RenamePopup, 'show' | 'type'>) => {
	popup.set({
		...data,
		show: true,
		type: 'rename',
		inputValue: data.inputValue || ''
	});
};

export const showInputPopup = (data: Omit<InputPopup, 'show' | 'type'>) => {
	popup.set({
		...data,
		show: true,
		type: 'input',
		inputValue: data.inputValue || '',
		inputType: data.inputType || 'text'
	});
};

export const hidePopup = () => {
	popup.update((state) => ({ ...state, show: false }));
};

// Handling Sidebar State
export interface SidebarData {
	collapsed: boolean;
	refresh: boolean;
	chatIds: string[];
}
const initialSidebarState: SidebarData = {
	collapsed: false,
	refresh: false,
	chatIds: []
};

export let sidebarState = writable<SidebarData>(initialSidebarState);

export const addChatId = (chatId: string) => {
	sidebarState.update((state) => ({
		...state,
		chatIds: [...state.chatIds, chatId]
	}));
};

export const removeChatId = (chatId: string) => {
	sidebarState.update((state) => ({
		...state,
		chatIds: state.chatIds.filter((id) => id !== chatId)
	}));
};

export const toggleSidebar = () => {
	sidebarState.update((state) => ({ ...state, collapsed: !state.collapsed }));
};

export const closeSidebar = () => {
	sidebarState.update((state) => ({ ...state, collapsed: true }));
};

export const openSidebar = () => {
	sidebarState.update((state) => ({ ...state, collapsed: false }));
};

export const refreshChatHistory = () => {
	sidebarState.update((state) => ({ ...state, refresh: true }));
};

// Handling Notification State
export interface NotificationData {
	show: boolean;
	is_error: boolean;
	title: string;
	description: string;
}
const initialNotificationState: NotificationData = {
	show: false,
	is_error: false,
	title: 'Test Notification',
	description: 'This is a test notification message.This is a test notification message.'
};

export let notificationState = writable<NotificationData>(initialNotificationState);

export const showNotification = (data: Omit<NotificationData, 'show'>) => {
	notificationState.set({
		...data,
		show: true
	});
};

export const hideNotification = () => {
	notificationState.update((state) => ({ ...state, show: false }));
};
