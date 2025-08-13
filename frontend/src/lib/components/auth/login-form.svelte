<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { Eye, EyeOff, Loader2 } from 'lucide-svelte';

	let email = $state('');
	let password = $state('');
	let showPassword = $state(false);
	let error = $state('');

	async function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';

		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}

		const result = await authStore.login({ email, password });
		
		if (!result.success) {
			error = result.error || 'Login failed';
		}
	}

	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}
</script>

<Card class="w-full max-w-md">
	<CardHeader class="space-y-1">
		<CardTitle class="text-2xl font-bold">Welcome back</CardTitle>
		<CardDescription>
			Enter your credentials to access your account
		</CardDescription>
	</CardHeader>
	<CardContent>
		<form on:submit={handleSubmit} class="space-y-4">
			{#if error}
				<div class="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
					{error}
				</div>
			{/if}

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input
					id="email"
					type="email"
					placeholder="Enter your email"
					bind:value={email}
					disabled={authStore.isLoading}
					required
				/>
			</div>

			<div class="space-y-2">
				<Label for="password">Password</Label>
				<div class="relative">
					<Input
						id="password"
						type={showPassword ? 'text' : 'password'}
						placeholder="Enter your password"
						bind:value={password}
						disabled={authStore.isLoading}
						required
					/>
					<Button
						type="button"
						variant="ghost"
						size="sm"
						class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
						on:click={togglePasswordVisibility}
						disabled={authStore.isLoading}
					>
						{#if showPassword}
							<EyeOff class="h-4 w-4" />
							<span class="sr-only">Hide password</span>
						{:else}
							<Eye class="h-4 w-4" />
							<span class="sr-only">Show password</span>
						{/if}
					</Button>
				</div>
			</div>

			<Button 
				type="submit" 
				class="w-full" 
				disabled={authStore.isLoading}
			>
				{#if authStore.isLoading}
					<Loader2 class="mr-2 h-4 w-4 animate-spin" />
				{/if}
				Sign In
			</Button>
		</form>

		<div class="mt-4 text-center text-sm">
			<a href="/auth/forgot-password" class="text-primary hover:underline">
				Forgot your password?
			</a>
		</div>

		<div class="mt-6 text-center text-sm text-muted-foreground">
			Don't have an account?{' '}
			<a href="/auth/register" class="text-primary hover:underline">
				Sign up
			</a>
		</div>
	</CardContent>
</Card>
