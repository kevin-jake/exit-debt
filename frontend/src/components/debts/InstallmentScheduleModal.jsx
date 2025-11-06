import { useState, useEffect, useRef } from 'react'
import { useDebtsStore } from '@stores/debtsStore'
import { formatCurrency, formatDate } from '@utils/formatters'

export const InstallmentScheduleModal = ({ debt, payments, onClose }) => {
  const [schedule, setSchedule] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)
  const [displayCount, setDisplayCount] = useState(12)
  const fetchPaymentSchedule = useDebtsStore((state) => state.fetchPaymentSchedule)
  const observerTarget = useRef(null)

  useEffect(() => {
    const loadSchedule = async () => {
      if (!debt?.id) return

      try {
        setIsLoading(true)
        setError(null)
        const scheduleData = await fetchPaymentSchedule(debt.id)

        // Update status for overdue payments based on due date
        const today = new Date()
        today.setHours(0, 0, 0, 0)

        const updatedSchedule = scheduleData.map((item) => {
          if (item.status === 'pending') {
            const dueDate = new Date(item.due_date)
            dueDate.setHours(0, 0, 0, 0)

            if (dueDate < today) {
              return { ...item, status: 'overdue' }
            }
          }
          return item
        })

        setSchedule(updatedSchedule || [])
      } catch (err) {
        console.error('Failed to load payment schedule:', err)
        setError(err.message || 'Failed to load payment schedule')
      } finally {
        setIsLoading(false)
      }
    }

    loadSchedule()
  }, [debt?.id, fetchPaymentSchedule])

  // Infinite scroll observer
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          const maxPayments = debt?.number_of_payments || schedule.length
          if (displayCount < maxPayments) {
            setDisplayCount((prev) => Math.min(prev + 12, maxPayments))
          }
        }
      },
      { threshold: 0.1 }
    )

    if (observerTarget.current) {
      observer.observe(observerTarget.current)
    }

    return () => {
      if (observerTarget.current) {
        observer.unobserve(observerTarget.current)
      }
    }
  }, [displayCount, schedule.length, debt?.number_of_payments])

  if (!debt) return null

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  // Display schedule based on current displayCount
  const maxPayments = debt.number_of_payments || schedule.length
  const displaySchedule = schedule.slice(0, Math.min(displayCount, maxPayments))

  const getStatusBadge = (status) => {
    const badges = {
      paid: 'bg-success/10 text-success',
      pending: 'bg-muted text-muted-foreground',
      overdue: 'bg-destructive/10 text-destructive',
      missed: 'bg-destructive/10 text-destructive',
    }
    const labels = {
      paid: 'Paid',
      pending: 'Pending',
      overdue: 'Overdue',
      missed: 'Missed',
    }
    return (
      <span
        className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${badges[status] || badges.pending}`}
      >
        {labels[status] || status}
      </span>
    )
  }

  return (
    <div
      className="fixed inset-0 z-[60] !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onClick={handleOverlayClick}
    >
      <div className="card my-8 w-full max-w-4xl overflow-hidden">
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-foreground">Payment Schedule</h2>
              <p className="mt-1 text-sm text-muted-foreground">
                {debt.installment_plan && debt.installment_plan !== 'onetime' && (
                  <span className="capitalize">
                    {debt.installment_plan === 'biweekly' ? 'Bi-weekly' : debt.installment_plan}{' '}
                    installments
                  </span>
                )}
              </p>
            </div>
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
          {/* Summary */}
          <div className="my-6 grid grid-cols-1 gap-4 rounded-lg border border-border bg-muted/50 p-4 md:grid-cols-3">
            <div className="text-center">
              <div className="text-sm text-muted-foreground">Total Payments</div>
              <div className="mt-1 text-lg font-bold text-foreground">
                {debt.number_of_payments || schedule.length}
              </div>
            </div>
            <div className="text-center">
              <div className="text-sm text-muted-foreground">Completed</div>
              <div className="mt-1 text-lg font-bold text-success">
                {schedule.filter((s) => s.status === 'paid').length}
              </div>
            </div>
            <div className="text-center">
              <div className="text-sm text-muted-foreground">Remaining</div>
              <div className="mt-1 text-lg font-bold text-primary">
                {(debt.number_of_payments || schedule.length) -
                  schedule.filter((s) => s.status === 'paid').length}
              </div>
            </div>
          </div>
          {isLoading ? (
            <div className="py-8 text-center">
              <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary border-r-transparent"></div>
              <p className="mt-2 text-sm text-muted-foreground">Loading payment schedule...</p>
            </div>
          ) : error ? (
            <div className="py-8 text-center text-destructive">
              <p>{error}</p>
              <button onClick={onClose} className="btn-secondary mt-4">
                Close
              </button>
            </div>
          ) : schedule.length === 0 ? (
            <div className="py-8 text-center text-muted-foreground">
              No payment schedule available for this debt.
            </div>
          ) : (
            <>
              {/* Mobile Card View */}
              <div className="space-y-4 md:hidden">
                {displaySchedule.map((item, index) => {
                  // Calculate remaining balance after this payment
                  const totalAmount = parseFloat(debt.total_amount || 0)

                  // Get the current payment index in the full schedule
                  const currentPaymentIndex = schedule.findIndex(
                    (s) => s.payment_number === item.payment_number
                  )

                  // Calculate cumulative scheduled amount up to and including this payment
                  const cumulativeScheduled = schedule
                    .slice(0, currentPaymentIndex + 1)
                    .reduce((sum, s) => sum + parseFloat(s.scheduled_amount || 0), 0)

                  // Remaining = Total - Cumulative Scheduled Amount
                  const remainingAfter = Math.max(0, totalAmount - cumulativeScheduled)

                  const scheduledAmount = parseFloat(item.scheduled_amount || 0)
                  const paidAmount = parseFloat(item.paid_amount || 0)
                  const stillOwed = parseFloat(item.amount || 0)

                  return (
                    <div
                      key={item.payment_number}
                      className={`rounded-lg border border-border p-4 ${
                        item.status === 'overdue' || item.status === 'missed' ? 'bg-muted/30' : ''
                      }`}
                    >
                      <div className="mb-3 flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-lg font-semibold text-foreground">
                            #{item.payment_number}
                          </span>
                          {getStatusBadge(item.status)}
                        </div>
                        <span className="text-sm text-muted-foreground">
                          {formatDate(item.due_date)}
                        </span>
                      </div>

                      <div className="space-y-2">
                        <div className="flex justify-between">
                          <span className="text-sm text-muted-foreground">Scheduled</span>
                          <span className="text-sm font-medium text-foreground">
                            {formatCurrency(scheduledAmount)}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-sm text-muted-foreground">Paid</span>
                          <span
                            className={`text-sm font-medium ${
                              paidAmount > 0
                                ? paidAmount >= scheduledAmount
                                  ? 'text-success'
                                  : 'text-warning'
                                : 'text-muted-foreground'
                            }`}
                          >
                            {formatCurrency(paidAmount)}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-sm text-muted-foreground">Still Owed</span>
                          <span className="text-sm font-medium">
                            {stillOwed > 0 ? (
                              <span className="text-warning">{formatCurrency(stillOwed)}</span>
                            ) : (
                              <span className="text-muted-foreground">{formatCurrency(0)}</span>
                            )}
                          </span>
                        </div>
                        <div className="flex justify-between border-t border-border pt-2">
                          <span className="text-sm font-medium text-muted-foreground">
                            Balance After
                          </span>
                          <span className="text-sm font-medium text-foreground">
                            {formatCurrency(remainingAfter)}
                          </span>
                        </div>
                      </div>
                    </div>
                  )
                })}
              </div>

              {/* Desktop Table View */}
              <div className="hidden overflow-x-auto md:block">
                <table className="w-full">
                  <thead className="border-b border-border">
                    <tr>
                      <th className="pb-3 text-left text-sm font-medium text-muted-foreground">
                        Payment #
                      </th>
                      <th className="pb-3 text-left text-sm font-medium text-muted-foreground">
                        Due Date
                      </th>
                      <th className="pb-3 text-left text-sm font-medium text-muted-foreground">
                        Status
                      </th>
                      <th className="pb-3 text-right text-sm font-medium text-muted-foreground">
                        Scheduled
                      </th>
                      <th className="pb-3 text-right text-sm font-medium text-muted-foreground">
                        Paid
                      </th>
                      <th className="pb-3 text-right text-sm font-medium text-muted-foreground">
                        Still Owed
                      </th>
                      <th className="pb-3 text-right text-sm font-medium text-muted-foreground">
                        Balance After
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border">
                    {displaySchedule.map((item, index) => {
                      // Calculate remaining balance after this payment
                      const totalAmount = parseFloat(debt.total_amount || 0)

                      // Get the current payment index in the full schedule
                      const currentPaymentIndex = schedule.findIndex(
                        (s) => s.payment_number === item.payment_number
                      )

                      // Calculate cumulative scheduled amount up to and including this payment
                      const cumulativeScheduled = schedule
                        .slice(0, currentPaymentIndex + 1)
                        .reduce((sum, s) => sum + parseFloat(s.scheduled_amount || 0), 0)

                      // Remaining = Total - Cumulative Scheduled Amount
                      const remainingAfter = Math.max(0, totalAmount - cumulativeScheduled)

                      const scheduledAmount = parseFloat(item.scheduled_amount || 0)
                      const paidAmount = parseFloat(item.paid_amount || 0)
                      const stillOwed = parseFloat(item.amount || 0)

                      return (
                        <tr
                          key={item.payment_number}
                          className={`${
                            item.status === 'overdue' || item.status === 'missed'
                              ? 'bg-muted/30'
                              : ''
                          }`}
                        >
                          <td className="py-3 text-sm text-foreground">
                            <span className="font-medium">#{item.payment_number}</span>
                          </td>
                          <td className="py-3 text-sm text-foreground">
                            {formatDate(item.due_date)}
                          </td>
                          <td className="py-3 text-sm">{getStatusBadge(item.status)}</td>
                          <td className="py-3 text-right text-sm font-medium text-foreground">
                            {formatCurrency(scheduledAmount)}
                          </td>
                          <td className="py-3 text-right text-sm font-medium text-foreground">
                            <span
                              className={`${
                                paidAmount > 0
                                  ? paidAmount >= scheduledAmount
                                    ? 'text-success'
                                    : 'text-warning'
                                  : 'text-muted-foreground'
                              }`}
                            >
                              {formatCurrency(paidAmount)}
                            </span>
                          </td>
                          <td className="py-3 text-right text-sm font-medium text-foreground">
                            {stillOwed > 0 ? (
                              <span className="text-warning">{formatCurrency(stillOwed)}</span>
                            ) : (
                              <span className="text-muted-foreground">{formatCurrency(0)}</span>
                            )}
                          </td>
                          <td className="py-3 text-right text-sm text-muted-foreground">
                            {formatCurrency(remainingAfter)}
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>

              {/* Infinite scroll trigger */}
              {displayCount < maxPayments && (
                <div ref={observerTarget} className="py-4 text-center">
                  <div className="inline-block h-6 w-6 animate-spin rounded-full border-2 border-solid border-primary border-r-transparent"></div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
