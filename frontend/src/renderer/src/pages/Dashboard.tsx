import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useDashboardSummary } from '@/hooks/useDashboard'
import { formatCurrency } from '@/lib/utils'
import { TrendingUp, TrendingDown, Wallet, RefreshCw } from 'lucide-react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from 'recharts'

export default function Dashboard() {
  const { data, isLoading, error } = useDashboardSummary()

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
        Gagal memuat data dashboard. Pastikan backend sudah berjalan.
      </div>
    )
  }

  const summary = data?.data

  const chartData = [
    { name: 'Pemasukan', amount: summary?.total_in || 0, fill: '#22c55e' },
    { name: 'Pengeluaran', amount: summary?.total_out || 0, fill: '#ef4444' }
  ]

  return (
    <div className="space-y-6">
      <h2 className="text-3xl font-bold tracking-tight">Dashboard</h2>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Pemasukan</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {formatCurrency(summary?.total_in || 0)}
            </div>
            <p className="text-xs text-muted-foreground">{summary?.count_in || 0} transaksi</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Pengeluaran</CardTitle>
            <TrendingDown className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {formatCurrency(summary?.total_out || 0)}
            </div>
            <p className="text-xs text-muted-foreground">{summary?.count_out || 0} transaksi</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo</CardTitle>
            <Wallet className="h-4 w-4 text-blue-500" />
          </CardHeader>
          <CardContent>
            <div
              className={`text-2xl font-bold ${
                (summary?.balance || 0) >= 0 ? 'text-blue-600' : 'text-red-600'
              }`}
            >
              {formatCurrency(summary?.balance || 0)}
            </div>
            <p className="text-xs text-muted-foreground">Saldo saat ini</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Belum Sync</CardTitle>
            <RefreshCw className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">
              {summary?.unsync_count || 0}
            </div>
            <p className="text-xs text-muted-foreground">Menunggu sinkronisasi</p>
          </CardContent>
        </Card>
      </div>

      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Ringkasan Keuangan</CardTitle>
        </CardHeader>
        <CardContent className="pl-2">
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis tickFormatter={(value) => formatCurrency(value)} />
              <Tooltip formatter={(value: number) => formatCurrency(value)} />
              <Bar dataKey="amount" fill="#8884d8" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
    </div>
  )
}
