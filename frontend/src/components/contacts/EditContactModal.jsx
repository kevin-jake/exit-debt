import { useState, useEffect } from 'react'
import { useForm, Controller } from 'react-hook-form'
import PhoneInput, { isValidPhoneNumber } from 'react-phone-number-input'
import 'react-phone-number-input/style.css'
import { useContactsStore } from '@stores/contactsStore'
import { useNotificationsStore } from '@stores/notificationsStore'

export const EditContactModal = ({ contact, onContactUpdated, onClose }) => {
  const [isLoading, setIsLoading] = useState(false)
  const [mouseDownOnOverlay, setMouseDownOnOverlay] = useState(false)
  const updateContact = useContactsStore((state) => state.updateContact)
  const showSuccess = useNotificationsStore((state) => state.success)
  const showError = useNotificationsStore((state) => state.error)

  const {
    register,
    handleSubmit,
    watch,
    control,
    formState: { errors },
  } = useForm({
    defaultValues: {
      name: contact.name,
      email: contact.email || '',
      phone: contact.phone || '',
      notes: contact.notes || '',
    },
  })

  const notes = watch('notes', contact.notes || '')

  useEffect(() => {
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

      await updateContact(contact.id, contactData)
      showSuccess(`Successfully updated contact "${data.name}"`)
      onContactUpdated()
    } catch (error) {
      showError(error.message || 'Failed to update contact. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleOverlayMouseDown = (e) => {
    // Track if mousedown happened on the overlay
    if (e.target === e.currentTarget && !isLoading) {
      setMouseDownOnOverlay(true)
    }
  }

  const handleOverlayClick = (e) => {
    // Only close if both mousedown and click happened on the overlay
    if (e.target === e.currentTarget && !isLoading && mouseDownOnOverlay) {
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
      <div
        className="my-8 w-full max-w-md overflow-hidden rounded-xl bg-card shadow-medium"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-xl font-semibold text-foreground">Edit Contact</h2>
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

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 p-6">
          <div>
            <label htmlFor="edit-name" className="label">
              Name *
            </label>
            <input
              id="edit-name"
              type="text"
              {...register('name', {
                required: 'Name is required',
                minLength: {
                  value: 2,
                  message: 'Name must be at least 2 characters',
                },
              })}
              className={`input ${errors.name ? 'border-destructive' : ''}`}
              disabled={isLoading}
            />
            {errors.name && <p className="mt-1 text-sm text-destructive">{errors.name.message}</p>}
          </div>

          <div>
            <label htmlFor="edit-email" className="label">
              Email
            </label>
            <input
              id="edit-email"
              type="email"
              {...register('email', {
                pattern: {
                  value: /\S+@\S+\.\S+/,
                  message: 'Please enter a valid email address',
                },
              })}
              className={`input ${errors.email ? 'border-destructive' : ''}`}
              disabled={isLoading}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="edit-phone" className="label">
              Phone
            </label>
            <Controller
              name="phone"
              control={control}
              rules={{
                validate: (value) => {
                  if (!value) return true // Optional field
                  return isValidPhoneNumber(value) || 'Please enter a valid phone number'
                },
              }}
              render={({ field: { onChange, value } }) => (
                <PhoneInput
                  id="edit-phone"
                  value={value}
                  onChange={onChange}
                  defaultCountry="PH"
                  international
                  className="input"
                  placeholder="Enter phone number"
                  disabled={isLoading}
                />
              )}
            />
            {errors.phone && (
              <p className="mt-1 text-sm text-destructive">{errors.phone.message}</p>
            )}
          </div>

          <div>
            <label htmlFor="edit-notes" className="label">
              Notes
            </label>
            <textarea
              id="edit-notes"
              {...register('notes', {
                maxLength: {
                  value: 500,
                  message: 'Notes cannot exceed 500 characters',
                },
              })}
              rows="3"
              className="input resize-none"
              disabled={isLoading}
            />
            <div className="mt-1 text-xs text-muted-foreground">{notes.length}/500 characters</div>
            {errors.notes && (
              <p className="mt-1 text-sm text-destructive">{errors.notes.message}</p>
            )}
          </div>

          <div className="flex items-center justify-end space-x-3 pt-4">
            <button type="button" onClick={onClose} className="btn-secondary" disabled={isLoading}>
              Cancel
            </button>
            <button
              type="submit"
              className="btn-primary disabled:cursor-not-allowed disabled:opacity-50"
              disabled={isLoading}
            >
              {isLoading ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
