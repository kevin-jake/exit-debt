// Export all stores for easier imports
export { authStore } from "./auth";
export { contactsStore } from "./contacts";
export { debtsStore } from "./debts";
export { paymentsStore } from "./payments";
export { notificationsStore } from "./notifications";

// Export types
export type { User } from "./auth";
export type { ContactsState } from "./contacts";
export type { DebtsState } from "./debts";
export type { PaymentsState } from "./payments";
export type { Notification, NotificationsState } from "./notifications";
