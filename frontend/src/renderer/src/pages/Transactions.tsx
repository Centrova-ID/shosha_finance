import { useState } from 'react'
import { useTransactions } from '@/hooks/useTransactions'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { formatCurrency, formatDate } from '@/lib/utils'
import { RefreshCw, ChevronLeft, ChevronRight } from 'lucide-react'
import TransactionSheet from '@/components/TransactionSheet'

export default function Transactions() {
  const [page, setPage] = useState(1)
  const limit = 10
  const { data, isLoading, error, refetch } = useTransactions(page, limit)

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center text-destructive">
        Gagal memuat data transaksi. Pastikan backend sudah berjalan.
      </div>
    )
  }

  const transactions = data?.data || []
  const meta = data?.meta

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">Daftar Transaksi</h2>
        <TransactionSheet onSuccess={() => refetch()} />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Riwayat Transaksi</CardTitle>
        </CardHeader>
        <CardContent>
          {transactions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              Belum ada transaksi. Mulai input transaksi pertama Anda!
            </div>
          ) : (
            <div className="space-y-4">
              <div className="rounded-md border">
                <table className="w-full">
                  <thead>
                    <tr className="border-b bg-muted/50">
                      <th className="p-3 text-left font-medium">Tanggal</th>
                      <th className="p-3 text-left font-medium">Unit</th>
                      <th className="p-3 text-left font-medium">Tipe</th>
                      <th className="p-3 text-left font-medium">Kategori</th>
                      <th className="p-3 text-left font-medium">Keterangan</th>
                      <th className="p-3 text-right font-medium">Jumlah</th>
                    </tr>
                  </thead>
                  <tbody>
                    {transactions.map((tx) => (
                      <tr key={tx.id} className="border-b last:border-0">
                        <td className="p-3 text-sm">{formatDate(tx.created_at)}</td>
                        <td className="p-3 text-sm">
                          <span className="inline-flex items-center rounded-full bg-secondary px-2 py-1 text-xs font-medium">
                            {tx.branch?.name || tx.branch_id.slice(0, 8)}
                          </span>
                        </td>
                        <td className="p-3">
                          <span
                            className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                              tx.type === 'IN'
                                ? 'bg-green-100 text-green-700'
                                : 'bg-red-100 text-red-700'
                            }`}
                          >
                            {tx.type === 'IN' ? 'Masuk' : 'Keluar'}
                          </span>
                        </td>
                        <td className="p-3 text-sm">{tx.category}</td>
                        <td className="p-3 text-sm text-muted-foreground">
                          {tx.description || '-'}
                        </td>
                        <td
                          className={`p-3 text-right font-medium ${
                            tx.type === 'IN' ? 'text-green-600' : 'text-red-600'
                          }`}
                        >
                          {tx.type === 'IN' ? '+' : '-'}
                          {formatCurrency(tx.amount)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {meta && meta.total_pages > 1 && (
                <div className="flex items-center justify-between">
                  <p className="text-sm text-muted-foreground">
                    Halaman {meta.page} dari {meta.total_pages} ({meta.total} total)
                  </p>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Sebelumnya
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPage((p) => p + 1)}
                      disabled={page >= meta.total_pages}
                    >
                      Selanjutnya
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
