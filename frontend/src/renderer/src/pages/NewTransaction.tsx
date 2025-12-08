import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCreateTransaction } from '@/hooks/useTransactions'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { toast } from '@/hooks/use-toast'
import { TransactionType } from '@/types'
import { ArrowLeft, Save } from 'lucide-react'

const categories = {
  IN: ['Penjualan', 'Setoran Modal', 'Piutang Dibayar', 'Lainnya'],
  OUT: ['Bahan Baku', 'Operasional', 'Gaji', 'Listrik', 'Gas', 'Transport', 'Lainnya']
}

export default function NewTransaction() {
  const navigate = useNavigate()
  const createMutation = useCreateTransaction()

  const [type, setType] = useState<TransactionType>('OUT')
  const [category, setCategory] = useState('')
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

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
        type,
        category,
        amount: amountNum,
        description: description || undefined
      })

      toast({
        title: 'Berhasil',
        description: 'Transaksi berhasil disimpan'
      })

      navigate('/transactions')
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
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h2 className="text-3xl font-bold tracking-tight">Input Transaksi Baru</h2>
      </div>

      <Card className="max-w-xl">
        <CardHeader>
          <CardTitle>Form Transaksi</CardTitle>
          <CardDescription>Masukkan detail transaksi keuangan</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <Label>Tipe Transaksi</Label>
              <div className="flex gap-4">
                <Button
                  type="button"
                  variant={type === 'IN' ? 'default' : 'outline'}
                  className={type === 'IN' ? 'bg-green-600 hover:bg-green-700' : ''}
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
                  className={type === 'OUT' ? 'bg-red-600 hover:bg-red-700' : ''}
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

            <Button
              type="submit"
              className="w-full"
              disabled={createMutation.isPending}
            >
              <Save className="mr-2 h-4 w-4" />
              {createMutation.isPending ? 'Menyimpan...' : 'Simpan Transaksi'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
