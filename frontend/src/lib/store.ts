import { writable } from 'svelte/store';
import type { Page, PageSettings } from './types';

export const EMPTY_PAGE = {
    data: [],
    page: 0,
    count: 0
} as Page;

export const pageStore = writable(EMPTY_PAGE);

export const pageSettings = writable({
    pageSize: 10,
    currentPage: 0,
    pages: 0
} as PageSettings);

pageStore.subscribe((value) => {
    pageSettings.update((settings) => {
        settings.pages = Math.ceil(value.count / settings.pageSize);
        return settings;
    });
});
