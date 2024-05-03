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
		try {
			const response = await fetch(apiUrl + 'identity/' + identityId);
			identity = await response.json();
		} catch (error) {
			console.log(error);
		}
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

<div id="container">
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

		<div id="membershipHeader">
			<h2>Memberships</h2>
			<button on:click={() => toggleDialog()}>Add membership</button>
		</div>

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
</div>

<style>
	#container {
		width: 60%;
		margin-left: auto;
		margin-right: auto;
	}

	#membershipHeader {
		display: flex;
		justify-content: space-between;
	}

	#membershipHeader button {
		height: 3.5em;
		border-radius: 8px;
		border: 1px solid #000;
		background-color: #2e83f2;
		color: #f1f1f1;
		margin-top: 1em;
	}

	table {
		width: 100%;
		text-align: left;
		border: 1px solid;
		border-collapse: collapse;
	}

	th,
	td {
		padding: 0.3em;
		border: 1px solid;
	}

	table tr:nth-child(even) {
		background-color: #dbdbdb;
	}
</style>
