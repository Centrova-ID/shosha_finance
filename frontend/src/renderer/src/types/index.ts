export type TransactionType = 'IN' | 'OUT'

export interface Branch {
  id: string
  code: string
  name: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Transaction {
  id: string
  branch_id: string
  type: TransactionType
  category: string
  amount: number
  description: string
  created_at: string
  branch?: Branch
}

export interface TransactionRequest {
  branch_id: string
  type: TransactionType
  category: string
  amount: number
  description?: string
}

export interface DashboardSummary {
  total_in: number
  total_out: number
  balance: number
  count_in: number
  count_out: number
  unsync_count: number
}

export interface SystemStatus {
  status: 'online' | 'offline'
  unsynced_count: number
  timestamp: string
}
