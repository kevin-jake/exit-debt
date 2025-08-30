<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { goto } from '$app/navigation';

	export let contact: any;

	const dispatch = createEventDispatcher();

	let activeDebts = 0;
	let totalOwed = 0;
	let totalOwing = 0;

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		loadContactDebts();
		return () => {
			document.body.style.overflow = 'auto';
		};
	});

	function loadContactDebts() {
		// Mock data - replace with actual API call
		activeDebts = contact.debtCount || 0;
		totalOwed = contact.totalOwed || 0;
		totalOwing = contact.totalOwing || 0;
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'long',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-PH', {
			style: 'currency',
			currency: 'PHP'
		}).format(amount);
	}

	function getTypeBadgeClass(type: string): string {
		return type === 'user_reference' 
			? 'bg-success/10 text-success' 
			: 'bg-primary/10 text-primary';
	}

	function closeModal() {
		dispatch('close');
	}

	function editContact() {
		dispatch('edit');
	}

	function deleteContact() {
		dispatch('delete');
	}

	function createDebtWithContact() {
		// Navigate to create debt page with contact pre-selected
		goto(`/debts/new?contactId=${contact.id}`);
		closeModal();
	}

	function viewAllDebts() {
		// Navigate to debts page filtered by this contact
		goto(`/debts?contactId=${contact.id}`);
		closeModal();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeModal();
		}
	}
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={closeModal}>
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-3xl w-full max-h-[90vh] overflow-hidden flex flex-col"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border">
			<div class="flex items-center justify-between">
				<div class="flex items-center space-x-4">
					<div class="w-12 h-12 bg-primary rounded-full flex items-center justify-center">
						<span class="text-primary-foreground font-medium">
							{contact.name.split(' ').map((n: string) => n[0]).join('')}
						</span>
					</div>
					<div>
						<h2 class="text-xl font-semibold text-foreground">{contact.name}</h2>
						<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getTypeBadgeClass(contact.type)}">
							{contact.type === 'user_reference' ? 'User Reference' : 'Regular Contact'}
						</span>
					</div>
				</div>
				<button on:click={closeModal} class="text-muted-foreground hover:text-foreground">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto">
			<div class="p-6 space-y-6">
				<!-- Contact Information -->
				<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
					<!-- Personal Details -->
					<div class="space-y-4">
						<h3 class="text-lg font-medium text-foreground mb-4">Personal Details</h3>
						
						<div class="space-y-3">
							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Full Name</label>
								<div class="text-sm text-foreground font-medium">{contact.name}</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Email Address</label>
								<div class="text-sm text-foreground">
									{#if contact.email}
										<a href="mailto:{contact.email}" class="text-primary hover:text-primary/80 flex items-center space-x-2">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
											</svg>
											<span>{contact.email}</span>
										</a>
									{:else}
										<span class="text-muted-foreground">Not provided</span>
									{/if}
								</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Phone Number</label>
								<div class="text-sm text-foreground">
									{#if contact.phone}
										<a href="tel:{contact.phone}" class="text-primary hover:text-primary/80 flex items-center space-x-2">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
											</svg>
											<span>{contact.phone}</span>
										</a>
									{:else}
										<span class="text-muted-foreground">Not provided</span>
									{/if}
								</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Facebook ID</label>
								<div class="text-sm text-foreground">
									{#if contact.facebookId}
										<a href="https://facebook.com/{contact.facebookId}" target="_blank" rel="noopener noreferrer" class="text-primary hover:text-primary/80 flex items-center space-x-2">
											<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
												<path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
											</svg>
											<span>{contact.facebookId}</span>
										</a>
									{:else}
										<span class="text-muted-foreground">Not provided</span>
									{/if}
								</div>
							</div>

							{#if contact.notes}
								<div>
									<label class="block text-sm font-medium text-muted-foreground mb-1">Notes</label>
									<div class="text-sm text-foreground bg-muted/50 p-3 rounded-lg">{contact.notes}</div>
								</div>
							{/if}
						</div>
					</div>

					<!-- System Information -->
					<div class="space-y-4">
						<h3 class="text-lg font-medium text-foreground mb-4">System Information</h3>
						
						<div class="space-y-3">
							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Created Date</label>
								<div class="text-sm text-foreground">{formatDate(contact.createdAt)}</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Last Updated</label>
								<div class="text-sm text-foreground">{formatDate(contact.updatedAt)}</div>
							</div>

							<div>
								<label class="block text-sm font-medium text-muted-foreground mb-1">Contact Type</label>
								<div class="text-sm">
									<span class="inline-flex px-3 py-1 text-sm font-medium rounded-full {getTypeBadgeClass(contact.type)}">
										{contact.type === 'user_reference' ? 'User Reference' : 'Regular Contact'}
									</span>
								</div>
							</div>

							{#if contact.type === 'user_reference'}
								<div>
									<label class="block text-sm font-medium text-muted-foreground mb-1">User Reference</label>
									<div class="text-sm text-foreground">
										<span class="text-primary">Linked to system user</span>
									</div>
								</div>
							{/if}
						</div>
					</div>
				</div>

				<!-- Debt Summary -->
				<div class="border-t border-border pt-6">
					<h3 class="text-lg font-medium text-foreground mb-4">Debt Summary</h3>
					
					<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
						<div class="card p-4">
							<div class="flex items-center justify-between">
								<div>
									<p class="text-sm font-medium text-muted-foreground">Active Debts</p>
									<p class="text-2xl font-bold text-foreground">{activeDebts}</p>
								</div>
								<div class="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center">
									<svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
									</svg>
								</div>
							</div>
						</div>

						<div class="card p-4">
							<div class="flex items-center justify-between">
								<div>
									<p class="text-sm font-medium text-muted-foreground">They Owe You</p>
									<p class="text-2xl font-bold text-success">{formatCurrency(totalOwed)}</p>
								</div>
								<div class="w-10 h-10 bg-success/10 rounded-lg flex items-center justify-center">
									<svg class="w-5 h-5 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
									</svg>
								</div>
							</div>
						</div>

						<div class="card p-4">
							<div class="flex items-center justify-between">
								<div>
									<p class="text-sm font-medium text-muted-foreground">You Owe Them</p>
									<p class="text-2xl font-bold text-destructive">{formatCurrency(totalOwing)}</p>
								</div>
								<div class="w-10 h-10 bg-destructive/10 rounded-lg flex items-center justify-center">
									<svg class="w-5 h-5 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"></path>
									</svg>
								</div>
							</div>
						</div>
					</div>

					<!-- Quick Actions -->
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<button 
							on:click={createDebtWithContact}
							class="flex items-center justify-center space-x-2 p-4 bg-primary/10 hover:bg-primary/20 rounded-lg transition-colors duration-200"
						>
							<svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
							</svg>
							<span class="font-medium text-primary">Create New Debt</span>
						</button>

						{#if activeDebts > 0}
							<button 
								on:click={viewAllDebts}
								class="flex items-center justify-center space-x-2 p-4 bg-secondary/10 hover:bg-secondary/20 rounded-lg transition-colors duration-200"
							>
								<svg class="w-5 h-5 text-secondary-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
								</svg>
								<span class="font-medium text-secondary-foreground">View All Debts</span>
							</button>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Footer Actions -->
		<div class="px-6 py-4 border-t border-border flex items-center justify-between">
			<div class="flex items-center space-x-3">
				<button on:click={editContact} class="btn-secondary">
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
					</svg>
					Edit Contact
				</button>
				<button on:click={deleteContact} class="btn-danger">
					<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
					</svg>
					Delete Contact
				</button>
			</div>
			<button on:click={closeModal} class="btn-primary">
				Close
			</button>
		</div>
	</div>
</div>
