<script lang="ts">
	import CreateDebtForm from '$lib/components/CreateDebtForm.svelte';
</script>

<svelte:head>
	<title>Add New Debt - DebtTracker</title>
</svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-foreground">Add New Debt</h1>
			<p class="text-muted-foreground mt-1">Create a new debt record with payment schedules and contact information</p>
		</div>
		<a href="/debts" class="btn-secondary">
			<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
			</svg>
			Back to Debts
		</a>
	</div>

	<!-- Form -->
	<CreateDebtForm />
</div>
