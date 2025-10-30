import { useEffect } from 'react'

export const ContactDetailsModal = ({ contact, onClose, onEdit, onDelete }) => {
  useEffect(() => {
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = 'auto'
    }
  }, [])

  useEffect(() => {
    const handleEscape = (e) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleEscape)
    return () => window.removeEventListener('keydown', handleEscape)
  }, [onClose])

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getInitials = (name) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  return (
    <div
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onClick={onClose}
    >
      <div
        className="my-8 w-full max-w-lg rounded-xl bg-card shadow-medium"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-xl font-semibold text-foreground">Contact Details</h2>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground">
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="p-6">
          {/* Contact Header */}
          <div className="mb-6 flex items-center space-x-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary">
              <span className="text-xl font-medium text-primary-foreground">
                {getInitials(contact.name)}
              </span>
            </div>
            <div>
              <h3 className="text-2xl font-bold text-foreground">{contact.name}</h3>
              <span className="inline-flex rounded-full bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
                Regular Contact
              </span>
            </div>
          </div>

          {/* Contact Information */}
          <div className="space-y-4">
            {contact.email && (
              <div className="flex items-start space-x-3">
                <svg
                  className="mt-0.5 h-5 w-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Email</p>
                  <p className="text-foreground">{contact.email}</p>
                </div>
              </div>
            )}

            {contact.phone && (
              <div className="flex items-start space-x-3">
                <svg
                  className="mt-0.5 h-5 w-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Phone</p>
                  <p className="text-foreground">{contact.phone}</p>
                </div>
              </div>
            )}

            {contact.notes && (
              <div className="flex items-start space-x-3">
                <svg
                  className="mt-0.5 h-5 w-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Notes</p>
                  <p className="text-foreground">{contact.notes}</p>
                </div>
              </div>
            )}

            <div className="flex items-start space-x-3">
              <svg
                className="mt-0.5 h-5 w-5 text-muted-foreground"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              <div>
                <p className="text-sm font-medium text-muted-foreground">Created</p>
                <p className="text-foreground">{formatDate(contact.created_at)}</p>
              </div>
            </div>

            {contact.updated_at !== contact.created_at && (
              <div className="flex items-start space-x-3">
                <svg
                  className="mt-0.5 h-5 w-5 text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-muted-foreground">Last Updated</p>
                  <p className="text-foreground">{formatDate(contact.updated_at)}</p>
                </div>
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="mt-6 flex items-center justify-end space-x-3 border-t border-border pt-4">
            <button onClick={onDelete} className="btn-destructive">
              Delete
            </button>
            <button onClick={onEdit} className="btn-primary">
              Edit
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
