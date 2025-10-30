import { useState, useEffect } from 'react'
import { usePaymentsStore } from '@stores/paymentsStore'
import {
  formatCurrency,
  formatDate,
  formatRelativeTime,
  getDaysUntilDue,
  getDueDateColor,
} from '@utils/formatters'

export const DebtDetailsModal = ({ debt, onClose, onEdit, onDelete }) => {
  const {
    fetchPayments,
    createPayment,
    uploadReceipt,
    isLoading: paymentsLoading,
  } = usePaymentsStore()
  const [debtPayments, setDebtPayments] = useState([])
  const [showAddPayment, setShowAddPayment] = useState(false)
  const [newPayment, setNewPayment] = useState({
    payment_date: new Date().toISOString().split('T')[0],
    amount: '',
    description: '',
  })
  const [receiptFile, setReceiptFile] = useState(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [viewingReceipt, setViewingReceipt] = useState(null)
  const [imageError, setImageError] = useState(false)
  const [imageBlob, setImageBlob] = useState(null)
  const [imageLoading, setImageLoading] = useState(false)

  useEffect(() => {
    loadPayments(debt.id)
  }, [debt.id])

  // Fetch receipt image with authorization header
  useEffect(() => {
    let isMounted = true
    let objectUrl = null

    const fetchReceiptImage = async () => {
      if (!viewingReceipt?.receipt_photo_url) return

      setImageLoading(true)
      setImageError(false)
      setImageBlob(null)

      try {
        const token = localStorage.getItem('token')
        const response = await fetch(viewingReceipt.receipt_photo_url, {
          headers: token ? { Authorization: `Bearer ${token}` } : {},
        })

        if (!response.ok) {
          throw new Error('Failed to load image')
        }

        const blob = await response.blob()
        objectUrl = URL.createObjectURL(blob)

        if (isMounted) {
          setImageBlob(objectUrl)
        } else {
          // Component unmounted before fetch completed, cleanup immediately
          URL.revokeObjectURL(objectUrl)
        }
      } catch (error) {
        console.error('Failed to load receipt image:', error)
        if (isMounted) {
          setImageError(true)
        }
      } finally {
        if (isMounted) {
          setImageLoading(false)
        }
      }
    }

    if (viewingReceipt) {
      fetchReceiptImage()
    } else {
      // Reset states when modal is closed
      setImageBlob(null)
      setImageError(false)
      setImageLoading(false)
    }

    // Cleanup blob URL when component unmounts or viewingReceipt changes
    return () => {
      isMounted = false
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl)
      }
    }
  }, [viewingReceipt])

  const loadPayments = async (debtId) => {
    console.log('debtId', debtId)
    try {
      const paymentsData = await fetchPayments(debtId)
      setDebtPayments(paymentsData || [])
    } catch (error) {
      console.error('Failed to load payments:', error)
      setDebtPayments([])
    }
  }

  const handleAddPayment = async (e) => {
    e.preventDefault()
    if (!newPayment.amount || parseFloat(newPayment.amount) <= 0) {
      alert('Please enter a valid amount')
      return
    }

    try {
      setIsSubmitting(true)
      // Convert date string to ISO 8601 datetime format
      const paymentDateTime = new Date(newPayment.payment_date + 'T12:00:00Z').toISOString()

      const payment = await createPayment(debt.id, {
        payment_date: paymentDateTime,
        amount: String(newPayment.amount),
        description: newPayment.description,
      })

      // Upload receipt if file is selected
      if (receiptFile && payment.id) {
        await uploadReceipt(payment.id, receiptFile)
      }

      // Reload payments
      await loadPayments(payment.debt_list_id)

      // Reset form
      setNewPayment({
        payment_date: new Date().toISOString().split('T')[0],
        amount: '',
        description: '',
      })
      setReceiptFile(null)
      setShowAddPayment(false)
    } catch (error) {
      console.error('Failed to add payment:', error)
      alert('Failed to add payment. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleFileChange = (e) => {
    const file = e.target.files?.[0]
    if (file) {
      // Validate file type
      const validTypes = ['image/jpeg', 'image/png', 'image/jpg', 'image/webp']
      if (!validTypes.includes(file.type)) {
        alert('Please select a valid image file (JPEG, PNG, or WebP)')
        return
      }
      // Validate file size (max 5MB)
      if (file.size > 5 * 1024 * 1024) {
        alert('File size must be less than 5MB')
        return
      }
      setReceiptFile(file)
    }
  }

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  const daysUntilDue = getDaysUntilDue(debt.due_date)
  // Calculate total paid and remaining balance
  const totalPaid = debtPayments.reduce((sum, payment) => sum + parseFloat(payment.amount || 0), 0)
  const remainingBalance = parseFloat(debt.total_amount || 0) - totalPaid

  return (
    <div
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onClick={handleOverlayClick}
    >
      <div className="card my-8 w-full max-w-2xl overflow-hidden">
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold text-foreground">Debt Details</h2>
            <button
              onClick={onClose}
              className="text-muted-foreground transition-colors hover:text-foreground"
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
          <div className="space-y-6">
            {/* Debt Type Badge */}
            <div className="flex items-center justify-between">
              <span
                className={`inline-flex rounded-full px-3 py-1 text-sm font-medium ${
                  debt.debt_type === 'i_owe'
                    ? 'bg-destructive/10 text-destructive'
                    : 'bg-success/10 text-success'
                }`}
              >
                {debt.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
              </span>
              <div className="text-sm text-muted-foreground">
                Updated {formatRelativeTime(debt.updated_at)}
              </div>
            </div>

            {/* Amount Section */}
            <div className="rounded-lg border border-border bg-muted/50 p-6">
              <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                <div className="text-center">
                  <div className="mb-2 text-sm text-muted-foreground">Total Amount</div>
                  <div
                    className={`text-2xl font-bold ${
                      debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'
                    }`}
                  >
                    {formatCurrency(parseFloat(debt.total_amount || 0))}
                  </div>
                </div>
                <div className="text-center">
                  <div className="mb-2 text-sm text-muted-foreground">Total Paid</div>
                  <div className="text-2xl font-bold text-primary">{formatCurrency(totalPaid)}</div>
                </div>
                <div className="text-center">
                  <div className="mb-2 text-sm text-muted-foreground">Remaining</div>
                  <div
                    className={`text-2xl font-bold ${
                      remainingBalance > 0
                        ? debt.debt_type === 'i_owe'
                          ? 'text-destructive'
                          : 'text-success'
                        : 'text-muted-foreground'
                    }`}
                  >
                    {formatCurrency(remainingBalance)}
                  </div>
                </div>
              </div>
            </div>

            {/* Contact Information */}
            <div className="rounded-lg border border-border p-4">
              <div className="mb-2 text-sm font-medium text-muted-foreground">Contact</div>
              <div className="flex items-center space-x-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <span className="text-lg font-medium text-primary">
                    {debt.contact?.name
                      ?.split(' ')
                      .map((n) => n[0])
                      .join('')
                      .toUpperCase() || 'U'}
                  </span>
                </div>
                <div>
                  <div className="font-medium text-foreground">
                    {debt.contact?.name || 'Unknown Contact'}
                  </div>
                  {debt.contact?.email && (
                    <div className="text-sm text-muted-foreground">{debt.contact.email}</div>
                  )}
                  {debt.contact?.phone && (
                    <div className="text-sm text-muted-foreground">{debt.contact.phone}</div>
                  )}
                </div>
              </div>
            </div>

            {/* Description */}
            <div className="rounded-lg border border-border p-4">
              <div className="mb-2 text-sm font-medium text-muted-foreground">Description</div>
              <div className="text-foreground">{debt.description || 'No description provided'}</div>
            </div>

            {/* Due Date */}
            {debt.due_date && (
              <div className="rounded-lg border border-border p-4">
                <div className="mb-2 text-sm font-medium text-muted-foreground">Due Date</div>
                <div className="flex items-center justify-between">
                  <div className="text-foreground">{formatDate(debt.due_date)}</div>
                  {daysUntilDue !== null && (
                    <span className={`text-sm font-medium ${getDueDateColor(daysUntilDue)}`}>
                      {daysUntilDue < 0
                        ? `Overdue by ${Math.abs(daysUntilDue)} days`
                        : daysUntilDue === 0
                          ? 'Due today'
                          : `Due in ${daysUntilDue} days`}
                    </span>
                  )}
                </div>
              </div>
            )}

            {/* Notes */}
            {debt.notes && (
              <div className="rounded-lg border border-border p-4">
                <div className="mb-2 text-sm font-medium text-muted-foreground">Notes</div>
                <div className="whitespace-pre-wrap text-foreground">{debt.notes}</div>
              </div>
            )}

            {/* Payment History */}
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
                  <form
                    onSubmit={handleAddPayment}
                    className="mt-3 space-y-3 rounded-lg bg-muted/30 p-3"
                  >
                    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                      <div>
                        <label className="mb-1 block text-xs font-medium text-muted-foreground">
                          Payment Date
                        </label>
                        <input
                          type="date"
                          value={newPayment.payment_date}
                          onChange={(e) =>
                            setNewPayment({ ...newPayment, payment_date: e.target.value })
                          }
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
                        onChange={(e) =>
                          setNewPayment({ ...newPayment, description: e.target.value })
                        }
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
                        onChange={handleFileChange}
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
                        onClick={() => {
                          setShowAddPayment(false)
                          setNewPayment({
                            payment_date: new Date().toISOString().split('T')[0],
                            amount: '',
                            description: '',
                          })
                          setReceiptFile(null)
                        }}
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
                <div className="py-8 text-center text-sm text-muted-foreground">
                  Loading payments...
                </div>
              ) : debtPayments.length === 0 ? (
                <div className="py-8 text-center text-sm text-muted-foreground">
                  No payments recorded yet
                </div>
              ) : (
                <div className="overflow-x-auto">
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
                                  onClick={() => {
                                    setViewingReceipt(payment)
                                    setImageError(false)
                                  }}
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
                          </tr>
                        ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            {/* Timestamps */}
            <div className="grid grid-cols-2 gap-4 rounded-lg border border-border p-4">
              <div>
                <div className="mb-1 text-sm font-medium text-muted-foreground">Created</div>
                <div className="text-sm text-foreground">{formatDate(debt.created_at)}</div>
              </div>
              <div>
                <div className="mb-1 text-sm font-medium text-muted-foreground">Last Updated</div>
                <div className="text-sm text-foreground">{formatDate(debt.updated_at)}</div>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="mt-6 flex space-x-3">
            <button onClick={onEdit} className="btn-secondary flex-1">
              <svg className="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                />
              </svg>
              Edit
            </button>
            <button onClick={onDelete} className="btn-destructive flex-1">
              <svg className="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
              Delete
            </button>
          </div>
        </div>
      </div>

      {/* Image Viewer Modal */}
      {viewingReceipt && (
        <div
          className="fixed inset-0 z-[60] flex items-center justify-center bg-black/80 p-4"
          onClick={() => setViewingReceipt(null)}
        >
          <div
            className="relative max-h-[90vh] max-w-4xl overflow-auto rounded-lg bg-card shadow-xl"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div className="flex items-center justify-between border-b border-border p-4">
              <div>
                <h3 className="text-lg font-semibold text-foreground">Payment Receipt</h3>
                <div className="mt-1 text-sm text-muted-foreground">
                  {formatDate(viewingReceipt.payment_date)} •{' '}
                  {formatCurrency(parseFloat(viewingReceipt.amount || 0))}
                </div>
              </div>
              <button
                onClick={() => setViewingReceipt(null)}
                className="text-muted-foreground transition-colors hover:text-foreground"
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

            {/* Image Content */}
            <div className="p-4">
              {imageLoading ? (
                <div className="flex min-h-[300px] items-center justify-center">
                  <div className="text-center">
                    <svg
                      className="mx-auto h-12 w-12 animate-spin text-primary"
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
                    <div className="mt-3 text-sm text-muted-foreground">Loading receipt...</div>
                  </div>
                </div>
              ) : imageError ? (
                <div className="flex min-h-[300px] flex-col items-center justify-center space-y-4 text-center">
                  <svg
                    className="h-16 w-16 text-destructive"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <div>
                    <div className="text-lg font-medium text-foreground">Failed to load image</div>
                    <div className="mt-1 text-sm text-muted-foreground">
                      The receipt photo could not be loaded. It may have been deleted or is
                      temporarily unavailable.
                    </div>
                  </div>
                  <a
                    href={viewingReceipt.receipt_photo_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="btn-secondary"
                  >
                    Try opening in new tab
                  </a>
                </div>
              ) : imageBlob ? (
                <div className="flex items-center justify-center">
                  <img
                    src={imageBlob}
                    alt="Payment receipt"
                    className="max-h-[70vh] w-auto rounded object-contain"
                    onError={() => setImageError(true)}
                  />
                </div>
              ) : null}
            </div>

            {/* Footer with Actions */}
            {!imageError && !imageLoading && imageBlob && (
              <div className="border-t border-border p-4">
                <div className="flex justify-end space-x-2">
                  <button
                    onClick={() => {
                      if (imageBlob) {
                        window.open(imageBlob, '_blank')
                      }
                    }}
                    className="btn-secondary"
                  >
                    <svg
                      className="mr-2 h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                      />
                    </svg>
                    Open in new tab
                  </button>
                  <a href={imageBlob} download="receipt.jpg" className="btn-primary">
                    <svg
                      className="mr-2 h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                      />
                    </svg>
                    Download
                  </a>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
