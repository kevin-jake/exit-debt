import { create } from 'zustand'
import { apiClient } from '@api/client'

export const useContactsStore = create((set, get) => ({
  // State
  contacts: [],
  selectedContact: null,
  isLoading: false,
  error: null,

  // Fetch all contacts
  fetchContacts: async () => {
    try {
      set({ isLoading: true, error: null })
      const contacts = await apiClient.getContacts()
      set({ contacts, isLoading: false })
      return contacts
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Fetch single contact
  fetchContact: async (id) => {
    try {
      set({ isLoading: true, error: null })
      const contact = await apiClient.getContact(id)
      set({ selectedContact: contact, isLoading: false })
      return contact
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Create contact
  createContact: async (contactData) => {
    try {
      set({ isLoading: true, error: null })
      const contact = await apiClient.createContact(contactData)
      set((state) => ({
        contacts: [...state.contacts, contact],
        isLoading: false,
      }))
      return contact
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Update contact
  updateContact: async (id, contactData) => {
    try {
      set({ isLoading: true, error: null })
      const updatedContact = await apiClient.updateContact(id, contactData)
      set((state) => ({
        contacts: state.contacts.map((c) => (c.id === id ? updatedContact : c)),
        selectedContact: state.selectedContact?.id === id ? updatedContact : state.selectedContact,
        isLoading: false,
      }))
      return updatedContact
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Delete contact
  deleteContact: async (id) => {
    try {
      set({ isLoading: true, error: null })
      await apiClient.deleteContact(id)
      set((state) => ({
        contacts: state.contacts.filter((c) => c.id !== id),
        selectedContact: state.selectedContact?.id === id ? null : state.selectedContact,
        isLoading: false,
      }))
    } catch (error) {
      set({ error: error.message, isLoading: false })
      throw error
    }
  },

  // Set selected contact
  setSelectedContact: (contact) => {
    set({ selectedContact: contact })
  },

  // Clear error
  clearError: () => {
    set({ error: null })
  },
}))

