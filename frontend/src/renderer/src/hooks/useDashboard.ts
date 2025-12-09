import { useQuery } from '@tanstack/react-query'
import { getDashboardSummary, getSystemStatus, DashboardParams } from '../api/dashboard'

export function useDashboardSummary(params?: DashboardParams) {
  return useQuery({
    queryKey: ['dashboard', params?.branchId, params?.date],
    queryFn: () => getDashboardSummary(params),
    refetchInterval: 30000
  })
}

export function useSystemStatus() {
  return useQuery({
    queryKey: ['system-status'],
    queryFn: getSystemStatus,
    refetchInterval: 10000
  })
}
