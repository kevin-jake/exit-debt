import { writable } from "svelte/store";
import { tokenManager, apiClient } from "$lib/api";

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  created_at: string;
  updated_at: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
  });

  return {
    subscribe,

    // Initialize auth state from localStorage and validate token
    async init() {
      const token = tokenManager.getToken();
      if (token) {
        try {
          // Validate token by making a request to the protected health endpoint
          const response = await apiClient.healthCheck();

          // If we get here, the token is valid
          // Since we don't have a /auth/me endpoint, we'll just set as authenticated
          // User info will need to be set when they log in
          update((state) => ({
            ...state,
            isAuthenticated: true,
            isLoading: false,
          }));
        } catch (error) {
          // Token is invalid, remove it
          console.log("Invalid token, removing from storage");
          tokenManager.removeToken();
          update((state) => ({
            ...state,
            user: null,
            isAuthenticated: false,
            isLoading: false,
          }));
        }
      } else {
        update((state) => ({
          ...state,
          isLoading: false,
        }));
      }
    },

    // Set user data after successful login
    setUser(user: User) {
      update((state) => ({
        ...state,
        user,
        isAuthenticated: true,
        isLoading: false,
      }));
    },

    // Clear user data on logout
    logout() {
      tokenManager.removeToken();
      set({
        user: null,
        isAuthenticated: false,
        isLoading: false,
      });
    },

    // Set loading state
    setLoading(isLoading: boolean) {
      update((state) => ({
        ...state,
        isLoading,
      }));
    },
  };
}

export const authStore = createAuthStore();
