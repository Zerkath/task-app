import { writable } from 'svelte/store';
import type { Page } from './types';

export const pageStore = writable({
    data: [],
    page: 0,
    count: 0
} as Page);

