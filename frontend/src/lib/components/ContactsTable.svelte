<script lang="ts">
	import { onMount } from 'svelte';
	import ContactDetailsModal from './ContactDetailsModal.svelte';
	import EditContactModal from './EditContactModal.svelte';
	import DeleteContactModal from './DeleteContactModal.svelte';

	type Contact = {
		id: number;
		name: string;
		email: string | null;
		phone: string | null;
		facebookId: string | null;
		notes: string | null;
		type: 'regular' | 'user_reference';
		createdAt: string;
		updatedAt: string;
		debtCount?: number;
		totalOwed?: number;
		totalOwing?: number;
	};

	let contacts: Contact[] = [];
	let filteredContacts: Contact[] = [];
	let selectedContact: Contact | null = null;
	let showDetailsModal = false;
	let showEditModal = false;
	let showDeleteDialog = false;
	let contactToDelete: Contact | null = null;

	// Filter and search state
	let searchQuery = '';
	let typeFilter = 'all';
	let sortBy = 'name';
	let sortOrder: 'asc' | 'desc' = 'asc';

	// Pagination
	let currentPage = 1;
	let itemsPerPage = 10;
	let totalPages = 1;

	onMount(() => {
		loadContacts();
	});

	function loadContacts() {
		// Mock data - replace with actual API call
		contacts = [
			{
				id: 1,
				name: 'Alice Johnson',
				email: 'alice@example.com',
				phone: '+63 917 123 4567',
				facebookId: 'alice.johnson',
				notes: 'Close friend from college',
				type: 'regular',
				createdAt: '2023-12-01T10:00:00Z',
				updatedAt: '2024-01-10T14:30:00Z',
				debtCount: 3,
				totalOwed: 5000,
				totalOwing: 1200
			},
			{
				id: 2,
				name: 'Bob Smith',
				email: 'bob.smith@example.com',
				phone: '+63 917 234 5678',
				facebookId: null,
				notes: 'Business partner',
				type: 'user_reference',
				createdAt: '2023-11-15T09:00:00Z',
				updatedAt: '2023-12-20T16:45:00Z',
				debtCount: 2,
				totalOwed: 0,
				totalOwing: 3500
			},
			{
				id: 3,
				name: 'Carol Davis',
				email: 'carol.davis@example.com',
				phone: null,
				facebookId: 'carol.davis.5',
				notes: null,
				type: 'regular',
				createdAt: '2023-10-20T11:30:00Z',
				updatedAt: '2023-10-20T11:30:00Z',
				debtCount: 1,
				totalOwed: 800,
				totalOwing: 0
			},
			{
				id: 4,
				name: 'David Wilson',
				email: null,
				phone: '+63 917 345 6789',
				facebookId: null,
				notes: 'Neighbor',
				type: 'regular',
				createdAt: '2024-01-01T08:00:00Z',
				updatedAt: '2024-01-05T10:15:00Z',
				debtCount: 0,
				totalOwed: 0,
				totalOwing: 0
			},
			{
				id: 5,
				name: 'Emma Brown',
				email: 'emma.brown@example.com',
				phone: '+63 917 456 7890',
				facebookId: 'emma.brown',
				notes: 'Family friend, very reliable',
				type: 'regular',
				createdAt: '2023-09-15T13:00:00Z',
				updatedAt: '2024-01-12T09:00:00Z',
				debtCount: 5,
				totalOwed: 12000,
				totalOwing: 0
			}
		];
		filterAndSortContacts();
	}

	function filterAndSortContacts() {
		let filtered = contacts;

		// Apply search filter
		if (searchQuery) {
			const query = searchQuery.toLowerCase();
			filtered = filtered.filter(contact => 
				contact.name.toLowerCase().includes(query) ||
				contact.email?.toLowerCase().includes(query) ||
				contact.phone?.includes(query) ||
				contact.facebookId?.toLowerCase().includes(query) ||
				contact.notes?.toLowerCase().includes(query)
			);
		}

		// Apply type filter
		if (typeFilter !== 'all') {
			filtered = filtered.filter(contact => contact.type === typeFilter);
		}

		// Apply sorting
		filtered.sort((a, b) => {
			let aValue: any = a[sortBy as keyof Contact];
			let bValue: any = b[sortBy as keyof Contact];

			if (sortBy === 'createdAt' || sortBy === 'updatedAt') {
				aValue = new Date(aValue).getTime();
				bValue = new Date(bValue).getTime();
			}

			if (typeof aValue === 'string' && typeof bValue === 'string') {
				aValue = aValue.toLowerCase();
				bValue = bValue.toLowerCase();
			}

			// Handle null values
			if (aValue === null) aValue = '';
			if (bValue === null) bValue = '';

			if (sortOrder === 'asc') {
				return aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
			} else {
				return aValue > bValue ? -1 : aValue < bValue ? 1 : 0;
			}
		});

		filteredContacts = filtered;
		totalPages = Math.ceil(filteredContacts.length / itemsPerPage);
		currentPage = Math.min(currentPage, totalPages || 1);
	}

	function handleSort(column: string) {
		if (sortBy === column) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortBy = column;
			sortOrder = 'asc';
		}
		filterAndSortContacts();
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = now.getTime() - date.getTime();
		const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

		if (diffDays === 0) {
			return 'Today';
		} else if (diffDays === 1) {
			return 'Yesterday';
		} else if (diffDays < 7) {
			return `${diffDays} days ago`;
		} else if (diffDays < 30) {
			const weeks = Math.floor(diffDays / 7);
			return `${weeks} week${weeks > 1 ? 's' : ''} ago`;
		} else if (diffDays < 365) {
			const months = Math.floor(diffDays / 30);
			return `${months} month${months > 1 ? 's' : ''} ago`;
		} else {
			return date.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			});
		}
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

	function viewContact(contact: Contact) {
		selectedContact = contact;
		showDetailsModal = true;
	}

	function editContact(contact: Contact) {
		selectedContact = contact;
		showEditModal = true;
	}

	function confirmDeleteContact(contact: Contact) {
		contactToDelete = contact;
		showDeleteDialog = true;
	}

	function deleteContact() {
		if (contactToDelete) {
			contacts = contacts.filter(c => c.id !== contactToDelete?.id);
			filterAndSortContacts();
			contactToDelete = null;
			showDeleteDialog = false;
		}
	}



	function handleContactUpdated(event: CustomEvent) {
		const updatedContact = event.detail;
		const index = contacts.findIndex(c => c.id === updatedContact.id);
		if (index !== -1) {
			contacts[index] = updatedContact;
			filterAndSortContacts();
		}
		showEditModal = false;
	}

	$: {
		// React to filter changes
		filterAndSortContacts();
	}

	$: paginatedContacts = filteredContacts.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);
