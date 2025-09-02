<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { contactsStore } from '$lib/stores/contacts';
	import { notificationsStore } from '$lib/stores/notifications';
	import { handleApiError } from '$lib/utils/error-handling';
	import { validateForm, contactValidationRules } from '$lib/utils/validation';
	import type { Contact, UpdateContactRequest } from '$lib/api';

	export let contact: Contact;

	const dispatch = createEventDispatcher();

	let formData = {
		name: contact.name,
		email: contact.email || '',
		phone: contact.phone || '',
		notes: contact.notes || ''
	};

	let originalData = { ...formData };
	let errors: { [key: string]: string } = {};
	let isLoading = false;
	let hasChanges = false;

	onMount(() => {
		// Prevent body scroll when modal is open
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = 'auto';
		};
	});



	function validateFormData(): boolean {
		const validation = validateForm(formData, contactValidationRules);
		errors = validation.errors;
		return validation.isValid;
	}

	async function handleSubmit() {
		if (!validateFormData()) return;

		isLoading = true;

		try {
			// Prepare the update data
			const updateData: UpdateContactRequest = {
				name: formData.name.trim(),
				email: formData.email?.trim() || undefined,
				phone: formData.phone?.trim() || undefined,
				notes: formData.notes?.trim() || undefined,
			};

			// Update contact using the API
			const updatedContact = await contactsStore.updateContact(contact.id, updateData);
			
			// Show success notification
			notificationsStore.success(
				'Contact Updated',
				`Successfully updated ${updatedContact.name}`
			);

			// Dispatch the updated contact to parent component
			dispatch('contact-updated', updatedContact);
			
		} catch (error) {
			console.error('Error updating contact:', error);
			const errorMessage = handleApiError(error, 'EditContactModal');
			errors.general = errorMessage;
			
			// Show error notification
			notificationsStore.error(
				'Contact Update Failed',
				errorMessage
			);
		} finally {
			isLoading = false;
		}
	}

	function handleClose() {
		if (hasChanges) {
			const confirmed = confirm('You have unsaved changes. Are you sure you want to close?');
			if (!confirmed) return;
		}
		dispatch('close');
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}

	// Watch for changes
	$: hasChanges = JSON.stringify(formData) !== JSON.stringify(originalData);
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- Modal Backdrop -->
<div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" on:click={handleClose}>
	<!-- Modal Content -->
	<div 
		class="bg-card rounded-xl shadow-medium max-w-md w-full max-h-[90vh] overflow-hidden"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="px-6 py-4 border-b border-border flex items-center justify-between">
			<h2 class="text-xl font-semibold text-foreground">Edit Contact</h2>
			<button on:click={handleClose} class="text-muted-foreground hover:text-foreground">
				<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>

		<!-- Form -->
		<form on:submit|preventDefault={handleSubmit} class="p-6 space-y-4">
			{#if errors.general}
				<div class="p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
					<p class="text-sm text-destructive">{errors.general}</p>
				</div>
			{/if}

			<!-- Name Field -->
			<div>
				<label for="contact-name" class="label">Name *</label>
				<input
					id="contact-name"
					type="text"
					bind:value={formData.name}
					class="input {errors.name ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.name !== originalData.name ? 'border-warning' : ''}"
					placeholder="Enter contact name"
					disabled={isLoading}
					required
				/>
				{#if errors.name}
					<p class="mt-1 text-sm text-destructive">{errors.name}</p>
				{/if}
			</div>

			<!-- Email Field -->
			<div>
				<label for="contact-email" class="label">Email</label>
				<input
					id="contact-email"
					type="email"
					bind:value={formData.email}
					class="input {errors.email ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.email !== originalData.email ? 'border-warning' : ''}"
					placeholder="Enter email address"
					disabled={isLoading}
				/>
				{#if errors.email}
					<p class="mt-1 text-sm text-destructive">{errors.email}</p>
				{/if}
			</div>

			<!-- Phone Field -->
			<div>
				<label for="contact-phone" class="label">Phone</label>
				<input
					id="contact-phone"
					type="tel"
					bind:value={formData.phone}
					class="input {errors.phone ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''} {formData.phone !== originalData.phone ? 'border-warning' : ''}"
					placeholder="Enter phone number"
					disabled={isLoading}
				/>
				{#if errors.phone}
					<p class="mt-1 text-sm text-destructive">{errors.phone}</p>
				{/if}
			</div>

			<!-- Notes Field -->
			<div>
				<label for="contact-notes" class="label">Notes</label>
				<textarea
					id="contact-notes"
					bind:value={formData.notes}
					rows="3"
					class="input resize-none {formData.notes !== originalData.notes ? 'border-warning' : ''}"
					placeholder="Additional notes about this contact..."
					disabled={isLoading}
				></textarea>
				<div class="mt-1 text-xs text-muted-foreground">
					{formData.notes.length}/500 characters
				</div>
			</div>

			<!-- Contact Type Info -->
			<!-- Note: Contact type information removed as it's not available in the current API -->

			<!-- Action Buttons -->
			<div class="flex items-center justify-end space-x-3 pt-4">
				<button
					type="button"
					on:click={handleClose}
					class="btn-secondary"
					disabled={isLoading}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
					disabled={isLoading || !formData.name.trim() || !hasChanges}
				>
					{#if isLoading}
						<svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-primary-foreground" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Saving...
					{:else}
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						Save Changes
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
