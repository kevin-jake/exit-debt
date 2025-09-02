<script lang="ts">
	import { apiClient } from '$lib/api';
	import { authStore } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	let healthStatus = $state('');
	let loginTestResult = $state('');
	let authStoreStatus = $state('');
	let isLoading = $state(false);

	async function testHealthCheck() {
		try {
			const result = await apiClient.healthCheck();
			healthStatus = `✅ Health check successful: ${result.service} v${result.version}`;
		} catch (error) {
			healthStatus = `❌ Health check failed: ${error instanceof Error ? error.message : 'Unknown error'}`;
		}
	}

	async function testLogin() {
		isLoading = true;
		loginTestResult = '';
		
		try {
			const result = await apiClient.login({
				email: 'test@example.com',
				password: 'password123'
			});
			loginTestResult = `✅ Login successful! User: ${result.user.first_name} ${result.user.last_name}`;
			
			// Test auth store update
			authStore.setUser(result.user);
			authStoreStatus = `✅ Auth store updated with user: ${result.user.first_name} ${result.user.last_name}`;
		} catch (error) {
			loginTestResult = `❌ Login failed: ${error instanceof Error ? error.message : 'Unknown error'}`;
		} finally {
			isLoading = false;
		}
	}

	function testAuthStore() {
		authStoreStatus = `Current auth state: ${JSON.stringify($authStore, null, 2)}`;
	}

	onMount(() => {
		testHealthCheck();
		testAuthStore();
	});
</script>

<svelte:head>
	<title>API Test - DebtTracker</title>
</svelte:head>

<div class="max-w-2xl mx-auto p-6">
	<h1 class="text-3xl font-bold mb-6">API Connection Test</h1>
	
	<div class="space-y-6">
		<!-- Health Check -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold mb-4">Backend Health Check</h2>
			<div class="mb-4">
				<button 
					on:click={testHealthCheck}
					class="btn-primary"
				>
					Test Health Check
				</button>
			</div>
			{#if healthStatus}
				<div class="p-4 rounded-lg bg-muted">
					<p class="text-sm">{healthStatus}</p>
				</div>
			{/if}
		</div>

		<!-- Login Test -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold mb-4">Login API Test</h2>
			<div class="mb-4">
				<button 
					on:click={testLogin}
					disabled={isLoading}
					class="btn-primary disabled:opacity-50"
				>
					{#if isLoading}
						Testing...
					{:else}
						Test Login API
					{/if}
				</button>
			</div>
			{#if loginTestResult}
				<div class="p-4 rounded-lg bg-muted">
					<p class="text-sm">{loginTestResult}</p>
				</div>
			{/if}
		</div>

		<!-- Auth Store Test -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold mb-4">Auth Store Test</h2>
			<div class="mb-4">
				<button 
					on:click={testAuthStore}
					class="btn-primary"
				>
					Check Auth Store
				</button>
			</div>
			{#if authStoreStatus}
				<div class="p-4 rounded-lg bg-muted">
					<pre class="text-sm whitespace-pre-wrap">{authStoreStatus}</pre>
				</div>
			{/if}
		</div>

		<!-- API Configuration -->
		<div class="card p-6">
			<h2 class="text-xl font-semibold mb-4">API Configuration</h2>
			<div class="space-y-2 text-sm">
				<p><strong>API Base URL:</strong> {import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'}</p>
				<p><strong>Environment:</strong> {import.meta.env.MODE}</p>
			</div>
		</div>
	</div>
</div>
