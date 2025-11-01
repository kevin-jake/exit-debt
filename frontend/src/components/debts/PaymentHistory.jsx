import { formatCurrency, formatDate } from '@utils/formatters'

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
}) => {
  const handleCancel = () => {
    setShowAddPayment(false)
    setNewPayment({
      payment_date: new Date().toISOString().split('T')[0],
      amount: '',
      description: '',
    })
    setReceiptFile(null)
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
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

