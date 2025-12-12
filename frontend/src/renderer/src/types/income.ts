export interface IncomeEntry {
  id: string
  branch_id: string
  branch_name: string
  branch_code: string
  date: string
  omzet: number // Decimal with 2 decimal places
  pemasukan_toru: number // Decimal with 2 decimal places
  pemasukan_cash: number // Decimal with 2 decimal places (auto-calculated)
  qris_bca: number // Decimal with 2 decimal places
  qris_bni: number // Decimal with 2 decimal places
  qris_bri: number // Decimal with 2 decimal places
  transfer_bca: number // Decimal with 2 decimal places
  transfer_bni: number // Decimal with 2 decimal places
  transfer_bri: number // Decimal with 2 decimal places
  total_payments: number // Decimal with 2 decimal places
  created_at: string
  updated_at: string
}

export interface IncomeEntryRequest {
  branch_id: string
  date: string
  omzet: number
  pemasukan_toru: number
  qris_bca: number
  qris_bni: number
  qris_bri: number
  transfer_bca: number
  transfer_bni: number
  transfer_bri: number
}

export interface ExpenseEntry {
  id: string
  branch_id: string
  branch_name: string
  branch_code: string
  date: string
  omzet: number
  pengeluaran_toru: number
  pengeluaran_cash: number
  qris_bca: number
  qris_bni: number
  qris_bri: number
  transfer_bca: number
  transfer_bni: number
  transfer_bri: number
  total_payments: number
  created_at: string
  updated_at: string
}

export interface ExpenseEntryRequest {
  branch_id: string
  date: string
  omzet: number
  pengeluaran_toru: number
  qris_bca: number
  qris_bni: number
  qris_bri: number
  transfer_bca: number
  transfer_bni: number
  transfer_bri: number
}

export interface IncomeRowData {
  branch_id: string
  branch_name: string
  omzet: string
  pemasukan_toru: string
  pemasukan_cash: number
  qris_bca: string
  qris_bni: string
  qris_bri: string
  transfer_bca: string
  transfer_bni: string
  transfer_bri: string
  difference?: number // positive = sisa, negative = kelebihan
  status?: 'over' | 'under'
  warning?: string
}
