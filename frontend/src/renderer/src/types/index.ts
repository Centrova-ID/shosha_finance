export type TransactionType = 'IN' | 'OUT'

export interface Transaction {
  id: string
  branch_id: string
  type: TransactionType
  category: string
  amount: number
  description: string
  created_at: string
  is_synced: boolean
  synced_at: string | null
}

export interface TransactionRequest {
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
  unsynced_count: number
  status: 'online' | 'offline'
  timestamp: string
}

export interface Branch {
  id: string
  code: string
  name: string
  api_key: string
  created_at: string
  updated_at: string
}
