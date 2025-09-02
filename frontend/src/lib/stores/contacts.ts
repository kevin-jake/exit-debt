import { writable } from "svelte/store";
import {
  apiClient,
  type Contact,
  type CreateContactRequest,
  type UpdateContactRequest,
} from "$lib/api";

export interface ContactsState {
  contacts: Contact[];
  isLoading: boolean;
  error: string | null;
  selectedContact: Contact | null;
}

function createContactsStore() {
  const { subscribe, set, update } = writable<ContactsState>({
    contacts: [],
    isLoading: false,
    error: null,
    selectedContact: null,
  });

  return {
    subscribe,

    async loadContacts() {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const contacts = await apiClient.getContacts();
        console.log(contacts);
        update((state) => ({ ...state, contacts, isLoading: false }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to load contacts",
          isLoading: false,
        }));
      }
    },

    async createContact(contactData: CreateContactRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const newContact = await apiClient.createContact(contactData);
        update((state) => ({
          ...state,
          contacts: [...state.contacts, newContact],
          isLoading: false,
        }));
        return newContact;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to create contact",
          isLoading: false,
        }));
        throw error;
      }
    },

    async updateContact(id: string, contactData: UpdateContactRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const updatedContact = await apiClient.updateContact(id, contactData);
        update((state) => ({
          ...state,
          contacts: state.contacts.map((contact) =>
            contact.id === id ? updatedContact : contact
          ),
          selectedContact:
            state.selectedContact?.id === id
              ? updatedContact
              : state.selectedContact,
          isLoading: false,
        }));
        return updatedContact;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to update contact",
          isLoading: false,
        }));
        throw error;
      }
    },

    async deleteContact(id: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        await apiClient.deleteContact(id);
        update((state) => ({
          ...state,
          contacts: state.contacts.filter((contact) => contact.id !== id),
          selectedContact:
            state.selectedContact?.id === id ? null : state.selectedContact,
          isLoading: false,
        }));
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to delete contact",
          isLoading: false,
        }));
        throw error;
      }
    },

    async getContact(id: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));
      try {
        const contact = await apiClient.getContact(id);
        update((state) => ({
          ...state,
          selectedContact: contact,
          isLoading: false,
        }));
        return contact;
      } catch (error) {
        update((state) => ({
          ...state,
          error:
            error instanceof Error ? error.message : "Failed to load contact",
          isLoading: false,
        }));
        throw error;
      }
    },

    setSelectedContact(contact: Contact | null) {
      update((state) => ({ ...state, selectedContact: contact }));
    },

    clearError() {
      update((state) => ({ ...state, error: null }));
    },

    reset() {
      set({
        contacts: [],
        isLoading: false,
        error: null,
        selectedContact: null,
      });
    },
  };
}

export const contactsStore = createContactsStore();
