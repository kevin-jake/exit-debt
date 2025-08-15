<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import CreateContactModal from './CreateContactModal.svelte';

	export let selectedContact: any = null;
	export let contacts: any[] = [];

	const dispatch = createEventDispatcher();

	let isOpen = false;
	let searchQuery = '';
	let showCreateModal = false;

	$: filteredContacts = contacts.filter(contact =>
		contact.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
		contact.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
		contact.phone?.includes(searchQuery)
	);

	function selectContact(contact: any) {
		selectedContact = contact;
		dispatch('select', contact);
		isOpen = false;
		searchQuery = '';
	}

	function openCreateModal() {
		showCreateModal = true;
		isOpen = false;
	}

	function handleContactCreated(event: CustomEvent) {
		const newContact = event.detail;
		contacts = [...contacts, newContact];
		selectContact(newContact);
		showCreateModal = false;
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.contact-selector')) {
			isOpen = false;
		}
	}
</script>

<svelte:window on:click={handleClickOutside} />

<div class="contact-selector relative">
	<label class="label">Contact *</label>
	
	<!-- Selected Contact Display / Search Input -->
	<div class="relative">
		{#if selectedContact}
			<div class="input flex items-center justify-between cursor-pointer" on:click={() => isOpen = !isOpen}>
				<div class="flex items-center space-x-3">
					<div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
						<span class="text-primary-foreground text-xs font-medium">
							{selectedContact.name.split(' ').map((n: string) => n[0]).join('')}
						</span>
					</div>
					<div>
						<div class="font-medium text-foreground">{selectedContact.name}</div>
						{#if selectedContact.email || selectedContact.phone}
							<div class="text-sm text-muted-foreground">
								{selectedContact.email || selectedContact.phone}
							</div>
						{/if}
					</div>
				</div>
				<svg class="w-5 h-5 text-muted-foreground transform transition-transform {isOpen ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
				</svg>
			</div>
		{:else}
			<div class="relative">
				<input
					type="text"
					bind:value={searchQuery}
					on:focus={() => isOpen = true}
					on:input={() => isOpen = true}
					class="input pr-10"
					placeholder="Search contacts or select..."
				/>
				<svg class="absolute right-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
			</div>
		{/if}
	</div>

	<!-- Dropdown -->
	{#if isOpen}
		<div class="absolute z-10 w-full mt-1 bg-card border border-border rounded-lg shadow-medium max-h-64 overflow-y-auto">
			<!-- Search Input (when contact is selected) -->
			{#if selectedContact}
				<div class="p-3 border-b border-border">
					<div class="relative">
						<input
							type="text"
							bind:value={searchQuery}
							class="w-full px-3 py-2 border border-input rounded-lg text-sm focus:border-ring focus:outline-none focus:ring-1 focus:ring-ring bg-background text-foreground placeholder:text-muted-foreground"
							placeholder="Search contacts..."
							autofocus
						/>
						<svg class="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
						</svg>
					</div>
				</div>
			{/if}

			<!-- Contact List -->
			<div class="max-h-48 overflow-y-auto">
				{#if filteredContacts.length > 0}
					{#each filteredContacts as contact (contact.id)}
						<button
							type="button"
							class="w-full px-4 py-3 text-left hover:bg-muted/50 focus:bg-muted/50 focus:outline-none border-b border-border last:border-b-0"
							on:click={() => selectContact(contact)}
						>
							<div class="flex items-center space-x-3">
								<div class="w-8 h-8 bg-primary rounded-full flex items-center justify-center">
									<span class="text-primary-foreground text-xs font-medium">
										{contact.name.split(' ').map((n: string) => n[0]).join('')}
									</span>
								</div>
								<div class="flex-1 min-w-0">
									<div class="font-medium text-foreground">{contact.name}</div>
									{#if contact.email || contact.phone}
										<div class="text-sm text-muted-foreground truncate">
											{contact.email || contact.phone}
										</div>
									{/if}
								</div>
							</div>
						</button>
					{/each}
				{:else}
					<div class="px-4 py-6 text-center text-muted-foreground">
						<svg class="mx-auto w-8 h-8 text-muted-foreground mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
						</svg>
						<p class="text-sm font-medium mb-1">No contacts found</p>
						<p class="text-xs text-muted-foreground/60">
							{searchQuery ? 'Try adjusting your search' : 'No contacts available'}
						</p>
					</div>
				{/if}
			</div>

			<!-- Create Contact Button -->
			<div class="p-3 border-t border-border">
				<button
					type="button"
					on:click={openCreateModal}
					class="w-full flex items-center justify-center space-x-2 px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 transition-colors duration-200"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
					</svg>
					<span class="text-sm font-medium">Create New Contact</span>
				</button>
			</div>
		</div>
	{/if}

	<!-- Clear Selection Button -->
	{#if selectedContact}
		<button
			type="button"
			on:click={() => { selectedContact = null; dispatch('select', null); }}
			class="absolute right-2 top-9 text-muted-foreground hover:text-foreground z-20"
			title="Clear selection"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
			</svg>
		</button>
	{/if}
</div>

<!-- Create Contact Modal -->
{#if showCreateModal}
	<CreateContactModal
		on:contact-created={handleContactCreated}
		on:close={() => showCreateModal = false}
	/>
{/if}
