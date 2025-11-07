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

/**
 * Get debt status based on due date
 * @param {string} dueDateString - ISO date string
 * @returns {object} Status object with label and color
 */
export const getDebtStatus = (dueDateString) => {
  if (!dueDateString) {
    return { label: 'Active', color: 'text-muted-foreground', bgColor: 'bg-muted/50' }
  }

  const daysUntilDue = getDaysUntilDue(dueDateString)

  if (daysUntilDue < 0) {
    return { label: 'Overdue', color: 'text-destructive', bgColor: 'bg-destructive/10' }
  }

  if (daysUntilDue === 0) {
    return { label: 'Due Today', color: 'text-warning', bgColor: 'bg-warning/10' }
  }

  if (daysUntilDue <= 7) {
    return { label: 'Due Soon', color: 'text-warning', bgColor: 'bg-warning/10' }
  }

  return { label: 'Active', color: 'text-blue-500', bgColor: 'bg-blue-500/10' }
}

/**
 * Convert a date string to ISO 8601 format
 * @param {string} dateStr - Date string in YYYY-MM-DD or ISO format
 * @returns {string|null} ISO formatted date string or null if invalid
 */
export const convertToISO = (dateStr) => {
  if (!dateStr || dateStr === '') return null
  try {
    // Handle both YYYY-MM-DD format and ISO format
    const date = new Date(dateStr)
    // Check if date is valid
    if (isNaN(date.getTime())) return null
    // Set to noon UTC to avoid timezone issues
    date.setUTCHours(12, 0, 0, 0)
    return date.toISOString()
  } catch (e) {
    console.error('Error converting date:', e)
    return null
  }
}
