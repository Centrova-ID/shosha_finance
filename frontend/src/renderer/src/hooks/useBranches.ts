import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { branchAPI, BranchRequest } from '@/api/branches'

export function useBranches() {
  return useQuery({
    queryKey: ['branches'],
    queryFn: () => branchAPI.getAll()
  })
}

export function useActiveBranches() {
  return useQuery({
    queryKey: ['branches', 'active'],
    queryFn: () => branchAPI.getActive()
  })
}

export function useCreateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: BranchRequest) => branchAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['branches'] })
    }
  })
}

export function useUpdateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: BranchRequest }) =>
      branchAPI.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['branches'] })
    }
  })
}

export function useDeleteBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => branchAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['branches'] })
    }
  })
}
