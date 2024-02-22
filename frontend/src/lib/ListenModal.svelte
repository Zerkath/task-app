<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { parseDefault } from '$lib/utils';
	import type { Message } from '$lib/types';

	export let id: string | undefined = undefined;
	export let message: Message;
	export let minimized: boolean;

	const dispatch = createEventDispatcher();
</script>

<section>
	{#if minimized}
		<button class={`${parseDefault(message?.status)}`} on:click={() => dispatch('expand')}>
			{parseDefault(id).split('-')[0]} - {parseDefault(message?.status)}</button
		>
	{:else}
		<button class={`${parseDefault(message?.status)}`} on:click={() => dispatch('minimize')}
			>Opened</button
		>
		<div class={`listen-modal ${parseDefault(message?.status)}`}>
			<div class="controls">
				<button on:click={() => dispatch('close')}>Close</button>
				<button on:click={() => dispatch('minimize')}>Minimize</button>
			</div>

			<p>Task {id}</p>

			{#if message}
				<p>{message.status}</p>
				<p>Completed at: {parseDefault(message.completedAt)}</p>
			{/if}
		</div>
	{/if}
</section>

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
