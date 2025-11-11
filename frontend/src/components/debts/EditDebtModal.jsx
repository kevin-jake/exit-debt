import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useDebtsStore } from '@stores/debtsStore'
import { useContactsStore } from '@stores/contactsStore'
import { useNotificationsStore } from '@stores/notificationsStore'
import { convertToISO } from '@utils/formatters'
import { PaymentFields } from './PaymentFields'

export const EditDebtModal = ({ debt, onClose, onDebtUpdated }) => {
  const updateDebt = useDebtsStore((state) => state.updateDebt)
  const { success, error } = useNotificationsStore()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [mouseDownOnOverlay, setMouseDownOnOverlay] = useState(false)

  // Determine payment type based on debt data
  const getPaymentType = (debtData) => {
    // If payment_type is explicitly set, use it
    if (debtData.payment_type) {
      return debtData.payment_type
    }
    // Otherwise, check if debt has installment data (with actual values, not empty strings)
    const hasInstallments =
      (debtData.number_of_payments && debtData.number_of_payments > 1) ||
      (debtData.installment_plan != 'onetime' && debtData.installment_plan.trim() !== '')

    return hasInstallments ? debtData.payment_type : 'onetime'
  }

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
    reset,
  } = useForm({
    defaultValues: {
      debt_type: debt.debt_type || 'to_pay',
      payment_type: getPaymentType(debt),
      total_amount: debt.total_amount || '',
      description: debt.description || '',
      due_date: debt.due_date ? debt.due_date.split('T')[0] : '',
      notes: debt.notes || '',
      number_of_payments: debt.number_of_payments || '',
      installment_amount: debt.installment_amount || '',
      installment_plan: debt.installment_plan || 'monthly',
      next_payment_date: debt.next_payment_date ? debt.next_payment_date.split('T')[0] : '',
      installment_calculation_method: 'by_count',
    },
  })

  // Reset form when debt changes
  useEffect(() => {
    const paymentType = getPaymentType(debt)
    console.log('Resetting form with payment_type:', paymentType, 'for debt:', debt.id)
    console.log('Debt data:', {
      number_of_payments: debt.number_of_payments,
      installment_plan: debt.installment_plan,
      payment_type: debt.payment_type,
    })

    reset({
      debt_type: debt.debt_type || 'to_pay',
      payment_type: paymentType,
      total_amount: debt.total_amount || '',
      description: debt.description || '',
      due_date: debt.due_date ? debt.due_date.split('T')[0] : '',
      notes: debt.notes || '',
      number_of_payments: debt.number_of_payments || '',
      installment_amount: debt.installment_amount || '',
      installment_plan: debt.installment_plan || 'monthly',
      next_payment_date: debt.next_payment_date ? debt.next_payment_date.split('T')[0] : '',
      installment_calculation_method: 'by_count',
    })
  }, [debt, reset])

  // Watch payment type and reset payment frequency when it changes
  const paymentType = watch('payment_type')
  useEffect(() => {
    // When payment type changes to installment, reset frequency to monthly
    if (paymentType === 'installment') {
      setValue('installment_plan', 'monthly')
    }
  }, [paymentType, setValue])

  const onSubmit = async (data) => {
    setIsSubmitting(true)
    try {
      // Convert due_date to ISO 8601 datetime format if provided
      const debtData = {
        ...data,
        total_amount: String(data.total_amount),
        installment_plan: data.payment_type === 'onetime' ? 'onetime' : data.installment_plan,
        number_of_payments:
          data.payment_type === 'installment' ? parseInt(data.number_of_payments) : 1,
      }

      const convertedDueDate = convertToISO(debtData.due_date)
      if (convertedDueDate) {
        debtData.due_date = convertedDueDate
      } else {
        delete debtData.due_date
      }

      await updateDebt(debt.id, debtData)
      success('Debt updated successfully')
      onDebtUpdated()
    } catch (err) {
      error(err.message || 'Failed to update debt')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleOverlayMouseDown = (e) => {
    // Track if mousedown happened on the overlay
    if (e.target === e.currentTarget && !isSubmitting) {
      setMouseDownOnOverlay(true)
    }
  }

  const handleOverlayClick = (e) => {
    // Only close if both mousedown and click happened on the overlay
    if (e.target === e.currentTarget && !isSubmitting && mouseDownOnOverlay) {
      onClose()
    }
    // Reset the flag
    setMouseDownOnOverlay(false)
  }

  return (
    <div
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4 pb-96"
      onMouseDown={handleOverlayMouseDown}
      onClick={handleOverlayClick}
    >
      <div className="card my-8 w-full max-w-lg overflow-visible">
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-semibold text-foreground">Edit Debt</h2>
            <button
              onClick={onClose}
              disabled={isSubmitting}
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

        <form onSubmit={handleSubmit(onSubmit)} className="p-6">
          <div className="space-y-4">
            {/* Debt Type */}
            <div>
              <label htmlFor="debt_type" className="mb-2 block text-sm font-medium text-foreground">
                Debt Type <span className="text-destructive">*</span>
              </label>
              <select
                id="debt_type"
                {...register('debt_type', { required: 'Debt type is required' })}
                className="input"
                disabled={isSubmitting}
              >
                <option value="to_pay">To Pay</option>
                <option value="to_receive">To Receive</option>
              </select>
              {errors.debt_type && (
                <p className="mt-1 text-sm text-destructive">{errors.debt_type.message}</p>
              )}
            </div>

            {/* Total Amount */}
            <div>
              <label
                htmlFor="total_amount"
                className="mb-2 block text-sm font-medium text-foreground"
              >
                Total Amount <span className="text-destructive">*</span>
              </label>
              <input
                id="total_amount"
                type="number"
                step="0.01"
                min="0"
                {...register('total_amount', {
                  required: 'Total amount is required',
                  min: { value: 0.01, message: 'Amount must be greater than 0' },
                })}
                className="input"
                disabled={isSubmitting}
                placeholder="0.00"
              />
              {errors.total_amount && (
                <p className="mt-1 text-sm text-destructive">{errors.total_amount.message}</p>
              )}
            </div>

            {/* Payment Fields Component */}
            <PaymentFields
              register={register}
              watch={watch}
              setValue={setValue}
              errors={errors}
              isSubmitting={isSubmitting}
            />

            {/* Description */}
            <div>
              <label
                htmlFor="description"
                className="mb-2 block text-sm font-medium text-foreground"
              >
                Description
              </label>
              <input
                id="description"
                type="text"
                {...register('description', {
                  maxLength: {
                    value: 255,
                    message: 'Description must be less than 255 characters',
                  },
                })}
                className="input"
                disabled={isSubmitting}
                placeholder="What is this debt for?"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-destructive">{errors.description.message}</p>
              )}
            </div>

            {/* Notes */}
            <div>
              <label htmlFor="notes" className="mb-2 block text-sm font-medium text-foreground">
                Notes
              </label>
              <textarea
                id="notes"
                rows={3}
                {...register('notes')}
                className="input resize-none"
                disabled={isSubmitting}
                placeholder="Additional notes or details"
              />
              {errors.notes && (
                <p className="mt-1 text-sm text-destructive">{errors.notes.message}</p>
              )}
            </div>
          </div>

          <div className="mt-6 flex space-x-3">
            <button
              type="button"
              onClick={onClose}
              disabled={isSubmitting}
              className="btn-secondary flex-1"
            >
              Cancel
            </button>
            <button type="submit" disabled={isSubmitting} className="btn-primary flex-1">
              {isSubmitting ? (
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
                  Updating...
                </span>
              ) : (
                'Update Debt'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
