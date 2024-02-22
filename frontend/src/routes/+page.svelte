<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { pageStore, pageSettings } from '$lib/store';
	import type { Message, PageSettings } from '$lib/types';
	import { getPage, changePage, newTask, deleteTask, parseDefault } from '$lib/utils';
	import ListenModal from '$lib/ListenModal.svelte';

	let socket: WebSocket | undefined = undefined;

	let listeningTo: string[] = [];
	let expanded: string | undefined = undefined;

	let listenData: Message[] = [];

	let pagination: PageSettings = {
		pages: 0,
		currentPage: 0,
		pageSize: 10
	};

	pageSettings.subscribe((value) => {
		pagination = value;
	});

	function listenTo(id: string) {
        if (expanded === undefined) {
            expanded = id;
        }
		for (const x of listeningTo) {
			if (x === id) {
				return; // don't add duplicates
			}
		}
		listeningTo = [...listeningTo, id]; // reassign to trigger reactivity
		updateListening();
	}

	function removeId(id: string) {
        if (expanded === id) {
            expanded = undefined;
        }
		listeningTo = listeningTo.filter((i) => i !== id);
		updateListening();
	}

	onMount(async () => {
		await getPage();

		// Plain socket, that will accept arrays of uuids
		socket = new WebSocket('ws://localhost:8080/task/listen');

		socket.onopen = () => {
			console.log('Socket opened');
		};

		socket.onmessage = (event) => {
			const data = JSON.parse(event.data);
			listenData = data;
			console.log('Socket message', data);
		};

		socket.onclose = () => {
			console.log('Socket closed');
		};

        while (true) {
            // This is a simple hack to recheck the current page, this update should come from the socket instead
            // TODO: Implement a proper way to update the page
            await new Promise((resolve) => setTimeout(resolve, 5000));
            getPage();
        }
	});

	function updateListening() {
		if (socket) {
			socket.send(JSON.stringify(listeningTo));
		}
	}

	onDestroy(() => {
		if (socket) {
			socket.close();
		}
	});
</script>

<h1>Demo application</h1>

<p>Page: {pagination.currentPage + 1} out of {pagination.pages}</p>

<button on:click={() => changePage(-1)} disabled={pagination.currentPage < 1}>{'<<'}</button>
<button on:click={() => changePage(1)} disabled={pagination.currentPage >= pagination.pages - 1}
	>{'>>'}</button
>
<button on:click={() => getPage()}>{'Refresh'}</button>

<input
	type="number"
	bind:value={$pageSettings.pageSize}
	min="5"
	max="100"
	step="5"
	on:change={() => {
		getPage();
	}}
/>

<button on:click={() => newTask().then((id) => listenTo(id))}>New task</button>

<p>Awaiting updates from following tasks</p>

{#if listeningTo.length === 0}
	<span>None</span>
{/if}

{#each listenData as entry}
	<ListenModal
		bind:id={entry.id}
		message={entry}
        minimized={expanded !== entry.id}
		on:close={() => removeId(entry.id)}
		on:expand={() => (expanded = entry.id)}
		on:minimize={() => {
			if (expanded === entry.id) {
				expanded = undefined;
			}
		}}
	/>
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
			<td>{parseDefault(item.completedAt)}</td>
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
		margin-top: 20px;
		border-collapse: collapse;
	}

	th,
	td {
		border: 1px solid black;
		padding: 4px;
		text-align: left;
	}

	tr:nth-child(even) {
		background-color: #f2f2f2;
	}

	.centered {
		text-align: center;
	}
</style>
