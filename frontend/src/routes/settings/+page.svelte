<script lang="ts">
  import { onMount } from 'svelte';
  import { settingsStore } from '$lib/stores/settings.svelte.js';
  
  // Local state
  let showSaveConfirmation = false;
  let showResetConfirmation = false;
  
  // Store subscriptions
  let settings: any = {};
  let isLoading = false;
  let hasUnsavedChanges = false;
  
  // Subscribe to store changes
  settingsStore.subscribe((value: any) => settings = value);
  settingsStore.isLoading.subscribe((value: boolean) => isLoading = value);
  settingsStore.hasUnsavedChanges.subscribe((value: boolean) => hasUnsavedChanges = value);
  
  // Currency options
  const currencies = [
    { code: 'USD', symbol: '$', name: 'US Dollar', example: '$1,234.56' },
    { code: 'EUR', symbol: '€', name: 'Euro', example: '€1.234,56' },
    { code: 'GBP', symbol: '£', name: 'British Pound', example: '£1,234.56' },
    { code: 'JPY', symbol: '¥', name: 'Japanese Yen', example: '¥123,456' },
    { code: 'CAD', symbol: 'C$', name: 'Canadian Dollar', example: 'C$1,234.56' },
    { code: 'AUD', symbol: 'A$', name: 'Australian Dollar', example: 'A$1,234.56' },
    { code: 'CHF', symbol: 'CHF', name: 'Swiss Franc', example: 'CHF 1\'234.56' },
    { code: 'CNY', symbol: '¥', name: 'Chinese Yuan', example: '¥1,234.56' }
  ];
  
  const dateFormats = [
    { value: 'MM/DD/YYYY', label: 'MM/DD/YYYY', example: '12/25/2024' },
    { value: 'DD/MM/YYYY', label: 'DD/MM/YYYY', example: '25/12/2024' },
    { value: 'YYYY-MM-DD', label: 'YYYY-MM-DD', example: '2024-12-25' }
  ];
  
  const timezones = Intl.supportedValuesOf('timeZone');
  
  // Methods
  function handleSettingChange() {
    // This will be handled by updateSetting method
  }
  
  async function handleSave() {
    await settingsStore.saveSettings();
    showSaveConfirmation = true;
    setTimeout(() => showSaveConfirmation = false, 3000);
  }
  
  function handleReset() {
    showResetConfirmation = true;
  }
  
  function confirmReset() {
    settingsStore.resetToDefaults();
    showResetConfirmation = false;
  }
  

  
  function handleCurrencyChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    if (target) {
      const selectedCurrency = currencies.find(c => c.code === target.value);
      if (selectedCurrency) {
        settingsStore.updateSetting('currency.code', selectedCurrency.code);
        settingsStore.updateSetting('currency.symbol', selectedCurrency.symbol);
      }
    }
  }
  
  function handleDueDateReminderChange(days: number, checked: boolean) {
    const currentReminders = [...(settings.notifications?.email?.dueDateReminders || [])];
    if (checked && !currentReminders.includes(days)) {
      currentReminders.push(days);
    } else if (!checked && currentReminders.includes(days)) {
      const index = currentReminders.indexOf(days);
      currentReminders.splice(index, 1);
    }
    settingsStore.updateSetting('notifications.email.dueDateReminders', currentReminders);
  }
  
  onMount(async () => {
    await settingsStore.loadSettings();
  });
</script>

<svelte:head>
  <title>Settings - Exit-Debt</title>
  <meta name="description" content="Configure your application preferences and basic settings" />
</svelte:head>

