import { useState } from 'react'
import { useCreateTransaction } from '@/hooks/useTransactions'
import { useActiveBranches } from '@/hooks/useBranches'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger
} from '@/components/ui/sheet'
import { toast } from '@/hooks/use-toast'
import { TransactionType } from '@/types'
import { PlusCircle, Save } from 'lucide-react'

const categories = {
  IN: ['Penjualan', 'Setoran Modal', 'Piutang Dibayar', 'Lainnya'],
  OUT: ['Bahan Baku', 'Operasional', 'Gaji', 'Listrik', 'Gas', 'Transport', 'Lainnya']
}

interface TransactionSheetProps {
  onSuccess?: () => void
}

export default function TransactionSheet({ onSuccess }: TransactionSheetProps) {
  const createMutation = useCreateTransaction()
  const { data: branchesData } = useActiveBranches()
  const [open, setOpen] = useState(false)
  const [branchId, setBranchId] = useState('')
  const [type, setType] = useState<TransactionType>('OUT')
  const [category, setCategory] = useState('')
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')

  const branches = branchesData?.data || []

  const resetForm = () => {
    setBranchId('')
    setType('OUT')
    setCategory('')
    setAmount('')
    setDescription('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!branchId) {
      toast({
        title: 'Error',
        description: 'Pilih unit terlebih dahulu',
        variant: 'destructive'
      })
      return
    }

    if (!category || !amount) {
      toast({
        title: 'Error',
        description: 'Kategori dan jumlah harus diisi',
        variant: 'destructive'
      })
      return
    }

    const amountNum = parseInt(amount.replace(/\D/g, ''), 10)
    if (isNaN(amountNum) || amountNum <= 0) {
      toast({
        title: 'Error',
        description: 'Jumlah harus berupa angka positif',
        variant: 'destructive'
      })
      return
    }

    try {
      await createMutation.mutateAsync({
        branch_id: branchId,
        type,
        category,
        amount: amountNum,
        description: description || undefined
      })

      toast({
        title: 'Berhasil',
        description: 'Transaksi berhasil disimpan'
      })

      resetForm()
      setOpen(false)
      onSuccess?.()
    } catch {
      toast({
        title: 'Error',
        description: 'Gagal menyimpan transaksi',
        variant: 'destructive'
      })
    }
  }

  const formatAmount = (value: string) => {
    const num = value.replace(/\D/g, '')
    return num.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  }

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button>
          <PlusCircle className="mr-2 h-4 w-4" />
          Input Transaksi
        </Button>
      </SheetTrigger>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>Input Transaksi Baru</SheetTitle>
          <SheetDescription>Masukkan detail transaksi keuangan</SheetDescription>
        </SheetHeader>
        <form onSubmit={handleSubmit} className="space-y-6 mt-6">
          <div className="space-y-2">
            <Label htmlFor="branch">Unit</Label>
            <Select value={branchId} onValueChange={setBranchId}>
              <SelectTrigger>
                <SelectValue placeholder="Pilih unit" />
              </SelectTrigger>
              <SelectContent>
                {branches.map((branch) => (
                  <SelectItem key={branch.id} value={branch.id}>
                    {branch.name} ({branch.code})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Tipe Transaksi</Label>
            <div className="flex gap-2">
              <Button
                type="button"
                variant={type === 'IN' ? 'default' : 'outline'}
                className={type === 'IN' ? 'bg-green-600 hover:bg-green-700 flex-1' : 'flex-1'}
                onClick={() => {
                  setType('IN')
                  setCategory('')
                }}
              >
                Pemasukan
              </Button>
              <Button
                type="button"
                variant={type === 'OUT' ? 'default' : 'outline'}
                className={type === 'OUT' ? 'bg-red-600 hover:bg-red-700 flex-1' : 'flex-1'}
                onClick={() => {
                  setType('OUT')
                  setCategory('')
                }}
              >
                Pengeluaran
              </Button>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="category">Kategori</Label>
            <Select value={category} onValueChange={setCategory}>
              <SelectTrigger>
                <SelectValue placeholder="Pilih kategori" />
              </SelectTrigger>
              <SelectContent>
                {categories[type].map((cat) => (
                  <SelectItem key={cat} value={cat}>
                    {cat}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="amount">Jumlah (Rp)</Label>
            <Input
              id="amount"
              type="text"
              placeholder="0"
              value={amount}
              onChange={(e) => setAmount(formatAmount(e.target.value))}
              className="text-lg font-medium"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Keterangan (Opsional)</Label>
            <Input
              id="description"
              type="text"
              placeholder="Tambahkan keterangan..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>

          <Button type="submit" className="w-full" disabled={createMutation.isPending}>
            <Save className="mr-2 h-4 w-4" />
            {createMutation.isPending ? 'Menyimpan...' : 'Simpan Transaksi'}
          </Button>
        </form>
      </SheetContent>
    </Sheet>
  )
}
