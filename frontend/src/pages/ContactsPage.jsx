import { useState } from 'react'
import { ContactsTable } from '@components/contacts/ContactsTable'
import { CreateContactModal } from '@components/contacts/CreateContactModal'

export const ContactsPage = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)

  const handleContactCreated = () => {
    setShowCreateModal(false)
  }

  return (
    <div className="mx-auto max-w-7xl space-y-6">
      {/* Header */}
      <div className="flex flex-col items-center justify-between space-y-2 md:flex-row">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Contact Management</h1>
          <p className="mt-1 text-muted-foreground">
            Organize and manage all your contacts and their information
          </p>
        </div>
        <button onClick={() => setShowCreateModal(true)} className="btn-primary">
          <svg className="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            />
          </svg>
          Create Contact
        </button>
      </div>

      {/* Contacts Table */}
      <ContactsTable />

      {/* Create Contact Modal */}
      {showCreateModal && (
        <CreateContactModal
          onContactCreated={handleContactCreated}
          onClose={() => setShowCreateModal(false)}
        />
      )}
    </div>
  )
}
