import { useEffect } from 'react'
import { formatCurrency } from '@utils/formatters'

/**
 * Calculate number of installments based on due date and payment frequency
 */
const calculateInstallmentsFromDate = (dueDate, frequency) => {
  if (!dueDate) return 0

  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const due = new Date(dueDate)
  due.setHours(0, 0, 0, 0)

  const diffTime = due - today
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays <= 0) return 0

  let installments = 0
  switch (frequency) {
    case 'weekly':
      installments = Math.floor(diffDays / 7)
      break
    case 'biweekly':
      installments = Math.floor(diffDays / 14)
      break
    case 'monthly':
      const monthsDiff =
        (due.getFullYear() - today.getFullYear()) * 12 + (due.getMonth() - today.getMonth())
      installments = monthsDiff
      break
    case 'quarterly':
      const quartersDiff =
        (due.getFullYear() - today.getFullYear()) * 4 +
        Math.floor(due.getMonth() / 3) -
        Math.floor(today.getMonth() / 3)
      installments = quartersDiff
      break
    default:
      installments = 0
  }

  return Math.max(0, installments)
}

/**
 * Calculate due date based on number of installments and payment frequency
 */
const calculateDueDateFromInstallments = (numberOfInstallments, frequency) => {
  if (!numberOfInstallments || numberOfInstallments < 1) return null

  const today = new Date()
  const dueDate = new Date(today)

  const count = parseInt(numberOfInstallments)

  switch (frequency) {
    case 'weekly':
      dueDate.setDate(dueDate.getDate() + count * 7)
      break
    case 'biweekly':
      dueDate.setDate(dueDate.getDate() + count * 14)
      break
    case 'monthly':
      dueDate.setMonth(dueDate.getMonth() + count)
      break
    case 'quarterly':
      dueDate.setMonth(dueDate.getMonth() + count * 3)
      break
    default:
      return null
  }

  return dueDate.toISOString().split('T')[0]
}

/**
 * PaymentFields component - handles payment type, installment details, and due date
 * @param {Object} props
 * @param {Function} props.register - react-hook-form register function
 * @param {Function} props.watch - react-hook-form watch function
 * @param {Function} props.setValue - react-hook-form setValue function
 * @param {Object} props.errors - react-hook-form errors object
 * @param {boolean} props.isSubmitting - form submission state
 */
