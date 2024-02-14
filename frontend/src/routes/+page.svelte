<script lang="ts">
	import { onMount } from 'svelte';
	import { pageData } from '$lib/store';
    import type { Message } from '$lib/types';
	import ListenModal from '$lib/ListenModal.svelte';

	let currentPage = 0;
	let pageSize = 10;

	let x: string | undefined = undefined;

	async function getPage(page: number) {
		const url = `http://localhost:8080/task?size=${pageSize}&page=${page}`;

		return fetch(url, {
			method: 'GET'
		})
			.then((response) => response.json())
			.then((data) => {
				pageData.set(data);
			})
			.catch((err) => {
				console.error(err);
				pageData.set([]);
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
				console.log(data);
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

	onMount(async () => {
		await getPage(currentPage);
	});
</script>

<h1>Demo application</h1>

<button on:click={() => changePage(-1)}>{'<<'}</button>
<button on:click={() => changePage(1)}>{'>>'}</button>
<label>Page: {currentPage + 1}</label>
<label>Page size:</label>
<input type="number" bind:value={pageSize} min="10" max="100" step="10" />
<button on:click={() => newTask()}>New task</button>

{#if x}
	<ListenModal bind:id={x} on:close={() => (x = undefined)} />
{/if}

<table>
	<tr>
		<th>id</th>
		<th>status</th>
		<th>createdAt</th>
		<th>restarts</th>
		<th>actions</th>
	</tr>
	{#each $pageData as item}
		<tr>
			<td>{item.id}</td>
			<td>{item.status}</td>
			<td>{item.createdAt}</td>
			<td>{item.restarts}</td>
			<td>
				<button on:click={() => deleteTask(item.id)}>Delete</button>
				<button on:click={() => (x = item.id)}>Listen</button>
			</td>
		</tr>
	{/each}
</table>

<style>
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
</style>
