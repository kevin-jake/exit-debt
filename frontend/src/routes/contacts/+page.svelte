<script lang="ts">
	import ContactsTable from '$lib/components/ContactsTable.svelte';
	import CreateContactModal from '$lib/components/CreateContactModal.svelte';

	let showCreateModal = false;

	function handleContactCreated(event: CustomEvent) {
		// In a real app, this would refresh the contacts table
		showCreateModal = false;
	}
</script>

<svelte:head>
	<title>Contacts - DebtTracker</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-foreground">Contact Management</h1>
			<p class="text-muted-foreground mt-1">Organize and manage all your contacts and their information</p>
		</div>
		<button 
			on:click={() => showCreateModal = true}
			class="btn-primary"
		>
			<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
			</svg>
			Create Contact
		</button>
	</div>

	<!-- Contacts Table -->
	<ContactsTable />
</div>

<!-- Create Contact Modal -->
{#if showCreateModal}
	<CreateContactModal
		on:contact-created={handleContactCreated}
		on:close={() => showCreateModal = false}
	/>
{/if}
