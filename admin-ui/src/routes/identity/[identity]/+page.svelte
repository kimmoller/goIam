<script lang="ts">
	import { onMount } from 'svelte';

	const apiUrl = 'http://172.19.0.8:8083/';

	let identity: ExtendedIdentity;
	const identityId = window.location.pathname.slice(-1);

	let isOpen = false;
	let dialogData: CreateMembership = {
		identityId: identityId,
		groupId: '',
		enabledAt: new Date().toISOString(),
		disabledAt: null
	};

	onMount(async () => {
		await getIdentity();
	});

	async function getIdentity() {
		const response = await fetch(apiUrl + 'identity/' + identityId);
		identity = await response.json();
	}

	function toggleDialog() {
		isOpen = !isOpen;
	}

	async function createMembership() {
		try {
			await fetch(apiUrl + 'membership', {
				method: 'POST',
				body: JSON.stringify(dialogData),
				headers: {
					'Content-type': 'application/json; charset=UTF-8'
				}
			});
			await getIdentity();
			toggleDialog();
			dialogData = {
				identityId: identityId,
				groupId: '',
				enabledAt: new Date().toISOString(),
				disabledAt: null
			};
		} catch (error) {
			console.log(error);
		}
	}
</script>

<a href="/"><h1>IAM Admin UI</h1></a>

{#if identity !== undefined}
	<div>
		<h2>{identity.firstName} {identity.lastName}</h2>
		<h3>{identity.email}</h3>
	</div>
	<h2>Accounts</h2>
	<table>
		<tr>
			<th>Username</th>
			<th>External system</th>
			<th>Enabled at</th>
			<th>Disabled at</th>
			<th>Deleted at</th>
		</tr>
		{#if identity.accounts.length > 0}
			{#each identity.accounts as account}
				<tr>
					<td>{account.username}</td>
					<td>{account.systemId}</td>
					<td>{account.enabledAt}</td>
					<td>{account.disabledAt}</td>
					<td>{account.deletedAt}</td>
				</tr>
			{/each}
		{/if}
	</table>
	<h2>Memberships</h2>
	<button on:click={() => toggleDialog()}>Add membership</button>

	<dialog open={isOpen}>
		<p>Add membership</p>
		<label for="group">Group</label>
		<input name="group" bind:value={dialogData.groupId} />
		<label for="enabledAt">Start date</label>
		<input name="enabledAt" bind:value={dialogData.enabledAt} />
		<label for="disabledAt">End date</label>
		<input name="disabledAt" bind:value={dialogData.disabledAt} />
		<button on:click={() => toggleDialog()}>Cancel</button>
		<button on:click={() => createMembership()} type="submit">Create</button>
	</dialog>

	<table>
		<tr>
			<th>Group</th>
			<th>Enabled at</th>
			<th>Disabled at</th>
		</tr>
		{#if identity.memberships.length > 0}
			{#each identity.memberships as membership}
				<tr>
					<td>{membership.group.name}</td>
					<td>{membership.enabledAt}</td>
					<td>{membership.disabledAt}</td>
				</tr>
			{/each}
		{/if}
	</table>
{/if}
