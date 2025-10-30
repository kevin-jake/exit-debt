/**
 * Format a number as currency
 * @param {number} amount - The amount to format
 * @param {string} currency - Currency code (default: 'PHP')
 * @returns {string} Formatted currency string
 */
export const formatCurrency = (amount, currency = 'PHP') => {
  return new Intl.NumberFormat('en-PH', {
    style: 'currency',
    currency: currency,
  }).format(amount)
}

/**
 * Format a date string to relative time (e.g., "2 days ago")
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted relative time
 */
export const formatRelativeTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInSeconds = Math.floor((now - date) / 1000)

  if (diffInSeconds < 60) return 'Just now'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} min ago`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`
  if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} days ago`

  return date.toLocaleDateString()
}

/**
 * Format a date string to a readable date
 * @param {string} dateString - ISO date string
 * @param {object} options - Intl.DateTimeFormat options
 * @returns {string} Formatted date
 */
export const formatDate = (dateString, options = {}) => {
  const date = new Date(dateString)
  const defaultOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    ...options,
  }
  return date.toLocaleDateString('en-US', defaultOptions)
}

/**
 * Format a date string to a short date (MM/DD/YYYY)
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted date
 */
export const formatShortDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

/**
 * Get days until a due date
 * @param {string} dueDateString - ISO date string
 * @returns {number} Number of days (negative if overdue)
 */
export const getDaysUntilDue = (dueDateString) => {
  if (!dueDateString) return null
  const dueDate = new Date(dueDateString)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  dueDate.setHours(0, 0, 0, 0)
  const diffTime = dueDate - today
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

/**
 * Get color class based on days until due
 * @param {number} daysUntilDue - Number of days
 * @returns {string} Tailwind color class
 */
export const getDueDateColor = (daysUntilDue) => {
  if (daysUntilDue === null) return 'text-muted-foreground'
  if (daysUntilDue < 0) return 'text-destructive'
  if (daysUntilDue === 0) return 'text-warning'
  if (daysUntilDue <= 7) return 'text-warning'
  return 'text-success'
}

/**
 * Get initials from a name
 * @param {string} name - Full name
 * @returns {string} Initials
 */
export const getInitials = (name) => {
  if (!name) return '?'
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
}

/**
 * Truncate text to a specified length
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
export const truncateText = (text, maxLength = 50) => {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

/**
 * Parse a string amount to float
 * @param {string|number} amount - Amount to parse
 * @returns {number} Parsed amount
 */
export const parseAmount = (amount) => {
  if (typeof amount === 'number') return amount
  return parseFloat(amount) || 0
}

