import { writable } from "svelte/store";
import {
  apiClient,
  type DebtList,
  type CreateDebtListRequest,
  type UpdateDebtListRequest,
} from "$lib/api";

export interface DebtsState {
  debts: DebtList[];
  isLoading: boolean;
  error: string | null;
  selectedDebt: DebtList | null;
}

function createDebtsStore() {
  const { subscribe, set, update } = writable<DebtsState>({
    debts: [],
    isLoading: false,
    error: null,
    selectedDebt: null,
  });

  return {
    subscribe,

    async loadDebts() {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const debts = await apiClient.getDebtLists();
        update((state) => ({ ...state, debts, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to load debts",
          isLoading: false,
        }));
      }
    },

    async createDebt(debtData: CreateDebtListRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const newDebt = await apiClient.createDebtList(debtData);
        update((state) => ({
          ...state,
          debts: [...state.debts, newDebt],
          isLoading: false,
        }));
        return newDebt;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to create debt",
          isLoading: false,
        }));
        throw error;
      }
    },

    async updateDebt(id: string, debtData: UpdateDebtListRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const updatedDebt = await apiClient.updateDebtList(id, debtData);
        update((state) => ({
          ...state,
          debts: state.debts.map((debt) =>
            debt.id === id ? updatedDebt : debt
          ),
          selectedDebt:
            state.selectedDebt?.id === id ? updatedDebt : state.selectedDebt,
          isLoading: false,
        }));
        return updatedDebt;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to update debt",
          isLoading: false,
        }));
        throw error;
      }
    },

    async deleteDebt(id: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        await apiClient.deleteDebtList(id);
        update((state) => ({
          ...state,
          debts: state.debts.filter((debt) => debt.id !== id),
          selectedDebt:
            state.selectedDebt?.id === id ? null : state.selectedDebt,
          isLoading: false,
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to delete debt",
          isLoading: false,
        }));
        throw error;
      }
    },

    async getDebt(id: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const debt = await apiClient.getDebtList(id);
        update((state) => ({
          ...state,
          selectedDebt: debt,
          isLoading: false,
        }));
        return debt;
      } catch (error) {
        update((state) => ({
          ...state,
          error: error instanceof Error ? error.message : "Failed to load debt",
          isLoading: false,
        }));
        throw error;
      }
    },

    setSelectedDebt(debt: DebtList | null) {
      update((state) => ({ ...state, selectedDebt: debt }));
    },

    clearError() {
      update((state) => ({ ...state, error: null }));
    },

    reset() {
      set({
        debts: [],
        isLoading: false,
        error: null,
        selectedDebt: null,
      });
    },
  };
}

export const debtsStore = createDebtsStore();
