<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	export let debt: any;

	const dispatch = createEventDispatcher();

	// State management
	let confirmationText = '';
	let confirmationChecked = false;
	let isLoading = false;
	let paymentCount = 0;
	let remainingPayments = 0;

	// Validation
	$: isConfirmationValid = confirmationText === 'DELETE' && confirmationChecked;

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		loadDebtDetails();
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function loadDebtDetails() {
		// Mock data - replace with actual API call
		// In real app, fetch payment history count and remaining payments
		paymentCount = Math.floor(Math.random() * 10) + 1;
		if (debt.installmentPlan !== 'onetime') {
			remainingPayments = Math.max(0, debt.numberOfPayments - paymentCount);
		}
	}

	function formatCurrency(amount: number, currency: string = 'PHP'): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: currency
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'long',
			day: 'numeric',
			year: 'numeric'
		});
	}

	function getStatusBadgeClass(status: string): string {
		switch (status) {
			case 'active':
				return 'bg-primary/10 text-primary';
			case 'settled':
				return 'bg-success/10 text-success';
			case 'archived':
				return 'bg-muted/50 text-muted-foreground';
			case 'overdue':
				return 'bg-destructive/10 text-destructive';
			default:
				return 'bg-muted/50 text-muted-foreground';
		}
	}

	async function handleDelete() {
		if (!isConfirmationValid) return;

		isLoading = true;

		try {
			// TODO: Replace with actual API call
			await new Promise(resolve => setTimeout(resolve, 1500));
			
			dispatch('confirm');
		} catch (error) {
			console.error('Error deleting debt:', error);
			// Handle error - in real app, show error message
		} finally {
			isLoading = false;
		}
	}

	function handleClose() {
		if (confirmationText && !isLoading) {
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
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border bg-destructive/5">
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-3">
					<div class="w-10 h-10 bg-destructive/10 rounded-full flex items-center justify-center">
						<svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
						</svg>
					</div>
					<h2 class="text-xl font-semibold text-foreground">Delete Debt List</h2>
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
					<div class="flex items-start space-x-3">
						<svg class="w-5 h-5 text-destructive mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.865-.833-2.635 0L4.178 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
						</svg>
						<div>
							<p class="font-medium text-destructive">Are you sure you want to delete this debt list?</p>
							<p class="text-sm text-destructive/80 mt-1">This action cannot be undone and will permanently remove all associated data.</p>
						</div>
					</div>
				</div>

				<!-- Debt Information Summary -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-foreground">Debt Information</h3>
					<div class="bg-muted/50 rounded-lg p-4 space-y-3">
						<div class="flex items-center justify-between">
							<span class="text-sm text-muted-foreground">Contact</span>
							<div class="flex items-center space-x-2">
								<div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center">
									<span class="text-primary-foreground text-xs font-medium">
										{debt.contactName.split(' ').map((n: string) => n[0]).join('')}
									</span>
								</div>
								<span class="text-sm font-medium text-foreground">{debt.contactName}</span>
							</div>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-muted-foreground">Type</span>
							<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {debt.type === 'owed_to_me' ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive'}">
								{debt.type === 'owed_to_me' ? 'Owed to Me' : 'I Owe'}
							</span>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-muted-foreground">Total Amount</span>
							<span class="text-sm font-medium text-foreground">{formatCurrency(debt.totalAmount, debt.currency)}</span>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-muted-foreground">Status</span>
							<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getStatusBadgeClass(debt.status)}">
								{debt.status.charAt(0).toUpperCase() + debt.status.slice(1)}
							</span>
						</div>
						<div class="flex items-center justify-between">
							<span class="text-sm text-muted-foreground">Due Date</span>
							<span class="text-sm text-foreground">{formatDate(debt.dueDate)}</span>
						</div>
					</div>
				</div>

				<!-- Impact Assessment -->
				<div class="space-y-4">
					<h3 class="text-sm font-medium text-destructive">This will permanently delete:</h3>
					<div class="space-y-2">
						<div class="flex items-center space-x-3">
							<svg class="w-4 h-4 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
							</svg>
							<span class="text-sm text-foreground">Debt list record</span>
						</div>
						<div class="flex items-center space-x-3">
							<svg class="w-4 h-4 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
							</svg>
							<span class="text-sm text-foreground">All payment history ({paymentCount} payment{paymentCount !== 1 ? 's' : ''})</span>
						</div>
						<div class="flex items-center space-x-3">
							<svg class="w-4 h-4 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
							</svg>
							<span class="text-sm text-foreground">Payment schedule</span>
						</div>
						{#if debt.notes}
							<div class="flex items-center space-x-3">
								<svg class="w-4 h-4 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
								</svg>
								<span class="text-sm text-foreground">Any associated notes</span>
							</div>
						{/if}
					</div>
				</div>

				<!-- Related Data Check -->
				{#if paymentCount > 0 || remainingPayments > 0}
					<div class="bg-warning/10 border border-warning/20 rounded-lg p-4">
						<div class="flex items-start space-x-3">
							<svg class="w-5 h-5 text-warning mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
							<div class="space-y-1">
								{#if paymentCount > 0}
									<p class="text-sm font-medium text-warning">{paymentCount} payment record{paymentCount !== 1 ? 's' : ''} will be lost</p>
								{/if}
								{#if remainingPayments > 0}
									<p class="text-sm font-medium text-warning">{remainingPayments} remaining payment{remainingPayments !== 1 ? 's' : ''} will be deleted</p>
								{/if}
							</div>
						</div>
					</div>
				{/if}

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
	</div>
</div>
