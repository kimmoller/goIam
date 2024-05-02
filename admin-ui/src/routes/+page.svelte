<script lang="ts">
	import { onMount } from 'svelte';

	const apiUrl = 'http://172.19.0.8:8083/';

	let data: Identity[];

	let isOpen = false;
	let dialogData: CreateIdentity = {
		firstName: '',
		lastName: '',
		email: ''
	};

	onMount(async () => {
		await getIdentities();
	});

	async function getIdentities() {
		try {
			const response = await fetch(apiUrl + 'identity');
			data = await response.json();
		} catch (error) {
			console.log(error);
		}
	}

	function toggleDialog() {
		isOpen = !isOpen;
	}

	async function createIdentity() {
		try {
			await fetch(apiUrl + 'identity', {
				method: 'POST',
				body: JSON.stringify(dialogData),
				headers: {
					'Content-type': 'application/json; charset=UTF-8'
				}
			});
			await getIdentities();
			toggleDialog();
			dialogData = {
				firstName: '',
				lastName: '',
				email: ''
			};
		} catch (error) {
			console.log(error);
		}
	}
</script>

<a href="/"><h1>IAM Admin UI</h1></a>
<label for="idSearch">Search for identity with id</label>
<input name="idSearch" />

<button on:click={() => toggleDialog()}>Create identity</button>

<dialog open={isOpen}>
	<p>Create identity</p>
	<label for="firstName">First name</label>
	<input name="firstName" bind:value={dialogData.firstName} />
	<label for="lastName">Last name</label>
	<input name="lastName" bind:value={dialogData.lastName} />
	<label for="email">Email</label>
	<input name="email" bind:value={dialogData.email} />
	<button on:click={() => toggleDialog()}>Cancel</button>
	<button on:click={() => createIdentity()} type="submit">Create</button>
</dialog>

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
