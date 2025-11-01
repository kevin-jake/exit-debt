import { useState, useEffect } from 'react'
import { usePaymentsStore } from '@stores/paymentsStore'
import { ImageViewerModal } from '@components/common/ImageViewerModal'
import { PaymentHistory } from './PaymentHistory'
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

  useEffect(() => {
    loadPayments(debt.id)
  }, [debt.id])

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
            <PaymentHistory
              debtPayments={debtPayments}
              paymentsLoading={paymentsLoading}
              showAddPayment={showAddPayment}
              setShowAddPayment={setShowAddPayment}
              newPayment={newPayment}
              setNewPayment={setNewPayment}
              receiptFile={receiptFile}
              setReceiptFile={setReceiptFile}
              isSubmitting={isSubmitting}
              onAddPayment={handleAddPayment}
              onFileChange={handleFileChange}
              onViewReceipt={(payment) => setViewingReceipt(payment)}
            />

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
      <ImageViewerModal payment={viewingReceipt} onClose={() => setViewingReceipt(null)} />
    </div>
  )
}
