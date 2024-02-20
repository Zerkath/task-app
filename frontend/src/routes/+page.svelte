<script lang="ts">
	import { onMount } from 'svelte';
	import { pageStore, pageSettings } from '$lib/store';
    import type { PageSettings } from '$lib/types';
    import { getPage, changePage, newTask, deleteTask, parseDefault } from '$lib/utils'
	import ListenModal from '$lib/ListenModal.svelte';

    let listeningTo: string[] = [];

    let pagination: PageSettings = {
        pages: 0,
        currentPage: 0,
        pageSize: 10,
    }

    pageSettings.subscribe((value) => {
        pagination = value;
    });

    function listenTo(id: string) {
        listeningTo = [...listeningTo, id]; // reassign to trigger reactivity
    }
    function removeId(id: string) {
        listeningTo = listeningTo.filter((i) => i !== id);
    }

	onMount(async () => {
		await getPage();
	});
</script>

<h1>Demo application</h1>

<p>Page: {pagination.currentPage + 1} out of {pagination.pages}</p>

<button on:click={() => changePage(-1)} disabled={pagination.currentPage < 1}>{'<<'}</button>
<button on:click={() => changePage(1)} disabled={pagination.currentPage >= pagination.pages - 1}>{'>>'}</button>
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
