import { useQuery } from '@tanstack/react-query'
import { getDashboardSummary, getSystemStatus } from '../api/dashboard'

export function useDashboardSummary() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: getDashboardSummary,
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
