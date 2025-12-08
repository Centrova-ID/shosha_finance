import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getTransactions, createTransaction } from '../api/transactions'
import { TransactionRequest } from '../types'

export function useTransactions(page: number = 1, limit: number = 10) {
  return useQuery({
    queryKey: ['transactions', page, limit],
    queryFn: () => getTransactions(page, limit)
  })
}

export function useCreateTransaction() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: TransactionRequest) => createTransaction(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
      queryClient.invalidateQueries({ queryKey: ['system-status'] })
    }
  })
}