export const PaymentFields = ({ register, watch, errors, isSubmitting, setValue }) => {
  const paymentType = watch('payment_type')
  const totalAmount = watch('total_amount')
  const numberOfInstallments = watch('number_of_payments')
  const dueDate = watch('due_date')
  const paymentFrequency = watch('installment_plan')
  const calculationMethod = watch('installment_calculation_method')

  // Calculate installments or due date based on method
  useEffect(() => {
    if (paymentType !== 'installment') return

    if (calculationMethod === 'by_date' && dueDate && paymentFrequency) {
      const calculated = calculateInstallmentsFromDate(dueDate, paymentFrequency)

      // If calculated installments is 0, automatically convert to one-time payment
      if (calculated === 0) {
        setValue('payment_type', 'one_time')
        setValue('number_of_payments', '')
        setValue('due_date', '') // Clear due date to allow switching back
        return
      }

      if (calculated !== parseInt(numberOfInstallments)) {
        setValue('number_of_payments', calculated > 0 ? calculated.toString() : '')
      }
    } else if (calculationMethod === 'by_count' && numberOfInstallments && paymentFrequency) {
      const calculated = calculateDueDateFromInstallments(numberOfInstallments, paymentFrequency)
      if (calculated && calculated !== dueDate) {
        setValue('due_date', calculated)
      }
    }
  }, [paymentType, calculationMethod, dueDate, numberOfInstallments, paymentFrequency, setValue])

  const showInstallmentFields = paymentType === 'installment'
  const calculatedInstallments = parseInt(numberOfInstallments) || 0
  const shouldBeOneTime = showInstallmentFields && calculatedInstallments <= 1

  return (
    <>
      {/* Payment Type */}
      <div>
        <label htmlFor="payment_type" className="mb-2 block text-sm font-medium text-foreground">
          Payment Type <span className="text-destructive">*</span>
        </label>
        <select
          id="payment_type"
          {...register('payment_type', { required: 'Payment type is required' })}
          className="input"
          disabled={isSubmitting}
        >
          <option value="one_time">One-Time Payment</option>
          <option value="installment">Installment Plan</option>
        </select>
        {errors.payment_type && (
          <p className="mt-1 text-sm text-destructive">{errors.payment_type.message}</p>
        )}
      </div>

      {/* Show warning if calculated installments is 1 */}
      {shouldBeOneTime && calculatedInstallments === 1 && (
        <div className="rounded-md border border-orange-500/50 bg-orange-500/10 p-3">
          <p className="text-sm text-orange-600 dark:text-orange-400">
            ⚠️ Based on your selected due date and payment frequency, this debt will only have 1
            installment. Consider selecting "One-Time Payment" instead or adjusting the due
            date/payment frequency for multiple installments.
          </p>
        </div>
      )}

      {/* Installment Fields - Only show when payment_type is 'installment' */}
      {showInstallmentFields && (
        <>
          {/* How would you like to pay */}
          <div>
            <label
              htmlFor="installment_calculation_method"
              className="mb-2 block text-sm font-medium text-foreground"
            >
              How would you like to pay? <span className="text-destructive">*</span>
            </label>
            <select
              id="installment_calculation_method"
              {...register('installment_calculation_method', {
                required: 'Please select a calculation method',
              })}
              className="input"
              disabled={isSubmitting}
            >
              <option value="by_count">Number of Installments</option>
              <option value="by_date">Due Date</option>
            </select>
            {errors.installment_calculation_method && (
              <p className="mt-1 text-sm text-destructive">
                {errors.installment_calculation_method.message}
              </p>
            )}
          </div>

          {/* Payment Frequency */}
          <div>
            <label
              htmlFor="installment_plan"
              className="mb-2 block text-sm font-medium text-foreground"
            >
              Payment Frequency <span className="text-destructive">*</span>
            </label>
            <select
              id="installment_plan"
              {...register('installment_plan', {
                required: paymentType === 'installment' && 'Payment frequency is required',
              })}
              className="input"
              disabled={isSubmitting}
            >
              <option value="">Select frequency</option>
              <option value="weekly">Weekly</option>
              <option value="biweekly">Bi-weekly</option>
              <option value="monthly">Monthly</option>
              <option value="quarterly">Quarterly</option>
            </select>
            {errors.installment_plan && (
              <p className="mt-1 text-sm text-destructive">{errors.installment_plan.message}</p>
            )}
          </div>

          {/* Conditional Field: Number of Installments (when by_count) */}
          {calculationMethod === 'by_count' && (
            <div>
              <label
                htmlFor="number_of_payments"
                className="mb-2 block text-sm font-medium text-foreground"
              >
                Number of Installments <span className="text-destructive">*</span>
              </label>
              <input
                id="number_of_payments"
                type="number"
                min="2"
                {...register('number_of_payments', {
                  required:
                    calculationMethod === 'by_count' && 'Number of installments is required',
                  min: { value: 2, message: 'Must have at least 2 installments' },
                })}
                className="input"
                disabled={isSubmitting}
                placeholder="e.g., 12"
              />
              {errors.number_of_payments && (
                <p className="mt-1 text-sm text-destructive">{errors.number_of_payments.message}</p>
              )}
            </div>
          )}

          {/* Conditional Field: Due Date (when by_date) */}
          {calculationMethod === 'by_date' && (
            <div>
              <label htmlFor="due_date" className="mb-2 block text-sm font-medium text-foreground">
                Due Date <span className="text-destructive">*</span>
              </label>
              <input
                id="due_date"
                type="date"
                {...register('due_date', {
                  required: calculationMethod === 'by_date' && 'Due date is required',
                })}
                className="input"
                disabled={isSubmitting}
              />
              {errors.due_date && (
                <p className="mt-1 text-sm text-destructive">{errors.due_date.message}</p>
              )}
            </div>
          )}

          {/* Display Calculated Value */}
          {calculationMethod === 'by_date' && numberOfInstallments && (
            <div className="rounded-md bg-muted/50 p-3">
              <p className="text-sm text-muted-foreground">
                Calculated Number of Installments: {numberOfInstallments}
              </p>
            </div>
          )}
          {/* Installment Amount (Calculated) */}

          {calculationMethod === 'by_count' && dueDate && (
            <div className="rounded-md bg-muted/50 p-3">
              <p className="text-sm text-muted-foreground">
                Calculated Final Due Date:{' '}
                {new Date(dueDate).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </p>
            </div>
          )}
          {totalAmount && numberOfInstallments > 0 && (
            <div className="rounded-md bg-muted/50 p-3">
              <p className="text-sm text-muted-foreground">
                Calculated Installment Amount:{' '}
                {formatCurrency(parseFloat(totalAmount) / parseInt(numberOfInstallments))} per
                installment
              </p>
            </div>
          )}
        </>
      )}

      {/* Due Date - Only show for one-time payments */}
      {paymentType === 'one_time' && (
        <div>
          <label htmlFor="due_date" className="mb-2 block text-sm font-medium text-foreground">
            Due Date
          </label>
          <input
            id="due_date"
            type="date"
            {...register('due_date')}
            className="input"
            disabled={isSubmitting}
          />
          {errors.due_date && (
            <p className="mt-1 text-sm text-destructive">{errors.due_date.message}</p>
          )}
        </div>
      )}
    </>
  )
}
