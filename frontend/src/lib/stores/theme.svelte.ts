import { browser } from '$app/environment';
import { getTheme, setTheme } from '$lib/utils';

class ThemeStore {
	theme = $state<'dark' | 'light'>(getTheme());

	constructor() {
		if (browser) {
			// Initialize theme on mount
			this.theme = getTheme();
			setTheme(this.theme);
		}
	}

	toggle() {
		this.theme = this.theme === 'dark' ? 'light' : 'dark';
		setTheme(this.theme);
	}

	set(newTheme: 'dark' | 'light') {
		this.theme = newTheme;
		setTheme(this.theme);
	}
}

export const themeStore = new ThemeStore();
