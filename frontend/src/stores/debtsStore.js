import { create } from 'zustand'
import { apiClient } from '@api/client'

export const useDebtsStore = create((set, get) => ({
  // State
  debts: [],
  selectedDebt: null,
  isLoading: false,
  error: null,

  // Fetch all debts
  fetchDebts: async () => {
    try {
      set({ isLoading: true, error: null })
      const debts = await apiClient.getDebtLists()
      set({ debts, isLoading: false })
      return debts
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Fetch single debt
  fetchDebt: async (id) => {
    try {
      set({ isLoading: true, error: null })
      const debt = await apiClient.getDebtList(id)
      set({ selectedDebt: debt, isLoading: false })
      return debt
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Create debt
  createDebt: async (debtData) => {
    try {
      set({ isLoading: true, error: null })
      const debt = await apiClient.createDebtList(debtData)
      set((state) => ({
        debts: [...state.debts, debt],
        isLoading: false,
      }))
      return debt
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Update debt
  updateDebt: async (id, debtData) => {
    try {
      set({ isLoading: true, error: null })
      const updatedDebt = await apiClient.updateDebtList(id, debtData)
      set((state) => ({
        debts: state.debts.map((d) => (d.id === id ? updatedDebt : d)),
        selectedDebt: state.selectedDebt?.id === id ? updatedDebt : state.selectedDebt,
        isLoading: false,
      }))
      return updatedDebt
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Delete debt
  deleteDebt: async (id) => {
    try {
      set({ isLoading: true, error: null })
      await apiClient.deleteDebtList(id)
      set((state) => ({
        debts: state.debts.filter((d) => d.id !== id),
        selectedDebt: state.selectedDebt?.id === id ? null : state.selectedDebt,
        isLoading: false,
      }))
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Set selected debt
  setSelectedDebt: (debt) => {
    set({ selectedDebt: debt })
  },

  // Clear error
  clearError: () => {
    set({ error: null })
  },
}))

