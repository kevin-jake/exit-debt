<script lang="ts">
	import LoginForm from '$lib/components/auth/login-form.svelte';
	import { APP_CONFIG } from '$lib/config/api.js';
</script>

<svelte:head>
	<title>Login - {APP_CONFIG.NAME}</title>
	<meta name="description" content="Sign in to your Exit-Debt account" />
</svelte:head>

<div class="flex w-full min-h-screen items-center justify-center bg-background px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8">
		<!-- Logo/Brand -->
		<div class="text-center">
			<div class="mx-auto h-12 w-12 rounded-lg bg-primary flex items-center justify-center">
				<span class="text-primary-foreground font-bold text-lg">ED</span>
			</div>
			<h1 class="mt-4 text-3xl font-bold text-foreground">{APP_CONFIG.NAME}</h1>
			<p class="mt-2 text-sm text-muted-foreground">Take control of your debts</p>
		</div>

		<!-- Login Form -->
		<LoginForm />
	</div>
</div>
