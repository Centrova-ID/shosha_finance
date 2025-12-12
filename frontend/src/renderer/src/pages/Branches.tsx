import { useState } from 'react'
import { useBranches, useCreateBranch, useUpdateBranch, useDeleteBranch } from '@/hooks/useBranches'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle
} from '@/components/ui/sheet'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from '@/components/ui/alert-dialog'
import { toast } from '@/hooks/use-toast'
import { Branch, BranchRequest } from '@/api/branches'
import { Plus, Pencil, Trash2, RefreshCw, Building2 } from 'lucide-react'

export default function Branches() {
  const { data, isLoading } = useBranches()
  const createMutation = useCreateBranch()
  const updateMutation = useUpdateBranch()
  const deleteMutation = useDeleteBranch()

  const [sheetOpen, setSheetOpen] = useState(false)
  const [editingBranch, setEditingBranch] = useState<Branch | null>(null)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [branchToDelete, setBranchToDelete] = useState<Branch | null>(null)
  const [formData, setFormData] = useState<BranchRequest>({
    code: '',
    name: '',
    description: ''
  })

  const resetForm = () => {
    setFormData({ code: '', name: '', description: '' })
    setEditingBranch(null)
  }

  const openCreateSheet = () => {
    resetForm()
    setSheetOpen(true)
  }

  const openEditSheet = (branch: Branch) => {
    setEditingBranch(branch)
    setFormData({
      code: branch.code,
      name: branch.name,
      description: branch.description || ''
    })
    setSheetOpen(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.code || !formData.name) {
      toast({
        title: 'Error',
        description: 'Kode dan nama harus diisi',
        variant: 'destructive'
      })
      return
    }

    try {
      if (editingBranch) {
        await updateMutation.mutateAsync({ id: editingBranch.id, data: formData })
        toast({ title: 'Berhasil', description: 'Unit berhasil diupdate' })
      } else {
        await createMutation.mutateAsync(formData)
        toast({ title: 'Berhasil', description: 'Unit berhasil dibuat' })
      }
      setSheetOpen(false)
      resetForm()
    } catch {
      toast({
        title: 'Error',
        description: 'Gagal menyimpan unit',
        variant: 'destructive'
      })
    }
  }

  const handleDelete = (branch: Branch) => {
    setBranchToDelete(branch)
    setConfirmOpen(true)
  }

  const confirmDelete = async () => {
    if (!branchToDelete) return

    try {
      await deleteMutation.mutateAsync(branchToDelete.id)
      toast({ title: 'Berhasil', description: 'Unit berhasil dihapus' })
    } catch {
      toast({
        title: 'Error',
        description: 'Gagal menghapus unit',
        variant: 'destructive'
      })
    } finally {
      setConfirmOpen(false)
      setBranchToDelete(null)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  const branches = data?.data || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">Kelola Unit</h2>
        <Button onClick={openCreateSheet}>
          <Plus className="mr-2 h-4 w-4" />
          Tambah Unit
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {branches.length === 0 ? (
          <Card className="col-span-full">
            <CardContent className="flex flex-col items-center justify-center py-8">
              <Building2 className="h-12 w-12 text-muted-foreground mb-4" />
              <p className="text-muted-foreground">Belum ada unit. Tambahkan unit pertama Anda!</p>
            </CardContent>
          </Card>
        ) : (
          branches.map((branch) => (
            <Card key={branch.id}>
              <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                <div>
                  <CardTitle className="text-lg">{branch.name}</CardTitle>
                  <p className="text-sm text-muted-foreground font-mono">{branch.code}</p>
                </div>
                <div className="flex gap-1">
                  <Button variant="ghost" size="icon" onClick={() => openEditSheet(branch)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(branch)}
                    className="text-destructive hover:text-destructive"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  {branch.description || 'Tidak ada deskripsi'}
                </p>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>{editingBranch ? 'Edit Unit' : 'Tambah Unit Baru'}</SheetTitle>
            <SheetDescription>
              {editingBranch ? 'Update informasi unit' : 'Buat unit baru untuk mengelompokkan transaksi'}
            </SheetDescription>
          </SheetHeader>
          <form onSubmit={handleSubmit} className="space-y-6 mt-6">
            <div className="space-y-2">
              <Label htmlFor="code">Kode Unit</Label>
              <Input
                id="code"
                placeholder="Contoh: DAPUR"
                value={formData.code}
                onChange={(e) => setFormData({ ...formData, code: e.target.value.toUpperCase() })}
                maxLength={20}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Nama Unit</Label>
              <Input
                id="name"
                placeholder="Contoh: Dapur Pusat"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Deskripsi (Opsional)</Label>
              <Input
                id="description"
                placeholder="Keterangan tentang unit ini"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </div>
            <Button
              type="submit"
              className="w-full"
              disabled={createMutation.isPending || updateMutation.isPending}
            >
              {createMutation.isPending || updateMutation.isPending
                ? 'Menyimpan...'
                : editingBranch
                  ? 'Update Unit'
                  : 'Tambah Unit'}
            </Button>
          </form>
        </SheetContent>
      </Sheet>

      <AlertDialog
        open={confirmOpen}
        onOpenChange={(open) => {
          setConfirmOpen(open)
          if (!open) setBranchToDelete(null)
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Hapus Unit</AlertDialogTitle>
            <AlertDialogDescription>
              {branchToDelete
                ? `Hapus unit "${branchToDelete.name}"? Tindakan ini tidak dapat dibatalkan.`
                : 'Hapus unit ini?'}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteMutation.isPending}>Batal</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              disabled={deleteMutation.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteMutation.isPending ? 'Menghapus...' : 'Hapus'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
