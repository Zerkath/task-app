import { writable } from 'svelte/store';
import type { Message } from './types';

export const pageData = writable([] as Message[]);

