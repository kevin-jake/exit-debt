<script lang="ts">
	import DebtTable from '$lib/components/DebtTable.svelte';
</script>

<svelte:head>
	<title>Debts - DebtTracker</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-foreground">Debt Management</h1>
			<p class="text-muted-foreground mt-1">Track and manage all your debts and payment schedules</p>
		</div>
	</div>

	<!-- Debt Table -->
	<DebtTable />
</div>
