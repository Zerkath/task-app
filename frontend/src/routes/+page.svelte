<script lang="ts">
    import { onMount } from 'svelte';
    import { pageData } from '$lib';

    async function getPage(page: number) {
        const limit = (1+page) * 10;
        const url = `http://localhost:8080/task?limit=${limit}&offset=${page}`;

        return fetch(url, {
            method: "GET",
            mode: "cors",
        })
        .then(response => response.json())
        .then(data => {
            pageData.set(data);
        })
        .catch(err => {
            console.error(err);
            pageData.set([]);
        });
    }

    onMount(async () => {
        await getPage(0);
    });

</script>

<h1>Demo application</h1>

<table>
    <tr>
        <th>id</th>
        <th>status</th>
        <th>createdAt</th>
        <th>restarts</th>
    </tr>
    {#each $pageData as item}
        <tr>
            <td>{item.id}</td>
            <td>{item.status}</td>
            <td>{item.createdAt}</td>
            <td>{item.restarts}</td>
        </tr>
    {/each}
</table>
