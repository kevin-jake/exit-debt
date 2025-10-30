import { useState } from 'react'
import { useDebtsStore } from '@stores/debtsStore'
import { useNotificationsStore } from '@stores/notificationsStore'
import { formatCurrency } from '@utils/formatters'

export const DeleteDebtModal = ({ debt, onConfirm, onClose }) => {
  const deleteDebt = useDebtsStore((state) => state.deleteDebt)
  const { success, error } = useNotificationsStore()
  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = async () => {
    setIsDeleting(true)
    try {
      await deleteDebt(debt.id)
      success('Debt deleted successfully')
      onConfirm()
    } catch (err) {
      error(err.message || 'Failed to delete debt')
    } finally {
      setIsDeleting(false)
    }
  }

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget && !isDeleting) {
      onClose()
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onClick={handleOverlayClick}
    >
      <div className="card my-8 w-full max-w-md overflow-hidden">
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold text-foreground">Delete Debt</h2>
            <button
              onClick={onClose}
              disabled={isDeleting}
              className="text-muted-foreground transition-colors hover:text-foreground disabled:cursor-not-allowed"
            >
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
        </div>

        <div className="p-6">
          <div className="mb-4 flex items-center justify-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
              <svg
                className="h-8 w-8 text-destructive"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
            </div>
          </div>

          <div className="mb-6 text-center">
            <p className="mb-2 text-muted-foreground">
              Are you sure you want to delete this debt? This action cannot be undone.
            </p>
            <div className="mt-4 rounded-lg border border-border bg-muted/50 p-4">
              <div className="mb-2 flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Contact:</span>
                <span className="font-medium text-foreground">
                  {debt.contact?.name || 'Unknown Contact'}
                </span>
              </div>
              <div className="mb-2 flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Description:</span>
                <span className="font-medium text-foreground">
                  {debt.description || 'No description'}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Amount:</span>
                <span
                  className={`font-semibold ${
                    debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'
                  }`}
                >
                  {formatCurrency(parseFloat(debt.total_amount || 0))}
                </span>
              </div>
            </div>
          </div>

          <div className="flex space-x-3">
            <button
              type="button"
              onClick={onClose}
              disabled={isDeleting}
              className="btn-secondary flex-1"
            >
              Cancel
            </button>
            <button onClick={handleDelete} disabled={isDeleting} className="btn-destructive flex-1">
              {isDeleting ? (
                <span className="flex items-center justify-center">
                  <svg className="mr-2 h-4 w-4 animate-spin" viewBox="0 0 24 24">
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                      fill="none"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                  Deleting...
                </span>
              ) : (
                'Delete Debt'
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
