import { apiClient, APIResponse } from './client'
import { DashboardSummary, SystemStatus } from '../types'

export interface DashboardParams {
  branchId?: string
  date?: string
}

export async function getDashboardSummary(params?: DashboardParams): Promise<APIResponse<DashboardSummary>> {
  const queryParams: Record<string, string> = {}
  if (params?.branchId) {
    queryParams.branch_id = params.branchId
  }
  if (params?.date) {
    queryParams.date = params.date
  }
  const response = await apiClient.get('/dashboard/summary', { params: queryParams })
  return response.data
}

export async function getSystemStatus(): Promise<APIResponse<SystemStatus>> {
  const response = await apiClient.get('/system/status')
  return response.data
}
