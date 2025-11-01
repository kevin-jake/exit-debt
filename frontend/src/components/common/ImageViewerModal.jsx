import { useState, useEffect } from 'react'
import { formatCurrency, formatDate } from '@utils/formatters'

export const ImageViewerModal = ({ payment, onClose }) => {
  const [imageError, setImageError] = useState(false)
  const [imageBlob, setImageBlob] = useState(null)
  const [imageLoading, setImageLoading] = useState(false)

  // Fetch receipt image with authorization header
  useEffect(() => {
    let isMounted = true
    let objectUrl = null

    const fetchReceiptImage = async () => {
      if (!payment?.receipt_photo_url) return

      setImageLoading(true)
      setImageError(false)
      setImageBlob(null)

      try {
        const token = localStorage.getItem('token')
        const response = await fetch(payment.receipt_photo_url, {
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

    if (payment) {
      fetchReceiptImage()
    }

    // Cleanup blob URL when component unmounts or payment changes
    return () => {
      isMounted = false
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl)
      }
    }
  }, [payment])

  if (!payment) return null

  return (
    <div
      className="fixed inset-0 z-[60] flex items-center justify-center bg-black/80 p-4"
      onClick={onClose}
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
              {formatDate(payment.payment_date)} • {formatCurrency(parseFloat(payment.amount || 0))}
            </div>
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
                  The receipt photo could not be loaded. It may have been deleted or is temporarily
                  unavailable.
                </div>
              </div>
              <a
                href={payment.receipt_photo_url}
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
  )
}

