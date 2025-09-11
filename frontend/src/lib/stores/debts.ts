import { writable } from "svelte/store";
import {
  apiClient,
  type DebtList,
  type CreateDebtListRequest,
  type UpdateDebtListRequest,
} from "../api";

// Extended debt interface for display purposes
export type Debt = DebtList & {
  contactName: string;
  remainingBalance: number;
  status: "active" | "settled" | "archived" | "overdue";
  dueDate: string;
  nextPayment: string;
};

interface DebtsState {
  debts: Debt[];
  isLoading: boolean;
  error: string | null;
  selectedDebt: Debt | null;
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
        const debtLists = await apiClient.getDebtLists();

        // Transform debt lists to include contact names and calculate derived fields
        const debts: Debt[] = debtLists.map((debtList) => {
          // Use contact information from the debt response
          const contactName = debtList.contact?.name || "Unknown Contact";

          // Use actual remaining balance from backend
          const remainingBalance = parseFloat(
            debtList.total_remaining_debt || debtList.total_amount
          );

          // Use actual status from backend
          const status =
            (debtList.status as
              | "active"
              | "settled"
              | "archived"
              | "overdue") || "active";

          // Use actual due date from backend
          const dueDate = debtList.due_date || debtList.created_at;
          const nextPayment = debtList.next_payment_date || debtList.created_at;

          return {
            ...debtList,
            contactName,
            remainingBalance,
            status,
            dueDate,
            nextPayment,
          };
        });

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
      try {
        const newDebtList = await apiClient.createDebtList(debtData);

        // Transform the new debt list to match our Debt interface
        const contactName = "Unknown Contact"; // We'll need to fetch contact details separately if needed
        const remainingBalance = parseFloat(
          newDebtList.total_remaining_debt || newDebtList.total_amount
        );
        const status =
          (newDebtList.status as
            | "active"
            | "settled"
            | "archived"
            | "overdue") || "active";
        const dueDate = newDebtList.due_date || newDebtList.created_at;
        const nextPayment =
          newDebtList.next_payment_date || newDebtList.created_at;

        const newDebt: Debt = {
          ...newDebtList,
          contactName,
          remainingBalance,
          status,
          dueDate,
          nextPayment,
        };

        update((state) => ({
          ...state,
          debts: [newDebt, ...state.debts], // Add to beginning of list
        }));

        return newDebt;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to create debt",
        }));
        throw error;
      }
    },

    async updateDebt(id: string, debtData: UpdateDebtListRequest) {
      try {
        const updatedDebtList = await apiClient.updateDebtList(id, debtData);

        // Transform the updated debt list
        const contactName = "Unknown Contact"; // We'll need to fetch contact details separately if needed
        const remainingBalance = parseFloat(
          updatedDebtList.total_remaining_debt || updatedDebtList.total_amount
        );
        const status =
          (updatedDebtList.status as
            | "active"
            | "settled"
            | "archived"
            | "overdue") || "active";
        const dueDate = updatedDebtList.due_date || updatedDebtList.created_at;
        const nextPayment =
          updatedDebtList.next_payment_date || updatedDebtList.created_at;

        const updatedDebt: Debt = {
          ...updatedDebtList,
          contactName,
          remainingBalance,
          status,
          dueDate,
          nextPayment,
        };

        update((state) => ({
          ...state,
          debts: state.debts.map((debt) =>
            debt.id === id ? updatedDebt : debt
          ),
        }));

        return updatedDebt;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to update debt",
        }));
        throw error;
      }
    },

    async deleteDebt(id: string) {
      try {
        await apiClient.deleteDebtList(id);
        update((state) => ({
          ...state,
          debts: state.debts.filter((debt) => debt.id !== id),
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to delete debt",
        }));
        throw error;
      }
    },

    setSelectedDebt(debt: Debt | null) {
      update((state) => ({ ...state, selectedDebt: debt }));
    },

    clearError() {
      update((state) => ({ ...state, error: null }));
    },

    // Helper function to get debt by ID
    getDebtById(id: string): Debt | undefined {
      let debt: Debt | undefined;
      update((state) => {
        debt = state.debts.find((d) => d.id === id);
        return state;
      });
      return debt;
    },
  };
}

export const debtsStore = createDebtsStore();
