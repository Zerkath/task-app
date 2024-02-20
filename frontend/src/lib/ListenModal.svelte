<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import type { Message } from '$lib/types';
	let socket: WebSocket | null = null;

	export let id: string | undefined = undefined;
	let minimized: boolean = false;

	let hidden = true;
	let closed = false;

	let message: Message;

	const dispatch = createEventDispatcher();

	onMount(() => {
		if (!id) {
			return;
		}
		socket = new WebSocket(`ws://localhost:8080/task/listen/${id}`);
		socket.onopen = () => {
			console.log('WebSocket connection established');
			hidden = false;
		};
		socket.onmessage = (event) => {
			console.log('Message from server', event.data);
			try {
				message = JSON.parse(event.data) as Message;
			} catch (e) {
				console.error(e);
			}
		};
		socket.onclose = () => {
			console.log('WebSocket connection closed');
			closed = true;
		};
	});

	function close() {
		if (socket) {
			socket.close();
		}
		dispatch('close');
	}

	function getDefault(x: string | undefined): string {
		return x ? x : '';
	}

</script>

{#if minimized}
	<button class={`${getDefault(message?.status)}`} on:click={() => (minimized = false)}>
		{id.split('-')[0]} - {getDefault(message?.status)}</button
	>
{:else}
	<button class={`${getDefault(message?.status)}`} on:click={() => (minimized = true)}
		>Opened</button
	>
	<div class={`listen-modal ${getDefault(message?.status)}`} {hidden}>
		<p>Task {id}</p>

		{#if message}
			<p>{message.status}</p>
			<p>Completed at: {getDefault(message.completedAt)}</p>
		{/if}

		<button on:click={close}>Close</button>
		<button on:click={() => (minimized = true)}>Minimize</button>
	</div>
{/if}

<style>
	.listen-modal {
		padding: 20px;
		position: fixed;
		margin-top: 50px;
		left: 50%;
		margin-left: -300px;
		width: 600px;
		background-color: white;
		border: 1px solid black;
		align-items: center;
	}

	.failed {
		color: white;
		background: red;
	}

	.completed {
		color: white;
		background: green;
	}

	.queued {
		color: black;
		background: grey;
	}

	.running {
		color: white;
		background: blue;
	}

	p {
		font-weight: bold;
	}
</style>
