import { apiClient, APIResponse } from './client'
import { DashboardSummary, SystemStatus } from '../types'

export async function getDashboardSummary(): Promise<APIResponse<DashboardSummary>> {
  const response = await apiClient.get('/dashboard/summary')
  return response.data
}

export async function getSystemStatus(): Promise<APIResponse<SystemStatus>> {
  const response = await apiClient.get('/system/status')
  return response.data
}
