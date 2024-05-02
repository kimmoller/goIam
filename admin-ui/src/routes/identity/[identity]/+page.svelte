<script lang="ts">
	import { onMount } from 'svelte';

	let identity: ExtendedIdentity;

	async function getIdentities() {
		const identityId = window.location.pathname.slice(-1);
		const response = await fetch('http://172.19.0.8:8083/identity/' + identityId);
		identity = await response.json();
		console.log(identity);
	}

	onMount(async () => {
		await getIdentities();
	});
</script>

<h1>IAM Admin UI</h1>

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
		{#each identity.accounts as account}
			<tr>
				<td>{account.username}</td>
				<td>{account.systemId}</td>
                <td>{account.enabledAt}</td>
                <td>{account.disabledAt}</td>
                <td>{account.deletedAt}</td>
			</tr>
		{/each}
	</table>
    <h2>Memberships</h2>
    <table>
		<tr>
			<th>Group</th>
			<th>Enabled at</th>
			<th>Disabled at</th>
		</tr>
		{#each identity.memberships as membership}
			<tr>
				<td>{membership.group.name}</td>
                <td>{membership.enabledAt}</td>
                <td>{membership.disabledAt}</td>
			</tr>
		{/each}
	</table>
{/if}
