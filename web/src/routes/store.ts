import { writable } from "svelte/store";
export const popupModule = writable({
    show: false,
    title: '',
    description: '',
    primaryButtonName: '',
    primaryButtonFunction: () => {}
});