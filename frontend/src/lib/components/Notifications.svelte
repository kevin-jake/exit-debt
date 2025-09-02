<script lang="ts">
	import { notificationsStore, type Notification } from '$lib/stores/notifications';
	import { onMount } from 'svelte';

	let notifications: Notification[] = [];

	onMount(() => {
		const unsubscribe = notificationsStore.subscribe((state) => {
			notifications = state.notifications;
		});

		return unsubscribe;
	});

	function removeNotification(id: string) {
		notificationsStore.remove(id);
	}

	function getIcon(type: Notification['type']) {
		switch (type) {
			case 'success':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
				</svg>`;
			case 'error':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>`;
			case 'warning':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
				</svg>`;
			case 'info':
				return `<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>`;
			default:
				return '';
		}
	}

	function getTypeClasses(type: Notification['type']) {
		switch (type) {
			case 'success':
				return 'bg-green-50 border-green-200 text-green-800';
			case 'error':
				return 'bg-red-50 border-red-200 text-red-800';
			case 'warning':
				return 'bg-yellow-50 border-yellow-200 text-yellow-800';
			case 'info':
				return 'bg-blue-50 border-blue-200 text-blue-800';
			default:
				return 'bg-gray-50 border-gray-200 text-gray-800';
		}
	}
</script>

<!-- Notifications Container -->
<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
	{#each notifications as notification (notification.id)}
		<div
			class="notification-enter"
			class:notification-exit={false}
			class:notification-exit-active={false}
		>
			<div
				class="flex items-start p-4 border rounded-lg shadow-lg {getTypeClasses(notification.type)}"
				role="alert"
			>
				<div class="flex-shrink-0 mr-3">
					{@html getIcon(notification.type)}
				</div>
				
				<div class="flex-1 min-w-0">
					<h4 class="text-sm font-medium mb-1">
						{notification.title}
					</h4>
					<p class="text-sm opacity-90">
						{notification.message}
					</p>
				</div>
				
				<button
					on:click={() => removeNotification(notification.id)}
					class="flex-shrink-0 ml-3 text-current opacity-60 hover:opacity-100 transition-opacity"
					aria-label="Close notification"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
		</div>
	{/each}
</div>

<style>
	.notification-enter {
		animation: slideIn 0.3s ease-out;
	}

	.notification-exit {
		animation: slideOut 0.3s ease-in;
	}

	@keyframes slideIn {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes slideOut {
		from {
			transform: translateX(0);
			opacity: 1;
		}
		to {
			transform: translateX(100%);
			opacity: 0;
		}
	}
</style>
