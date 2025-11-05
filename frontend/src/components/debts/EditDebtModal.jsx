import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useDebtsStore } from '@stores/debtsStore'
import { useContactsStore } from '@stores/contactsStore'
import { useNotificationsStore } from '@stores/notificationsStore'

export const EditDebtModal = ({ debt, onClose, onDebtUpdated }) => {
  const updateDebt = useDebtsStore((state) => state.updateDebt)
  const { contacts, fetchContacts } = useContactsStore()
  const { success, error } = useNotificationsStore()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [mouseDownOnOverlay, setMouseDownOnOverlay] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      debt_type: debt.debt_type || 'i_owe',
      total_amount: debt.total_amount || '',
      description: debt.description || '',
      due_date: debt.due_date ? debt.due_date.split('T')[0] : '',
      notes: debt.notes || '',
    },
  })

  const onSubmit = async (data) => {
    setIsSubmitting(true)
    try {
      // Convert due_date to ISO 8601 datetime format if provided
      const debtData = {
        ...data,
        total_amount: String(data.total_amount),
      }

      if (debtData.due_date && debtData.due_date.trim() !== '') {
        debtData.due_date = new Date(debtData.due_date + 'T12:00:00Z').toISOString()
      } else {
        // Remove empty due_date field
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
      className="fixed inset-0 z-50 !mt-0 flex items-start justify-center overflow-y-auto bg-black/60 p-4"
      onMouseDown={handleOverlayMouseDown}
      onClick={handleOverlayClick}
    >
      <div className="card my-8 w-full max-w-lg overflow-hidden">
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
                <option value="i_owe">I Owe</option>
                <option value="owed_to_me">Owed to Me</option>
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

            {/* Due Date */}
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
