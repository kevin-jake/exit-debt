import { browser } from '$app/environment';

export type Theme = 'light' | 'dark' | 'system';

class ThemeStore {
	private currentTheme = $state<Theme>('dark'); // Default to dark as per requirements
	private systemPreference = $state<'light' | 'dark'>('dark');

	constructor() {
		if (browser) {
			this.initializeTheme();
			this.watchSystemPreference();
		}
	}

	get theme(): Theme {
		return this.currentTheme;
	}

	get resolvedTheme(): 'light' | 'dark' {
		if (this.currentTheme === 'system') {
			return this.systemPreference;
		}
		return this.currentTheme;
	}

	get isDark(): boolean {
		return this.resolvedTheme === 'dark';
	}

	setTheme(theme: Theme): void {
		this.currentTheme = theme;
		this.applyTheme();
		this.saveToStorage();
	}

	toggleTheme(): void {
		const newTheme = this.currentTheme === 'dark' ? 'light' : 'dark';
		this.setTheme(newTheme);
	}

	private initializeTheme(): void {
		// Try to get theme from localStorage
		const stored = localStorage.getItem('exit-debt-theme') as Theme;

		if (stored && ['light', 'dark', 'system'].includes(stored)) {
			this.currentTheme = stored;
		}

		// Detect system preference
		this.detectSystemPreference();
		this.applyTheme();
	}

	private detectSystemPreference(): void {
		if (window.matchMedia) {
			this.systemPreference = window.matchMedia('(prefers-color-scheme: dark)').matches
				? 'dark'
				: 'light';
		}
	}

	private watchSystemPreference(): void {
		if (window.matchMedia) {
			const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
			mediaQuery.addEventListener('change', (e) => {
				this.systemPreference = e.matches ? 'dark' : 'light';
				if (this.currentTheme === 'system') {
					this.applyTheme();
				}
			});
		}
	}

	private applyTheme(): void {
		const resolvedTheme = this.resolvedTheme;
		const root = document.documentElement;

		// Remove existing theme classes
		root.classList.remove('light', 'dark');

		// Add current theme class
		root.classList.add(resolvedTheme);

		// Update meta theme-color for mobile browsers
		const metaThemeColor = document.querySelector('meta[name="theme-color"]');
		if (metaThemeColor) {
			metaThemeColor.setAttribute('content', resolvedTheme === 'dark' ? '#0f172a' : '#ffffff');
		}
	}

	private saveToStorage(): void {
		localStorage.setItem('exit-debt-theme', this.currentTheme);
	}
}

export const themeStore = new ThemeStore();
