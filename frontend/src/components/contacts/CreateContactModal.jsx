import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useContactsStore } from '@stores/contactsStore'
import { useNotificationsStore } from '@stores/notificationsStore'

export const CreateContactModal = ({ onContactCreated, onClose }) => {
  const [isLoading, setIsLoading] = useState(false)
  const createContact = useContactsStore((state) => state.createContact)
  const showSuccess = useNotificationsStore((state) => state.success)
  const showError = useNotificationsStore((state) => state.error)

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm()

  const notes = watch('notes', '')

  useEffect(() => {
    // Prevent body scroll when modal is open
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

  const onSubmit = async (data) => {
    setIsLoading(true)

    try {
      const contactData = {
        name: data.name.trim(),
        email: data.email?.trim() || undefined,
        phone: data.phone?.trim() || undefined,
        notes: data.notes?.trim() || undefined,
      }

      const newContact = await createContact(contactData)
      showSuccess(`Successfully created contact "${newContact.name}"`)
      onContactCreated()
    } catch (error) {
      showError(error.message || 'Failed to create contact. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      onClick={onClose}
    >
      <div
        className="w-full max-w-md overflow-hidden rounded-xl bg-card shadow-medium"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-xl font-semibold text-foreground">Create New Contact</h2>
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

        {/* Form */}
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 p-6">
          {/* Name Field */}
          <div>
            <label htmlFor="contact-name" className="label">
              Name *
            </label>
            <input
              id="contact-name"
              type="text"
              {...register('name', {
                required: 'Name is required',
                minLength: {
                  value: 2,
                  message: 'Name must be at least 2 characters',
                },
              })}
              className={`input ${errors.name ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}`}
              placeholder="Enter contact name"
              disabled={isLoading}
            />
            {errors.name && <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>}
          </div>

          {/* Email Field */}
          <div>
            <label htmlFor="contact-email" className="label">
              Email
            </label>
            <input
              id="contact-email"
              type="email"
              {...register('email', {
                pattern: {
                  value: /\S+@\S+\.\S+/,
                  message: 'Please enter a valid email address',
                },
              })}
              className={`input ${errors.email ? 'border-destructive focus:border-destructive focus:ring-destructive' : ''}`}
              placeholder="Enter email address"
              disabled={isLoading}
            />
            {errors.email && <p className="mt-1 text-sm text-destructive">{errors.email.message}</p>}
          </div>

          {/* Phone Field */}
          <div>
            <label htmlFor="contact-phone" className="label">
              Phone
            </label>
            <input
              id="contact-phone"
              type="tel"
              {...register('phone')}
              className="input"
              placeholder="Enter phone number"
              disabled={isLoading}
            />
          </div>

          {/* Notes Field */}
          <div>
            <label htmlFor="contact-notes" className="label">
              Notes
            </label>
            <textarea
              id="contact-notes"
              {...register('notes', {
                maxLength: {
                  value: 500,
                  message: 'Notes cannot exceed 500 characters',
                },
              })}
              rows="3"
              className="input resize-none"
              placeholder="Additional notes about this contact..."
              disabled={isLoading}
            />
            <div className="mt-1 text-xs text-muted-foreground">{notes.length}/500 characters</div>
            {errors.notes && <p className="mt-1 text-sm text-destructive">{errors.notes.message}</p>}
          </div>

          {/* Action Buttons */}
          <div className="flex items-center justify-end space-x-3 pt-4">
            <button type="button" onClick={onClose} className="btn-secondary" disabled={isLoading}>
              Cancel
            </button>
            <button
              type="submit"
              className="btn-primary disabled:cursor-not-allowed disabled:opacity-50"
              disabled={isLoading}
            >
              {isLoading ? (
                <span className="flex items-center">
                  <svg
                    className="-ml-1 mr-3 h-4 w-4 animate-spin text-primary-foreground"
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
                  Creating...
                </span>
              ) : (
                <span className="flex items-center">
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
                      d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                    />
                  </svg>
                  Create Contact
                </span>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