</script>

<div class="space-y-6">
	<!-- Header with Search and Filters -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
		<div class="flex-1 max-w-md">
			<div class="relative">
				<svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search contacts..."
					class="input pl-10"
				/>
			</div>
		</div>

		<div class="flex items-center space-x-4">
			<select bind:value={typeFilter} class="input">
				<option value="all">All Types</option>
				<option value="regular">Regular Contact</option>
				<option value="user_reference">User Reference</option>
			</select>
		</div>
	</div>

	<!-- Desktop Table -->
	<div class="hidden lg:block card overflow-hidden">
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead class="bg-muted/50 border-b border-border">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('name')}>
							Name
							{#if sortBy === 'name'}
								<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
							{/if}
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Contact Info
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('type')}>
							Type
							{#if sortBy === 'type'}
								<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
							{/if}
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Debt Summary
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider cursor-pointer" on:click={() => handleSort('createdAt')}>
							Created
							{#if sortBy === 'createdAt'}
								<span class="ml-1">{sortOrder === 'asc' ? '↑' : '↓'}</span>
							{/if}
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Actions
						</th>
					</tr>
				</thead>
				<tbody class="bg-card divide-y divide-border">
					{#each paginatedContacts as contact (contact.id)}
						<tr class="hover:bg-muted/30 cursor-pointer transition-colors duration-200" on:click={() => viewContact(contact)}>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="flex items-center">
									<div class="w-10 h-10 bg-primary rounded-full flex items-center justify-center mr-3">
										<span class="text-primary-foreground text-sm font-medium">
											{contact.name.split(' ').map(n => n[0]).join('')}
										</span>
									</div>
									<div>
										<div class="text-sm font-medium text-foreground">{contact.name}</div>
										{#if contact.notes}
											<div class="text-sm text-muted-foreground truncate max-w-48">{contact.notes}</div>
										{/if}
									</div>
								</div>
							</td>
							<td class="px-6 py-4">
								<div class="space-y-1 text-sm">
									{#if contact.email}
										<div class="flex items-center space-x-2">
											<svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
											</svg>
											<span class="text-muted-foreground">{contact.email}</span>
										</div>
									{/if}
									{#if contact.phone}
										<div class="flex items-center space-x-2">
											<svg class="w-4 h-4 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
											</svg>
											<span class="text-muted-foreground">{contact.phone}</span>
										</div>
									{/if}
									{#if contact.facebookId}
										<div class="flex items-center space-x-2">
											<svg class="w-4 h-4 text-muted-foreground" fill="currentColor" viewBox="0 0 24 24">
												<path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
											</svg>
											<span class="text-muted-foreground">{contact.facebookId}</span>
										</div>
									{/if}
									{#if !contact.email && !contact.phone && !contact.facebookId}
										<span class="text-muted-foreground/60">No contact info</span>
									{/if}
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getTypeBadgeClass(contact.type)}">
									{contact.type === 'user_reference' ? 'User Reference' : 'Regular Contact'}
								</span>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								{#if contact.debtCount && contact.debtCount > 0}
									<div class="text-sm">
										<div class="font-medium text-foreground">{contact.debtCount} debt{contact.debtCount > 1 ? 's' : ''}</div>
										<div class="text-xs text-muted-foreground">
											{#if contact.totalOwed && contact.totalOwed > 0}
												<span class="text-success">+{formatCurrency(contact.totalOwed)}</span>
											{/if}
											{#if contact.totalOwing && contact.totalOwing > 0}
												<span class="text-destructive ml-2">-{formatCurrency(contact.totalOwing)}</span>
											{/if}
										</div>
									</div>
								{:else}
									<span class="text-sm text-muted-foreground">No debts</span>
								{/if}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm text-muted-foreground">
								{formatDate(contact.createdAt)}
							</td>
							<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
								<div class="flex items-center space-x-2" on:click|stopPropagation on:keydown|stopPropagation role="group">
									<button
										on:click={() => viewContact(contact)}
										class="text-primary hover:text-primary/80 p-1"
										title="View Details"
										aria-label="View details for {contact.name}"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
										</svg>
									</button>
									<button
										on:click={() => editContact(contact)}
										class="text-secondary hover:text-secondary/80 p-1"
										title="Edit"
										aria-label="Edit {contact.name}"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
										</svg>
									</button>
									<button
										on:click={() => confirmDeleteContact(contact)}
										class="text-destructive hover:text-destructive/80 p-1"
										title="Delete"
										aria-label="Delete {contact.name}"
									>
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
										</svg>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<!-- Mobile Card Layout -->
	<div class="lg:hidden space-y-4">
		{#each paginatedContacts as contact (contact.id)}
			<div class="card p-4" on:click={() => viewContact(contact)} on:keydown={(e) => e.key === 'Enter' && viewContact(contact)} role="button" tabindex="0">
				<div class="flex items-start justify-between mb-3">
					<div class="flex items-center space-x-3">
						<div class="w-10 h-10 bg-primary rounded-full flex items-center justify-center">
							<span class="text-primary-foreground text-sm font-medium">
								{contact.name.split(' ').map(n => n[0]).join('')}
							</span>
						</div>
						<div>
							<div class="font-medium text-foreground">{contact.name}</div>
							<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {getTypeBadgeClass(contact.type)}">
								{contact.type === 'user_reference' ? 'User Reference' : 'Regular Contact'}
							</span>
						</div>
					</div>
					<span class="text-xs text-muted-foreground">{formatDate(contact.createdAt)}</span>
				</div>
				
				<div class="space-y-2 text-sm mb-3">
					{#if contact.email}
						<div class="flex items-center space-x-2 text-muted-foreground">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
							</svg>
							<span>{contact.email}</span>
						</div>
					{/if}
					{#if contact.phone}
						<div class="flex items-center space-x-2 text-muted-foreground">
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
							</svg>
							<span>{contact.phone}</span>
						</div>
					{/if}
				</div>

				{#if contact.debtCount && contact.debtCount > 0}
					<div class="text-sm mb-3">
						<span class="font-medium">{contact.debtCount} debt{contact.debtCount > 1 ? 's' : ''}</span>
						{#if contact.totalOwed && contact.totalOwed > 0}
							<span class="text-success ml-2">+{formatCurrency(contact.totalOwed)}</span>
						{/if}
						{#if contact.totalOwing && contact.totalOwing > 0}
							<span class="text-destructive ml-2">-{formatCurrency(contact.totalOwing)}</span>
						{/if}
					</div>
				{/if}

				<div class="flex justify-end space-x-2" on:click|stopPropagation on:keydown|stopPropagation role="group">
					<button on:click={() => viewContact(contact)} class="btn-secondary text-xs px-3 py-1">View</button>
					<button on:click={() => editContact(contact)} class="btn-secondary text-xs px-3 py-1">Edit</button>
				</div>
			</div>
		{/each}
	</div>

	<!-- Pagination -->
	{#if totalPages > 1}
		<div class="flex items-center justify-between">
			<div class="text-sm text-muted-foreground">
				Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, filteredContacts.length)} of {filteredContacts.length} contacts
			</div>
			<div class="flex items-center space-x-2">
				<button
					on:click={() => currentPage = Math.max(1, currentPage - 1)}
					disabled={currentPage === 1}
					class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Previous
				</button>
				<span class="text-sm text-muted-foreground">
					Page {currentPage} of {totalPages}
				</span>
				<button
					on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
					disabled={currentPage === totalPages}
					class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Next
				</button>
			</div>
		</div>
	{/if}

	<!-- Empty State -->
	{#if filteredContacts.length === 0}
		<div class="text-center py-12">
			<svg class="mx-auto w-12 h-12 text-muted-foreground mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
			</svg>
			<h3 class="text-lg font-medium text-foreground mb-2">No contacts found</h3>
			<p class="text-muted-foreground mb-4">
				{searchQuery || typeFilter !== 'all' 
					? 'Try adjusting your filters or search query.'
					: 'Get started by adding your first contact.'}
			</p>
		</div>
	{/if}
</div>

<!-- Modals -->
{#if showDetailsModal && selectedContact}
	<ContactDetailsModal
		contact={selectedContact}
		on:close={() => { showDetailsModal = false; selectedContact = null; }}
		on:edit={() => { showDetailsModal = false; if (selectedContact) editContact(selectedContact); }}
		on:delete={() => { showDetailsModal = false; if (selectedContact) confirmDeleteContact(selectedContact); }}
	/>
{/if}

{#if showEditModal && selectedContact}
	<EditContactModal
		contact={selectedContact}
		on:close={() => { showEditModal = false; selectedContact = null; }}
		on:contact-updated={handleContactUpdated}
	/>
{/if}

{#if showDeleteDialog && contactToDelete}
	<DeleteContactModal
		contact={contactToDelete}
		on:confirm={deleteContact}
		on:close={() => { showDeleteDialog = false; contactToDelete = null; }}
	/>
{/if}
