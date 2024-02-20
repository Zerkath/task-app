import { pageStore, pageSettings, EMPTY_PAGE } from './store';
import type { Message } from './types';

let currentPage = 0;
let pageSize = 10;

pageStore.subscribe((value) => {
    currentPage = value.page;
});

pageSettings.subscribe((value) => {
    pageSize = value.pageSize;
    currentPage = value.currentPage;
});

export function parseDefault(x: string | undefined): string {
    return x || '';
}

export async function getPage() {

    const url = `http://localhost:8080/task?size=${pageSize}&page=${currentPage}`;

    return fetch(url, {
        method: 'GET'
    })
        .then((response) => response.json())
        .then((data) => {
            if (typeof data === 'object' && !data.data) {
                pageStore.set(EMPTY_PAGE);
                return; // should check if the response is valid
            }
            pageStore.set(data);
        })
        .catch((err) => {
            console.error(err);
            pageStore.set({
                data: [] as Message[],
                page: 0,
                count: 0
            });
        });
}

export async function changePage(offSet: number) {
    if (currentPage + offSet < 0) {
        return;
    }

    pageSettings.update((settings) => {
        settings.currentPage += offSet;
        return settings;
    });

    await getPage();
}

export async function deleteTask(id: string) {

    const url = `http://localhost:8080/task/${id}`;

    return fetch(url, {
        method: 'DELETE'
    })
        .then((_) => {
            getPage();
        })
        .catch((err) => {
            console.error(err);
        });
}
export async function newTask(): Promise<string> {
    const url = `http://localhost:8080/task`;

    return fetch(url, {
        method: 'POST'
    })
        .then((response) => response.json())
        .then((data) => {
            let x = data.id;
            getPage();
            return x
        })
        .catch((err) => {
            console.error(err);
        });
}

