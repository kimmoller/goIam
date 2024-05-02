<script lang="ts">
	import { onMount } from 'svelte';

	let data: Identity[];

	async function getIdentities() {
		const response = await fetch('http://172.19.0.8:8083/identity');
		data = await response.json();
		console.log(data);
	}

	onMount(async () => {
		await getIdentities();
	});
</script>

<h1>IAM Admin UI</h1>
<label for="idSearch">Search for identity with id</label>
<input name="idSearch" />

{#if data !== undefined}
	<table>
        <tr>
            <th>First name</th>
            <th>Last name</th>
            <th>Email</th>
        </tr>
		{#each data as identity}
			<tr>
				<td><a href="/identity/{identity.id}">{identity.firstName}</a></td>
                <td>{identity.lastName}</td>
                <td>{identity.email}</td>
            </tr>
		{/each}
    </table>
{/if}
