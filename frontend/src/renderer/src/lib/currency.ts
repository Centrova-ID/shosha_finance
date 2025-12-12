/**
 * Parse numeric string to number, removing all non-digits
 */
export const parseNumber = (value: string): number => {
  const parsed = parseInt(value.replace(/\D/g, ''), 10)
  return isNaN(parsed) ? 0 : parsed
}

/**
 * Format number as Indonesian currency (Rp)
 */
export const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

/**
 * Format number with thousand separators for input display
 */
export const formatNumberInput = (value: string): string => {
  const num = parseNumber(value)
  return new Intl.NumberFormat('id-ID').format(num)
}

/**
 * Handle currency input with proper formatting while preserving editing
 */
export const handleCurrencyInput = (value: string): string => {
  // Remove all non-digit characters
  const digitsOnly = value.replace(/\D/g, '')
  if (!digitsOnly) return ''
  
  // Format with thousand separators
  return new Intl.NumberFormat('id-ID').format(parseInt(digitsOnly, 10))
}
