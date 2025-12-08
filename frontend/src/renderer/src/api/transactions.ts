import { apiClient, APIResponse, PaginatedResponse } from './client'
import { Transaction, TransactionRequest } from '../types'

export async function getTransactions(
  page: number = 1,
  limit: number = 10
): Promise<PaginatedResponse<Transaction[]>> {
  const response = await apiClient.get('/transactions', {
    params: { page, limit }
  })
  return response.data
}

export async function createTransaction(
  data: TransactionRequest
): Promise<APIResponse<Transaction>> {
  const response = await apiClient.post('/transactions', data)
  return response.data
}

export async function getTransaction(id: string): Promise<APIResponse<Transaction>> {
  const response = await apiClient.get(`/transactions/${id}`)
  return response.data
}
