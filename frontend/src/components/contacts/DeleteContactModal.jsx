import { useState, useEffect } from 'react'
import { useContactsStore } from '@stores/contactsStore'
import { useDebtsStore } from '@stores/debtsStore'
import { useNotificationsStore } from '@stores/notificationsStore'

export const DeleteContactModal = ({ contact, onConfirm, onClose }) => {
  const [isLoading, setIsLoading] = useState(false)
  const deleteContact = useContactsStore((state) => state.deleteContact)
  const debts = useDebtsStore((state) => state.debts)
  const showSuccess = useNotificationsStore((state) => state.success)
  const showError = useNotificationsStore((state) => state.error)

  // Get debts associated with this contact
  const associatedDebts = debts.filter((debt) => debt.contact_id === contact.id)
  const hasAssociatedDebts = associatedDebts.length > 0
  const displayDebts = associatedDebts.slice(0, 10)
  const remainingDebtsCount = Math.max(0, associatedDebts.length - 10)

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

  const handleDelete = async () => {
    setIsLoading(true)

    try {
      await deleteContact(contact.id)
      showSuccess(`Successfully deleted contact "${contact.name}"`)
      onConfirm()
    } catch (error) {
      showError(error.message || 'Failed to delete contact. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onClick={onClose}
    >
      <div
        className="my-8 w-full max-w-md overflow-hidden rounded-xl bg-card shadow-medium"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-xl font-semibold text-foreground">Delete Contact</h2>
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
          <div className="mb-6">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
              <svg
                className="h-6 w-6 text-destructive"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3l-6.928-12c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
            </div>
            <h3 className="mb-2 text-center text-lg font-medium text-foreground">
              {hasAssociatedDebts ? 'Cannot Delete Contact' : 'Are you sure?'}
            </h3>
            <p className="text-center text-sm text-muted-foreground">
              {hasAssociatedDebts ? (
                <>
                  This contact cannot be deleted because there are{' '}
                  <span className="font-medium">{associatedDebts.length}</span> debt list
                  {associatedDebts.length === 1 ? '' : 's'} associated with{' '}
                  <span className="font-medium">"{contact.name}"</span>. Please delete or reassign
                  the following debt list{associatedDebts.length === 1 ? '' : 's'} first:
                </>
              ) : (
                <>
                  This will permanently delete <span className="font-medium">"{contact.name}"</span>{' '}
                  and all associated data. This action cannot be undone.
                </>
              )}
            </p>
          </div>

          {hasAssociatedDebts && (
            <div className="mb-6 max-h-60 overflow-y-auto rounded-lg border border-border bg-muted/30 p-4">
              <ul className="space-y-2">
                {displayDebts.map((debt) => (
                  <li
                    key={debt.id}
                    className="flex items-start justify-between rounded-md bg-card p-3 text-sm"
                  >
                    <div className="flex-1">
                      <div className="font-medium text-foreground">
                        {debt.debt_type === 'owed_to_me' ? '↓ Owed to me' : '↑ I owe'}
                      </div>
                      {debt.description && (
                        <div className="mt-1 line-clamp-1 text-xs text-muted-foreground">
                          {debt.description}
                        </div>
                      )}
                    </div>
                    <div className="ml-3 text-right">
                      <div className="font-semibold text-foreground">
                        {debt.currency} {parseFloat(debt.total_amount).toLocaleString()}
                      </div>
                      <div className="text-xs capitalize text-muted-foreground">{debt.status}</div>
                    </div>
                  </li>
                ))}
              </ul>
              {remainingDebtsCount > 0 && (
                <div className="mt-3 text-center text-xs text-muted-foreground">
                  and {remainingDebtsCount} more debt list{remainingDebtsCount === 1 ? '' : 's'}...
                </div>
              )}
            </div>
          )}

          <div className="flex items-center justify-end space-x-3">
            <button type="button" onClick={onClose} className="btn-secondary" disabled={isLoading}>
              {hasAssociatedDebts ? 'Close' : 'Cancel'}
            </button>
            {!hasAssociatedDebts && (
              <button
                onClick={handleDelete}
                className="btn-destructive disabled:cursor-not-allowed disabled:opacity-50"
                disabled={isLoading}
              >
                {isLoading ? 'Deleting...' : 'Delete Contact'}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
