import { create } from 'zustand'
import { apiClient } from '@api/client'

export const usePaymentsStore = create((set, get) => ({
  // State
  payments: [],
  upcomingPayments: [],
  selectedPayment: null,
  isLoading: false,
  error: null,

  // Fetch payments for a specific debt
  fetchPayments: async (debtId) => {
    try {
      set({ isLoading: true, error: null })
      const payments = await apiClient.getPayments(debtId)
      set({ payments, isLoading: false })
      return payments
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Fetch upcoming payments
  fetchUpcomingPayments: async () => {
    try {
      set({ isLoading: true, error: null })
      const upcomingPayments = await apiClient.getUpcomingPayments()
      set({ upcomingPayments, isLoading: false })
      return upcomingPayments
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Create payment
  createPayment: async (debtId, paymentData) => {
    try {
      set({ isLoading: true, error: null })
      const payment = await apiClient.createPayment(debtId, paymentData)
      set((state) => ({
        payments: [...state.payments, payment],
        isLoading: false,
      }))
      return payment
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Update payment
  updatePayment: async (paymentId, paymentData) => {
    try {
      set({ isLoading: true, error: null })
      const updatedPayment = await apiClient.updatePayment(paymentId, paymentData)
      set((state) => ({
        payments: state.payments.map((p) => (p.id === paymentId ? updatedPayment : p)),
        selectedPayment:
          state.selectedPayment?.id === paymentId ? updatedPayment : state.selectedPayment,
        isLoading: false,
      }))
      return updatedPayment
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Delete payment
  deletePayment: async (paymentId) => {
    try {
      set({ isLoading: true, error: null })
      await apiClient.deletePayment(paymentId)
      set((state) => ({
        payments: state.payments.filter((p) => p.id !== paymentId),
        selectedPayment: state.selectedPayment?.id === paymentId ? null : state.selectedPayment,
        isLoading: false,
      }))
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Upload receipt photo
  uploadReceipt: async (paymentId, file) => {
    try {
      set({ isLoading: true, error: null })
      const result = await apiClient.uploadReceiptPhoto(paymentId, file)
      set({ isLoading: false })
      return result
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Set selected payment
  setSelectedPayment: (payment) => {
    set({ selectedPayment: payment })
  },

  // Clear error
  clearError: () => {
    set({ error: null })
  },
}))
