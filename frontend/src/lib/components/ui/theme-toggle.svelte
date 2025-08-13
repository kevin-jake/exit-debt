<script lang="ts">
	import { Sun, Moon, Monitor } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { themeStore } from '$lib/stores/theme.svelte.js';

	let { variant = 'ghost', size = 'sm', ...props } = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Button builders={[builder]} {variant} {size} class="w-9 px-0" {...props}>
			{#if themeStore.isDark}
				<Moon class="h-[1.2rem] w-[1.2rem]" />
			{:else}
				<Sun class="h-[1.2rem] w-[1.2rem]" />
			{/if}
			<span class="sr-only">Toggle theme</span>
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Item 
			on:click={() => themeStore.setTheme('light')}
			class="cursor-pointer"
		>
			<Sun class="mr-2 h-4 w-4" />
			<span>Light</span>
		</DropdownMenu.Item>
		<DropdownMenu.Item 
			on:click={() => themeStore.setTheme('dark')}
			class="cursor-pointer"
		>
			<Moon class="mr-2 h-4 w-4" />
			<span>Dark</span>
		</DropdownMenu.Item>
		<DropdownMenu.Item 
			on:click={() => themeStore.setTheme('system')}
			class="cursor-pointer"
		>
			<Monitor class="mr-2 h-4 w-4" />
			<span>System</span>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
