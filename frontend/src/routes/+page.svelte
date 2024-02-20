<script lang="ts">
	import { onMount } from 'svelte';
	import { pageStore } from '$lib/store';
	import type { Message } from '$lib/types';
	import ListenModal from '$lib/ListenModal.svelte';

	let currentPage = 0;
	let pageSize = 10;
	let pages = 0;

    let listeningTo: string[] = [];

	async function getPage(pageNum: number) {
        pages = Math.ceil($pageStore.count / pageSize);
		const url = `http://localhost:8080/task?size=${pageSize}&page=${pageNum}`;

		return fetch(url, {
			method: 'GET'
		})
			.then((response) => response.json())
			.then((data) => {
				if (typeof data === 'object' && !data.data) {
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

	async function changePage(offSet: number) {
		if (currentPage + offSet < 0) {
			return;
		}
		currentPage += offSet;
		await getPage(currentPage);
	}

	async function newTask() {
		const url = `http://localhost:8080/task`;

		return fetch(url, {
			method: 'POST'
		})
			.then((response) => response.json())
			.then((data) => {
                listenTo(data.id);
				getPage(currentPage);
			})
			.catch((err) => {
				console.error(err);
			});
	}

	async function deleteTask(id: string) {
		const url = `http://localhost:8080/task/${id}`;

		return fetch(url, {
			method: 'DELETE'
		})
			.then((_) => {
				getPage(currentPage);
			})
			.catch((err) => {
				console.error(err);
			});
	}

    function listenTo(id: string) {
        listeningTo = [...listeningTo, id]; // reassign to trigger reactivity
    }
    function removeId(id: string) {
        listeningTo = listeningTo.filter((i) => i !== id);
    }

	onMount(async () => {
		await getPage(currentPage);
        pages = Math.ceil($pageStore.count / pageSize);
	});
</script>

<h1>Demo application</h1>

<p>Page: {currentPage + 1} out of {pages}</p>

<button on:click={() => changePage(-1)} disabled={currentPage < 1}>{'<<'}</button>
<button on:click={() => changePage(1)} disabled={currentPage >= pages - 1}>{'>>'}</button>
<button on:click={() => getPage(currentPage)}>{'Refresh'}</button>

<input
	type="number"
	bind:value={pageSize}
	min="10"
	max="100"
	step="10"
	on:change={() => {
		getPage(currentPage);
	}}
/>

<button on:click={() => newTask()}>New task</button>

{#each listeningTo as id}
    <ListenModal bind:id={id} on:close={() => removeId(id)} />
{/each}

<table>
	<tr>
		<th>id</th>
		<th>status</th>
		<th>createdAt</th>
        <th>completedAt</th>
		<th>restarts</th>
		<th>actions</th>
	</tr>
	{#each $pageStore.data as item}
		<tr>
			<td class="centered">{item.id}</td>
			<td class={`centered ${item.status}`}>{item.status}</td>
			<td>{item.createdAt}</td>
            <td>{item.completedAt}</td>
			<td class="centered">{item.restarts}</td>
			<td class="centered">
				<button on:click={() => deleteTask(item.id)}>Delete</button>
				<button on:click={() => listenTo(item.id)}>Listen</button>
			</td>
		</tr>
	{/each}
</table>

<style lang="scss">
	table {
		width: 100%;
		border-collapse: collapse;
	}

	th,
	td {
		border: 1px solid black;
		padding: 8px;
		text-align: left;
	}

	tr:nth-child(even) {
		background-color: #f2f2f2;
	}

    .centered {
        text-align: center;
    }
</style>
