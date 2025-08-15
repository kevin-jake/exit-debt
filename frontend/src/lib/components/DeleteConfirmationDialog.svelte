<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	export let debtName: string;

	const dispatch = createEventDispatcher();

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function confirm() {
		dispatch('confirm');
	}

	function cancel() {
		dispatch('cancel');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			cancel();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={cancel}>
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-md w-full"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border">
			<div class="flex items-center">
				<div class="w-10 h-10 bg-destructive/10 rounded-full flex items-center justify-center mr-4">
					<svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
					</svg>
				</div>
				<div>
					<h3 class="text-lg font-medium text-foreground">Delete Debt</h3>
					<p class="text-sm text-muted-foreground">This action cannot be undone</p>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="px-6 py-4">
			<p class="text-foreground">
				Are you sure you want to delete the debt record for <span class="font-medium">{debtName}</span>? 
				This will permanently remove all associated payment history and cannot be recovered.
			</p>
		</div>

		<!-- Footer Actions -->
		<div class="px-6 py-4 border-t border-border flex items-center justify-end space-x-3">
			<button on:click={cancel} class="btn-secondary">
				Cancel
			</button>
			<button on:click={confirm} class="btn-danger">
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
				</svg>
				Delete Debt
			</button>
		</div>
	</div>
</div>
