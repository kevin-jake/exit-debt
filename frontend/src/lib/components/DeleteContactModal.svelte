<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { goto } from '$app/navigation';

	export let contact: any;

	const dispatch = createEventDispatcher();

	// State management
	let confirmationText = '';
	let confirmationChecked = false;
	let isLoading = false;
	let relatedDebts: any[] = [];
	let hasRelatedDebts = false;
	let checkingDebts = true;

	// Validation
	$: isConfirmationValid = confirmationText === 'DELETE' && confirmationChecked;

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		checkRelatedDebts();
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	async function checkRelatedDebts() {
		checkingDebts = true;
		try {
			// TODO: Replace with actual API call
			// GET /api/v1/contacts/:id/debts
			await new Promise(resolve => setTimeout(resolve, 1000));

			// Mock data - simulate checking for related debts
			const mockDebts = Math.random() > 0.5 ? [
				{
					id: 1,
					type: 'owed_to_me',
					totalAmount: 2500.00,
					remainingBalance: 1800.00,
					status: 'active',
					description: 'Personal loan',
					currency: 'PHP'
				},
				{
					id: 2,
					type: 'i_owe',
					totalAmount: 1200.00,
					remainingBalance: 800.00,
					status: 'overdue',
					description: 'Car repair',
					currency: 'PHP'
				}
			] : [];

			relatedDebts = mockDebts;
			hasRelatedDebts = relatedDebts.length > 0;
		} catch (error) {
			console.error('Error checking related debts:', error);
			// In real app, show error message
		} finally {
			checkingDebts = false;
		}
	}

	function formatCurrency(amount: number, currency: string = 'PHP'): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: currency
		}).format(amount);
	}

	function getStatusBadgeClass(status: string): string {
		switch (status) {
			case 'active':
				return 'bg-primary/10 text-primary';
			case 'settled':
				return 'bg-success/10 text-success';
			case 'overdue':
				return 'bg-destructive/10 text-destructive';
			default:
				return 'bg-muted/50 text-muted-foreground';
		}
	}

	function getDebtTypeBadgeClass(type: string): string {
		return type === 'owed_to_me' 
			? 'bg-success/10 text-success' 
			: 'bg-destructive/10 text-destructive';
	}

	async function handleDelete() {
		if (!isConfirmationValid || hasRelatedDebts) return;

		isLoading = true;

		try {
			// TODO: Replace with actual API call
			// DELETE /api/v1/contacts/:id
			await new Promise(resolve => setTimeout(resolve, 1500));
			
			dispatch('confirm');
		} catch (error) {
			console.error('Error deleting contact:', error);
			// Handle error - in real app, show error message
		} finally {
			isLoading = false;
		}
	}

	function handleViewDebts() {
		dispatch('close');
		goto(`/debts?contactId=${contact.id}`);
	}

	function handleViewDebt(debtId: number) {
		dispatch('close');
		// In real app, this would open the debt details modal or navigate to debt page
		goto(`/debts?debtId=${debtId}`);
	}

	function handleClose() {
		if (confirmationText && !isLoading && !hasRelatedDebts) {
			const confirmed = confirm('Are you sure you want to cancel? Your confirmation will be lost.');
			if (!confirmed) return;
		}
		dispatch('close');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && !isLoading) {
			handleClose();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" role="dialog" aria-modal="true">
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-md w-full max-h-[90vh] overflow-hidden flex flex-col"
		on:click|stopPropagation
		on:keydown|stopPropagation
		role="document"
	>
		{#if checkingDebts}
			<!-- Loading State -->
			<div class="p-8 text-center">
				<svg class="animate-spin h-8 w-8 text-primary mx-auto mb-4" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				<p class="text-muted-foreground">Checking for related debts...</p>
			</div>
		{:else if hasRelatedDebts}
			<!-- Cannot Delete Scenario -->
			<!-- Header -->
			<div class="px-6 py-4 border-b border-border bg-warning/5">
				<div class="flex items-center justify-between">
					<div class="flex items-center space-x-3">
						<div class="w-10 h-10 bg-warning/10 rounded-full flex items-center justify-center">
							<svg class="w-6 h-6 text-warning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
							</svg>
						</div>
						<div>
							<h2 class="text-xl font-semibold text-foreground">Cannot Delete Contact</h2>
							<p class="text-sm text-muted-foreground">This contact has active debt lists</p>
						</div>
					</div>
					<button 
						on:click={handleClose} 
						class="text-muted-foreground hover:text-foreground"
						aria-label="Close modal"
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
				</div>
			</div>

			<!-- Content -->
			<div class="flex-1 overflow-y-auto">
				<div class="p-6 space-y-6">
					<!-- Explanation -->
					<div class="bg-warning/10 border border-warning/20 rounded-lg p-4">
						<p class="text-sm text-warning">
							To delete this contact, you must first settle or delete all associated debt lists.
						</p>
					</div>

					<!-- Contact Info -->
					<div class="flex items-center space-x-3 pb-4 border-b border-border">
						<div class="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
							<span class="text-primary-foreground font-medium">
								{contact.name.split(' ').map((n: string) => n[0]).join('')}
							</span>
						</div>
						<div>
							<div class="font-medium text-foreground">{contact.name}</div>
							{#if contact.email}
								<div class="text-sm text-muted-foreground">{contact.email}</div>
							{/if}
						</div>
					</div>

					<!-- Related Debts Section -->
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<h3 class="text-sm font-medium text-foreground">Related Debt Lists</h3>
							<span class="text-sm text-muted-foreground">Total: {relatedDebts.length} debt{relatedDebts.length !== 1 ? 's' : ''}</span>
						</div>
						
						<div class="space-y-3">
							{#each relatedDebts as debt (debt.id)}
								<button
									on:click={() => handleViewDebt(debt.id)}
									class="w-full text-left p-4 bg-muted/50 hover:bg-muted/70 rounded-lg transition-colors duration-200"
								>
									<div class="flex items-center justify-between mb-2">
										<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getDebtTypeBadgeClass(debt.type)}">
											{debt.type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
										</span>
										<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusBadgeClass(debt.status)}">
											{debt.status.charAt(0).toUpperCase() + debt.status.slice(1)}
										</span>
									</div>
									<div class="space-y-1">
										{#if debt.description}
											<p class="text-sm font-medium text-foreground">{debt.description}</p>
										{/if}
										<div class="flex items-center justify-between text-sm">
											<span class="text-muted-foreground">Total: {formatCurrency(debt.totalAmount, debt.currency)}</span>
											<span class="text-muted-foreground">Remaining: {formatCurrency(debt.remainingBalance, debt.currency)}</span>
										</div>
									</div>
									<div class="mt-2 text-xs text-primary">
										Click to view debt details →
									</div>
								</button>
							{/each}
						</div>
					</div>
				</div>
			</div>

			<!-- Footer Actions -->
			<div class="px-6 py-4 border-t border-border flex items-center justify-between">
				<button
					on:click={handleViewDebts}
					class="btn-secondary"
				>
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
					</svg>
					View All Debts
				</button>
				<button
					on:click={handleClose}
					class="btn-primary"
				>
					Close
				</button>
			</div>
		{:else}
			<!-- Can Delete Scenario -->
			<!-- Header -->
			<div class="px-6 py-4 border-b border-border bg-destructive/5">
				<div class="flex items-center justify-between">
					<div class="flex items-center space-x-3">
						<div class="w-10 h-10 bg-destructive/10 rounded-full flex items-center justify-center">
							<svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
							</svg>
						</div>
						<h2 class="text-xl font-semibold text-foreground">Delete Contact</h2>
					</div>
					<button 
						on:click={handleClose} 
						class="text-muted-foreground hover:text-foreground"
						aria-label="Close modal"
						disabled={isLoading}
					>
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
				</div>
			</div>

			<!-- Content -->
			<div class="flex-1 overflow-y-auto">
				<div class="p-6 space-y-6">
					<!-- Warning Message -->
					<div class="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
						<p class="text-sm font-medium text-destructive">Are you sure you want to delete this contact?</p>
						<p class="text-sm text-destructive/80 mt-1">This action cannot be undone.</p>
					</div>

					<!-- Contact Information -->
					<div class="space-y-4">
						<h3 class="text-sm font-medium text-foreground">Contact Information</h3>
						<div class="bg-muted/50 rounded-lg p-4">
							<div class="flex items-center space-x-3 mb-4">
								<div class="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
									<span class="text-primary-foreground font-medium">
										{contact.name.split(' ').map((n: string) => n[0]).join('')}
									</span>
								</div>
								<div>
									<div class="font-medium text-foreground">{contact.name}</div>
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {contact.type === 'user_reference' ? 'bg-success/10 text-success' : 'bg-primary/10 text-primary'}">
										{contact.type === 'user_reference' ? 'User Reference' : 'Regular Contact'}
									</span>
								</div>
							</div>
							
							<div class="space-y-2">
								{#if contact.email}
									<div class="flex items-center space-x-2 text-sm">
										<svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
										</svg>
										<span class="text-muted-foreground">{contact.email}</span>
									</div>
								{/if}
								{#if contact.phone}
									<div class="flex items-center space-x-2 text-sm">
										<svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
										</svg>
										<span class="text-muted-foreground">{contact.phone}</span>
									</div>
								{/if}
								{#if !contact.email && !contact.phone}
									<p class="text-sm text-muted-foreground/60">No additional contact information</p>
								{/if}
							</div>
						</div>
					</div>

					<!-- Confirmation Requirements -->
					<div class="space-y-4">
						<div>
							<label for="confirmation-text" class="label">Type 'DELETE' to confirm *</label>
							<input
								id="confirmation-text"
								type="text"
								bind:value={confirmationText}
								class="input {confirmationText && confirmationText !== 'DELETE' ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}"
								placeholder="Type DELETE"
								disabled={isLoading}
								autocomplete="off"
							/>
							{#if confirmationText && confirmationText !== 'DELETE'}
								<p class="mt-1 text-sm text-destructive">Please type 'DELETE' exactly</p>
							{/if}
						</div>

						<label class="flex items-start space-x-3 cursor-pointer">
							<input
								type="checkbox"
								bind:checked={confirmationChecked}
								class="w-4 h-4 text-destructive focus:ring-destructive border-input rounded mt-0.5"
								disabled={isLoading}
							/>
							<span class="text-sm text-foreground">I understand this action cannot be undone</span>
						</label>
					</div>
				</div>
			</div>

			<!-- Footer Actions -->
			<div class="px-6 py-4 border-t border-border flex items-center justify-end space-x-3">
				<button
					type="button"
					on:click={handleClose}
					class="btn-secondary"
					disabled={isLoading}
				>
					Cancel
				</button>
				<button
					on:click={handleDelete}
					class="btn-danger disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={!isConfirmationValid || isLoading}
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-danger-foreground" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Deleting...
					{:else}
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
						Delete Permanently
					{/if}
				</button>
			</div>
		{/if}
	</div>
</div>
