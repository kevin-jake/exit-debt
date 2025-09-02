import { writable } from "svelte/store";
import {
  apiClient,
  type Payment,
  type CreatePaymentRequest,
  type PaymentSchedule,
} from "$lib/api";

export interface PaymentsState {
  payments: Payment[];
  paymentSchedules: Record<string, PaymentSchedule>;
  overdueItems: Payment[];
  dueSoonItems: Payment[];
  upcomingPayments: Payment[];
  isLoading: boolean;
  error: string | null;
  selectedPayment: Payment | null;
}

function createPaymentsStore() {
  const { subscribe, set, update } = writable<PaymentsState>({
    payments: [],
    paymentSchedules: {},
    overdueItems: [],
    dueSoonItems: [],
    upcomingPayments: [],
    isLoading: false,
    error: null,
    selectedPayment: null,
  });

  return {
    subscribe,

    async loadPayments(debtId: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const payments = await apiClient.getPayments(debtId);
        update((state) => ({ ...state, payments, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to load payments",
          isLoading: false,
        }));
      }
    },

    async createPayment(debtId: string, paymentData: CreatePaymentRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const newPayment = await apiClient.createPayment(debtId, paymentData);
        update((state) => ({
          ...state,
          payments: [...state.payments, newPayment],
          isLoading: false,
        }));
        return newPayment;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to create payment",
          isLoading: false,
        }));
        throw error;
      }
    },

    async verifyPayment(paymentId: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        await apiClient.verifyPayment(paymentId);
        update((state) => ({
          ...state,
          payments: state.payments.map((payment) =>
            payment.id === paymentId
              ? { ...payment, status: "verified" as const }
              : payment
          ),
          selectedPayment:
            state.selectedPayment?.id === paymentId
              ? { ...state.selectedPayment, status: "verified" as const }
              : state.selectedPayment,
          isLoading: false,
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to verify payment",
          isLoading: false,
        }));
        throw error;
      }
    },

    async rejectPayment(paymentId: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        await apiClient.rejectPayment(paymentId);
        update((state) => ({
          ...state,
          payments: state.payments.map((payment) =>
            payment.id === paymentId
              ? { ...payment, status: "rejected" as const }
              : payment
          ),
          selectedPayment:
            state.selectedPayment?.id === paymentId
              ? { ...state.selectedPayment, status: "rejected" as const }
              : state.selectedPayment,
          isLoading: false,
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to reject payment",
          isLoading: false,
        }));
        throw error;
      }
    },

    async uploadReceipt(paymentId: string, file: File) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        await apiClient.uploadReceipt(paymentId, file);
        update((state) => ({ ...state, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to upload receipt",
          isLoading: false,
        }));
        throw error;
      }
    },

    async loadPaymentSchedule(debtId: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const schedule = await apiClient.getPaymentSchedule(debtId);
        update((state) => ({
          ...state,
          paymentSchedules: { ...state.paymentSchedules, [debtId]: schedule },
          isLoading: false,
        }));
        return schedule;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error
              ? error.message
              : "Failed to load payment schedule",
          isLoading: false,
        }));
        throw error;
      }
    },

    async loadOverdueItems() {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const overdueItems = await apiClient.getOverdueItems();
        update((state) => ({ ...state, overdueItems, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error
              ? error.message
              : "Failed to load overdue items",
          isLoading: false,
        }));
      }
    },

    async loadDueSoonItems() {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const dueSoonItems = await apiClient.getDueSoonItems();
        update((state) => ({ ...state, dueSoonItems, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error
              ? error.message
              : "Failed to load due soon items",
          isLoading: false,
        }));
      }
    },

    async loadUpcomingPayments() {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const upcomingPayments = await apiClient.getUpcomingPayments();
        update((state) => ({ ...state, upcomingPayments, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error
              ? error.message
              : "Failed to load upcoming payments",
          isLoading: false,
        }));
      }
    },

    setSelectedPayment(payment: Payment | null) {
      update((state) => ({ ...state, selectedPayment: payment }));
    },

    clearError() {
      update((state) => ({ ...state, error: null }));
    },

    reset() {
      set({
        payments: [],
        paymentSchedules: {},
        overdueItems: [],
        dueSoonItems: [],
        upcomingPayments: [],
        isLoading: false,
        error: null,
        selectedPayment: null,
      });
    },
  };
}

export const paymentsStore = createPaymentsStore();
