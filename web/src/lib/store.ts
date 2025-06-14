import { writable } from 'svelte/store';

export type PopupType = 'confirmation' | 'rename' | 'input';

export interface BasePopup {
	show: boolean;
	type: PopupType;
	title: string;
	description: string;
	primaryButtonName: string;
	primaryButtonFunction: () => void;
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
	inputValue: ''
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
