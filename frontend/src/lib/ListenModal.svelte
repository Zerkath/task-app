<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
    import { parseDefault } from '$lib/utils';
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

</script>

{#if minimized}
	<button class={`${parseDefault(message?.status)}`} on:click={() => (minimized = false)}>
		{parseDefault(id).split('-')[0]} - {parseDefault(message?.status)}</button
	>
{:else}
	<button class={`${parseDefault(message?.status)}`} on:click={() => (minimized = true)}
		>Opened</button
	>
	<div class={`listen-modal ${parseDefault(message?.status)}`} {hidden}>
        <div class="controls">
            <button on:click={close}>Close</button>
            <button on:click={() => (minimized = true)}>Minimize</button>
        </div>

		<p>Task {id}</p>

		{#if message}
			<p>{message.status}</p>
			<p>Completed at: {parseDefault(message.completedAt)}</p>
		{/if}

	</div>
{/if}

<style lang="scss">
	.listen-modal {
        border-radius: 5px;
		padding: 8px 12px;
		position: fixed;
		margin-top: 50px;
		left: 50%;
		margin-left: -200px;
		width: 400px;
        box-shadow: 0 0 10px 0 black;
		border: 1px solid black;
		align-items: center;
	}

    .controls {
        display: flex;
        flex-direction: row-reverse;
        * {
            margin-left: 5px;
            background-color: #3f3f3f;
            border-radius: 5px;
            border: none;
            color: #dfdfdf;
        }
    }

</style>