<div class="max-w-4xl mx-auto p-4 sm:p-6">
  <!-- Header -->
  <div class="mb-6 sm:mb-8">
    <h1 class="text-2xl sm:text-3xl font-bold text-foreground mb-2">Settings & Preferences</h1>
    <p class="text-muted-foreground">Customize your Exit-Debt experience with basic preferences</p>
  </div>
  
  <!-- Save/Reset Actions -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
    <div class="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3">
      <button
        on:click={handleSave}
        disabled={!hasUnsavedChanges || isLoading}
        class="btn-primary w-full sm:w-auto"
      >
        {isLoading ? 'Saving...' : 'Save Changes'}
      </button>
      
      <button
        on:click={handleReset}
        class="btn-secondary w-full sm:w-auto"
      >
        Reset to Defaults
      </button>
    </div>
    
    {#if showSaveConfirmation}
      <div class="text-green-600 bg-green-50 dark:bg-green-900/20 px-4 py-2 rounded-lg text-sm">
        ✓ Settings saved successfully
      </div>
    {/if}
  </div>
  
  <!-- Settings Sections -->
  <div class="space-y-4 sm:space-y-6">
    <!-- General Settings -->
    <div class="card">
      <div class="p-4 sm:p-6">
        <div class="flex items-center space-x-3 mb-4 sm:mb-6">
          <svg class="w-5 h-5 sm:w-6 sm:h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <h2 class="text-lg sm:text-xl font-semibold text-foreground">General Settings</h2>
        </div>
        
        <div class="space-y-4 sm:space-y-6">
          <!-- Currency Settings -->
          <div>
            <h3 class="text-base sm:text-lg font-medium text-foreground mb-3 sm:mb-4">Currency Preferences</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
              <div>
                <label class="label">Default Currency</label>
                <select
                  value={settings.currency?.code}
                  on:change={handleCurrencyChange}
                  class="input"
                >
                  {#each currencies as currency}
                    <option value={currency.code}>
                      {currency.code} - {currency.name} ({currency.symbol})
                    </option>
                  {/each}
                </select>
                {#if settings.currency?.code}
                  {@const selectedCurrency = currencies.find(c => c.code === settings.currency.code)}
                  <p class="text-sm text-muted-foreground mt-2">
                    Example: {selectedCurrency?.example}
                  </p>
                {/if}
              </div>
              
              <div>
                <label class="label">Decimal Places</label>
                <select
                  value={settings.currency?.decimalPlaces}
                  on:change={(e) => {
                    const target = e.target as HTMLSelectElement;
                    if (target && target.value) {
                      settingsStore.updateSetting('currency.decimalPlaces', parseInt(target.value));
                    }
                  }}
                  class="input"
                >
                  <option value={0}>0 (Whole numbers)</option>
                  <option value={2}>2 (Standard)</option>
                  <option value={3}>3 (High precision)</option>
                </select>
              </div>
            </div>
          </div>
          
          <!-- Date and Time Settings -->
          <div>
            <h3 class="text-base sm:text-lg font-medium text-foreground mb-3 sm:mb-4">Date & Time</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
              <div>
                <label class="label">Date Format</label>
                <select
                  value={settings.dateFormat}
                  on:change={(e) => {
                    const target = e.target as HTMLSelectElement;
                    if (target && target.value) {
                      settingsStore.updateSetting('dateFormat', target.value);
                    }
                  }}
                  class="input"
                >
                  {#each dateFormats as format}
                    <option value={format.value}>{format.label} - {format.example}</option>
                  {/each}
                </select>
              </div>
              
              <div>
                <label class="label">Time Zone</label>
                <select
                  value={settings.timezone}
                  on:change={(e) => {
                    const target = e.target as HTMLSelectElement;
                    if (target && target.value) {
                      settingsStore.updateSetting('timezone', target.value);
                    }
                  }}
                  class="input"
                >
                  {#each timezones as tz}
                    <option value={tz}>{tz}</option>
                  {/each}
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Basic Notifications -->
    <div class="card">
      <div class="p-4 sm:p-6">
        <div class="flex items-center space-x-3 mb-4 sm:mb-6">
          <svg class="w-5 h-5 sm:w-6 sm:h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM4 19h6v-6H4v6zM4 5h6V4a1 1 0 00-1-1H5a1 1 0 00-1 1v1zM4 19h6v-6H4v6zM4 5h6V4a1 1 0 00-1-1H5a1 1 0 00-1 1v1z" />
          </svg>
          <h2 class="text-lg sm:text-xl font-semibold text-foreground">Email Notifications</h2>
        </div>
        
        <div class="space-y-4 sm:space-y-6">
          <div class="space-y-4">
            <label class="flex items-center space-x-3">
              <input
                type="checkbox"
                checked={settings.notifications?.email?.enabled}
                                  on:change={(e) => {
                    const target = e.target as HTMLInputElement;
                    if (target && typeof target.checked === 'boolean') {
                      settingsStore.updateSetting('notifications.email.enabled', target.checked);
                    }
                  }}
                class="rounded border-border text-primary focus:ring-primary"
              />
              <span>Enable email notifications</span>
            </label>
            
            {#if settings.notifications?.email?.enabled}
              <div class="ml-6 space-y-3">
                <div>
                  <label class="label">Due Date Reminders</label>
                  <div class="flex flex-wrap gap-2">
                    {#each [1, 3, 7, 14] as days}
                      <label class="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          checked={settings.notifications?.email?.dueDateReminders?.includes(days)}
                          on:change={(e) => {
                            const target = e.target as HTMLInputElement;
                            if (target && typeof target.checked === 'boolean') {
                              handleDueDateReminderChange(days, target.checked);
                            }
                          }}
                          class="rounded border-border text-primary focus:ring-primary"
                        />
                        <span class="text-sm">{days} day{days === 1 ? '' : 's'} before</span>
                      </label>
                    {/each}
                  </div>
                </div>
                
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    checked={settings.notifications?.email?.paymentConfirmations}
                    on:change={(e) => {
                      const target = e.target as HTMLInputElement;
                      if (target && typeof target.checked === 'boolean') {
                        settingsStore.updateSetting('notifications.email.paymentConfirmations', target.checked);
                      }
                    }}
                    class="rounded border-border text-primary focus:ring-primary"
                  />
                  <span>Payment confirmations</span>
                </label>
                
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    checked={settings.notifications?.email?.newDebtLists}
                    on:change={(e) => {
                      const target = e.target as HTMLInputElement;
                      if (target && typeof target.checked === 'boolean') {
                        settingsStore.updateSetting('notifications.email.newDebtLists', target.checked);
                      }
                    }}
                    class="rounded border-border text-primary focus:ring-primary"
                  />
                  <span>New debt list notifications</span>
                </label>
                
                <div>
                  <label class="label">Summary Frequency</label>
                  <select
                    value={settings.notifications?.email?.summaries}
                    on:change={(e) => {
                      const target = e.target as HTMLSelectElement;
                      if (target && target.value) {
                        settingsStore.updateSetting('notifications.email.summaries', target.value);
                      }
                    }}
                    class="input"
                  >
                    <option value="never">Never</option>
                    <option value="weekly">Weekly</option>
                    <option value="monthly">Monthly</option>
                  </select>
                </div>
                
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    checked={settings.notifications?.email?.marketing}
                    on:change={(e) => {
                      const target = e.target as HTMLInputElement;
                      if (target && typeof target.checked === 'boolean') {
                        settingsStore.updateSetting('notifications.email.marketing', target.checked);
                      }
                    }}
                    class="rounded border-border text-primary focus:ring-primary"
                  />
                  <span>Marketing emails</span>
                </label>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
    
    <!-- Display Preferences -->
    <div class="card">
      <div class="p-4 sm:p-6">
        <div class="flex items-center space-x-3 mb-4 sm:mb-6">
          <svg class="w-5 h-5 sm:w-6 sm:h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
          </svg>
          <h2 class="text-lg sm:text-xl font-semibold text-foreground">Display Preferences</h2>
        </div>
        
        <div class="space-y-4 sm:space-y-6">
          <div>
            <label class="label">Theme</label>
            <select
              value={settings.display?.theme}
              on:change={(e) => {
                const target = e.target as HTMLSelectElement;
                if (target && target.value) {
                  settingsStore.updateSetting('display.theme', target.value);
                }
              }}
              class="input"
            >
              <option value="light">Light</option>
              <option value="dark">Dark</option>
              <option value="system">System</option>
            </select>
            <p class="text-sm text-muted-foreground mt-2">
              Choose your preferred theme appearance
            </p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Data Management -->
    <div class="card">
      <div class="p-4 sm:p-6">
        <div class="flex items-center space-x-3 mb-4 sm:mb-6">
          <svg class="w-5 h-5 sm:w-6 sm:h-6 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8 1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8-4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
          </svg>
          <h2 class="text-lg sm:text-xl font-semibold text-foreground">Data Management</h2>
        </div>
        
        <div class="space-y-4 sm:space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6">
            <div>
              <label class="label">Auto Backup</label>
              <select
                value={settings.dataManagement?.backupFrequency}
                on:change={(e) => {
                  const target = e.target as HTMLSelectElement;
                  if (target && target.value) {
                    settingsStore.updateSetting('dataManagement.backupFrequency', target.value);
                  }
                }}
                class="input"
              >
                <option value="never">Never</option>
                <option value="daily">Daily</option>
                <option value="weekly">Weekly</option>
                <option value="monthly">Monthly</option>
              </select>
            </div>
          </div>
          
          <div class="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3">
            <button class="btn-secondary w-full sm:w-auto">
              Export Data
            </button>
            <button class="btn-secondary w-full sm:w-auto">
              Import Data
            </button>
            <button class="btn-secondary w-full sm:w-auto">
              Backup Now
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Reset Confirmation Modal -->
{#if showResetConfirmation}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
    <div class="bg-background border border-border rounded-lg p-6 max-w-md w-full">
      <h3 class="text-lg font-semibold text-foreground mb-2">Reset to Defaults</h3>
      <p class="text-muted-foreground mb-6">
        This will reset all your settings to their default values. This action cannot be undone.
      </p>
      <div class="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3">
        <button
          on:click={() => showResetConfirmation = false}
          class="btn-secondary flex-1"
        >
          Cancel
        </button>
        <button
          on:click={confirmReset}
          class="btn-danger flex-1"
        >
          Reset Settings
        </button>
      </div>
    </div>
  </div>
{/if}
