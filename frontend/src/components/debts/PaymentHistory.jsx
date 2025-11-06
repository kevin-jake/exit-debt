import { useState } from 'react'
import { formatCurrency, formatDate } from '@utils/formatters'
import { apiClient } from '@api/client'

export const PaymentHistory = ({
  debtPayments,
  paymentsLoading,
  showAddPayment,
  setShowAddPayment,
  newPayment,
  setNewPayment,
  receiptFile,
  setReceiptFile,
  isSubmitting,
  onAddPayment,
  onFileChange,
  onViewReceipt,
  debtType,
  onPaymentStatusChange,
}) => {
  const [verifyingPaymentId, setVerifyingPaymentId] = useState(null)
  const [rejectingPaymentId, setRejectingPaymentId] = useState(null)
  const [deletingPaymentId, setDeletingPaymentId] = useState(null)

  const handleCancel = () => {
    setShowAddPayment(false)
    setNewPayment({
      payment_date: new Date().toISOString().split('T')[0],
      amount: '',
      description: '',
    })
    setReceiptFile(null)
  }

  const handleVerifyPayment = async (payment) => {
    try {
      setVerifyingPaymentId(payment.id)
      await apiClient.verifyPayment(payment.id)
      // Refresh the payments list
      if (onPaymentStatusChange) {
        await onPaymentStatusChange()
      }
    } catch (error) {
      console.error('Error verifying payment:', error)
      alert(`Failed to verify payment: ${error.message}`)
    } finally {
      setVerifyingPaymentId(null)
    }
  }

  const handleRejectPayment = async (payment) => {
    try {
      setRejectingPaymentId(payment.id)
      await apiClient.rejectPayment(payment.id)
      // Refresh the payments list
      if (onPaymentStatusChange) {
        await onPaymentStatusChange()
      }
    } catch (error) {
      console.error('Error rejecting payment:', error)
      alert(`Failed to reject payment: ${error.message}`)
    } finally {
      setRejectingPaymentId(null)
    }
  }

  const handleDeletePayment = async (payment) => {
    if (
      !confirm(
        `Are you sure you want to delete this payment of ${formatCurrency(parseFloat(payment.amount || 0))}?`
      )
    ) {
      return
    }

    try {
      setDeletingPaymentId(payment.id)
      await apiClient.deletePayment(payment.id)
      // Refresh the payments list
      if (onPaymentStatusChange) {
        await onPaymentStatusChange()
      }
    } catch (error) {
      console.error('Error deleting payment:', error)
      alert(`Failed to delete payment: ${error.message}`)
    } finally {
      setDeletingPaymentId(null)
    }
  }

  return (
    <div className="rounded-lg border border-border">
      <div className="border-b border-border p-4">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-medium text-muted-foreground">
            Payment History ({debtPayments.length})
          </div>
          <button
            onClick={() => setShowAddPayment(!showAddPayment)}
            className="text-sm text-primary hover:text-primary/80"
          >
            {showAddPayment ? '− Cancel' : '+ Add Payment'}
          </button>
        </div>

        {/* Add Payment Form */}
        {showAddPayment && (
          <form onSubmit={onAddPayment} className="mt-3 space-y-3 rounded-lg bg-muted/30 p-3">
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div>
                <label className="mb-1 block text-xs font-medium text-muted-foreground">
                  Payment Date
                </label>
                <input
                  type="date"
                  value={newPayment.payment_date}
                  onChange={(e) => setNewPayment({ ...newPayment, payment_date: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="mb-1 block text-xs font-medium text-muted-foreground">
                  Amount
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={newPayment.amount}
                  onChange={(e) => setNewPayment({ ...newPayment, amount: e.target.value })}
                  className="input"
                  placeholder="0.00"
                  required
                />
              </div>
            </div>
            <div>
              <label className="mb-1 block text-xs font-medium text-muted-foreground">
                Description (Optional)
              </label>
              <input
                type="text"
                value={newPayment.description}
                onChange={(e) => setNewPayment({ ...newPayment, description: e.target.value })}
                className="input"
                placeholder="Payment description..."
              />
            </div>
            <div>
              <label className="mb-1 block text-xs font-medium text-muted-foreground">
                Receipt Photo (Optional)
              </label>
              <input
                type="file"
                accept="image/jpeg,image/png,image/jpg,image/webp"
                onChange={onFileChange}
                className="input"
              />
              {receiptFile && (
                <div className="mt-1 text-xs text-muted-foreground">
                  Selected: {receiptFile.name}
                </div>
              )}
            </div>
            <div className="flex justify-end space-x-2">
              <button
                type="button"
                onClick={handleCancel}
                className="btn-secondary"
                disabled={isSubmitting}
              >
                Cancel
              </button>
              <button type="submit" className="btn-primary" disabled={isSubmitting}>
                {isSubmitting ? 'Adding...' : 'Add Payment'}
              </button>
            </div>
          </form>
        )}
      </div>

      {paymentsLoading ? (
        <div className="py-8 text-center text-sm text-muted-foreground">Loading payments...</div>
      ) : debtPayments.length === 0 ? (
        <div className="py-8 text-center text-sm text-muted-foreground">
          No payments recorded yet
        </div>
      ) : (
        <>
          {/* Mobile Card View */}
          <div className="space-y-3 p-4 md:hidden">
            {debtPayments
              .sort((a, b) => new Date(b.payment_date) - new Date(a.payment_date))
              .map((payment) => (
                <div
                  key={payment.id}
                  className="rounded-lg border border-border bg-card p-4 hover:bg-muted/30"
                >
                  <div className="mb-3 flex items-start justify-between">
                    <div>
                      <div className="text-lg font-semibold text-foreground">
                        {formatCurrency(parseFloat(payment.amount || 0))}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {formatDate(payment.payment_date)}
                      </div>
                    </div>
                    <span
                      className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                        payment.status === 'verified'
                          ? 'bg-success/10 text-success'
                          : payment.status === 'pending'
                            ? 'bg-warning/10 text-warning'
                            : payment.status === 'disputed'
                              ? 'bg-destructive/10 text-destructive'
                              : 'bg-muted/50 text-muted-foreground'
                      }`}
                    >
                      {payment.status || 'pending'}
                    </span>
                  </div>

                  {payment.description && (
                    <div className="mb-3">
                      <div className="text-xs text-muted-foreground">Description</div>
                      <div className="text-sm text-foreground">{payment.description}</div>
                    </div>
                  )}

                  <div className="flex items-center justify-between border-t border-border pt-3">
                    <div className="flex items-center gap-3">
                      {payment.receipt_photo_url ? (
                        <button
                          onClick={() => onViewReceipt(payment)}
                          className="inline-flex items-center gap-1 text-xs text-primary hover:text-primary/80"
                          title="View receipt"
                        >
                          <svg
                            className="h-4 w-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth="2"
                              d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                            />
                          </svg>
                          Receipt
                        </button>
                      ) : (
                        <span className="text-xs text-muted-foreground">No receipt</span>
                      )}
                    </div>

                    {debtType === 'owed_to_me' && (
                      <div className="flex items-center gap-2">
                        {payment.status === 'pending' && (
                          <>
                            <button
                              onClick={() => handleVerifyPayment(payment)}
                              disabled={
                                verifyingPaymentId === payment.id ||
                                rejectingPaymentId === payment.id ||
                                deletingPaymentId === payment.id
                              }
                              className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-success transition-colors hover:bg-success/10 disabled:cursor-not-allowed disabled:opacity-50"
                              title="Verify payment"
                            >
                              {verifyingPaymentId === payment.id ? (
                                <svg
                                  className="h-4 w-4 animate-spin"
                                  fill="none"
                                  viewBox="0 0 24 24"
                                >
                                  <circle
                                    className="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                  />
                                  <path
                                    className="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                  />
                                </svg>
                              ) : (
                                <svg
                                  className="h-5 w-5"
                                  fill="none"
                                  stroke="currentColor"
                                  viewBox="0 0 24 24"
                                >
                                  <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    d="M5 13l4 4L19 7"
                                  />
                                </svg>
                              )}
                            </button>
                            <button
                              onClick={() => handleRejectPayment(payment)}
                              disabled={
                                verifyingPaymentId === payment.id ||
                                rejectingPaymentId === payment.id ||
                                deletingPaymentId === payment.id
                              }
                              className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-destructive transition-colors hover:bg-destructive/10 disabled:cursor-not-allowed disabled:opacity-50"
                              title="Reject payment"
                            >
                              {rejectingPaymentId === payment.id ? (
                                <svg
                                  className="h-4 w-4 animate-spin"
                                  fill="none"
                                  viewBox="0 0 24 24"
                                >
                                  <circle
                                    className="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                  />
                                  <path
                                    className="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                  />
                                </svg>
                              ) : (
                                <svg
                                  className="h-5 w-5"
                                  fill="none"
                                  stroke="currentColor"
                                  viewBox="0 0 24 24"
                                >
                                  <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    d="M6 18L18 6M6 6l12 12"
                                  />
                                </svg>
                              )}
                            </button>
                          </>
                        )}
                        <button
                          onClick={() => handleDeletePayment(payment)}
                          disabled={
                            verifyingPaymentId === payment.id ||
                            rejectingPaymentId === payment.id ||
                            deletingPaymentId === payment.id
                          }
                          className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive disabled:cursor-not-allowed disabled:opacity-50"
                          title="Delete payment"
                        >
                          {deletingPaymentId === payment.id ? (
                            <svg className="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                              <circle
                                className="opacity-25"
                                cx="12"
                                cy="12"
                                r="10"
                                stroke="currentColor"
                                strokeWidth="4"
                              />
                              <path
                                className="opacity-75"
                                fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                              />
                            </svg>
                          ) : (
                            <svg
                              className="h-5 w-5"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                              />
                            </svg>
                          )}
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
          </div>

          {/* Desktop Table View */}
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full">
              <thead className="border-b border-border bg-muted/50">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Date
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Amount
                  </th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-muted-foreground">
                    Description
                  </th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-muted-foreground">
                    Status
                  </th>
                  <th className="px-4 py-3 text-center text-xs font-medium text-muted-foreground">
                    Receipt
                  </th>
                  {debtType === 'owed_to_me' && (
                    <th className="px-4 py-3 text-center text-xs font-medium text-muted-foreground">
                      Actions
                    </th>
                  )}
                </tr>
              </thead>
              <tbody className="divide-y divide-border">
                {debtPayments
                  .sort((a, b) => new Date(b.payment_date) - new Date(a.payment_date))
                  .map((payment) => (
                    <tr key={payment.id} className="hover:bg-muted/30">
                      <td className="px-4 py-3 text-sm text-foreground">
                        {formatDate(payment.payment_date)}
                      </td>
                      <td className="px-4 py-3 text-sm font-medium text-foreground">
                        {formatCurrency(parseFloat(payment.amount || 0))}
                      </td>
                      <td className="px-4 py-3 text-sm text-muted-foreground">
                        {payment.description || '—'}
                      </td>
                      <td className="px-4 py-3 text-center">
                        <span
                          className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                            payment.status === 'verified'
                              ? 'bg-success/10 text-success'
                              : payment.status === 'pending'
                                ? 'bg-warning/10 text-warning'
                                : payment.status === 'disputed'
                                  ? 'bg-destructive/10 text-destructive'
                                  : 'bg-muted/50 text-muted-foreground'
                          }`}
                        >
                          {payment.status || 'pending'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-center">
                        {payment.receipt_photo_url ? (
                          <button
                            onClick={() => onViewReceipt(payment)}
                            className="mx-auto inline-flex items-center justify-center rounded p-1 transition-colors hover:bg-primary/10"
                            title="View receipt"
                          >
                            <svg
                              className="h-5 w-5 text-primary"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                              />
                            </svg>
                          </button>
                        ) : (
                          <span className="text-muted-foreground/50">—</span>
                        )}
                      </td>
                      {debtType === 'owed_to_me' && (
                        <td className="px-4 py-3 text-center">
                          <div className="flex items-center justify-center gap-2">
                            {payment.status === 'pending' && (
                              <>
                                <button
                                  onClick={() => handleVerifyPayment(payment)}
                                  disabled={
                                    verifyingPaymentId === payment.id ||
                                    rejectingPaymentId === payment.id ||
                                    deletingPaymentId === payment.id
                                  }
                                  className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-success transition-colors hover:bg-success/10 disabled:cursor-not-allowed disabled:opacity-50"
                                  title="Verify payment"
                                >
                                  {verifyingPaymentId === payment.id ? (
                                    <svg
                                      className="h-4 w-4 animate-spin"
                                      fill="none"
                                      viewBox="0 0 24 24"
                                    >
                                      <circle
                                        className="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        strokeWidth="4"
                                      />
                                      <path
                                        className="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                      />
                                    </svg>
                                  ) : (
                                    <svg
                                      className="h-5 w-5"
                                      fill="none"
                                      stroke="currentColor"
                                      viewBox="0 0 24 24"
                                    >
                                      <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth="2"
                                        d="M5 13l4 4L19 7"
                                      />
                                    </svg>
                                  )}
                                </button>
                                <button
                                  onClick={() => handleRejectPayment(payment)}
                                  disabled={
                                    verifyingPaymentId === payment.id ||
                                    rejectingPaymentId === payment.id ||
                                    deletingPaymentId === payment.id
                                  }
                                  className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-destructive transition-colors hover:bg-destructive/10 disabled:cursor-not-allowed disabled:opacity-50"
                                  title="Reject payment"
                                >
                                  {rejectingPaymentId === payment.id ? (
                                    <svg
                                      className="h-4 w-4 animate-spin"
                                      fill="none"
                                      viewBox="0 0 24 24"
                                    >
                                      <circle
                                        className="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        strokeWidth="4"
                                      />
                                      <path
                                        className="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                      />
                                    </svg>
                                  ) : (
                                    <svg
                                      className="h-5 w-5"
                                      fill="none"
                                      stroke="currentColor"
                                      viewBox="0 0 24 24"
                                    >
                                      <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth="2"
                                        d="M6 18L18 6M6 6l12 12"
                                      />
                                    </svg>
                                  )}
                                </button>
                              </>
                            )}
                            <button
                              onClick={() => handleDeletePayment(payment)}
                              disabled={
                                verifyingPaymentId === payment.id ||
                                rejectingPaymentId === payment.id ||
                                deletingPaymentId === payment.id
                              }
                              className="inline-flex h-8 w-8 items-center justify-center rounded p-1 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive disabled:cursor-not-allowed disabled:opacity-50"
                              title="Delete payment"
                            >
                              {deletingPaymentId === payment.id ? (
                                <svg
                                  className="h-4 w-4 animate-spin"
                                  fill="none"
                                  viewBox="0 0 24 24"
                                >
                                  <circle
                                    className="opacity-25"
                                    cx="12"
                                    cy="12"
                                    r="10"
                                    stroke="currentColor"
                                    strokeWidth="4"
                                  />
                                  <path
                                    className="opacity-75"
                                    fill="currentColor"
                                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                  />
                                </svg>
                              ) : (
                                <svg
                                  className="h-5 w-5"
                                  fill="none"
                                  stroke="currentColor"
                                  viewBox="0 0 24 24"
                                >
                                  <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth="2"
                                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                  />
                                </svg>
                              )}
                            </button>
                          </div>
                        </td>
                      )}
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  )
}
