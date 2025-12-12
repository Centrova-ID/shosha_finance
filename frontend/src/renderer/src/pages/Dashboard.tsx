import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useDashboardSummary } from '@/hooks/useDashboard'
import { useActiveBranches } from '@/hooks/useBranches'
import { useAuth } from '@/contexts/AuthContext'
import { formatCurrency } from '@/lib/utils'
import { TrendingUp, TrendingDown, Wallet, RefreshCw, Users, Building2, Filter, Calendar } from 'lucide-react'
import TransactionSheet from '@/components/TransactionSheet'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from 'recharts'

// Get today's date in YYYY-MM-DD format
const getTodayDate = () => {
  const today = new Date()
  return today.toISOString().split('T')[0]
}

export default function Dashboard() {
  const { user } = useAuth()
  const [selectedBranch, setSelectedBranch] = useState<string>('all')
  const [selectedDate, setSelectedDate] = useState<string>(getTodayDate())
  const { data: branchesData } = useActiveBranches()
  const { data, isLoading, error, refetch } = useDashboardSummary({
    branchId: selectedBranch === 'all' ? undefined : selectedBranch,
    date: selectedDate || undefined
  })

  const branches = branchesData?.data || []

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
    {
      name: 'Ringkasan',
      pemasukan: summary?.total_in || 0,
      pengeluaran: summary?.total_out || 0
    }
  ]

  const getRoleGreeting = () => {
    switch (user?.role) {
      case 'admin':
        return 'Administrator Dashboard'
      case 'manager':
        return 'Manager Dashboard'
      case 'staff':
        return 'Staff Dashboard'
      default:
        return 'Dashboard'
    }
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header Section */}
      <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
        <div className="space-y-1">
          <h2 className="text-3xl font-bold tracking-tight">{getRoleGreeting()}</h2>
          <p className="text-muted-foreground">Selamat datang, {user?.name}</p>
        </div>
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
          <Card className="p-3">
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <Input
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="h-8 w-[140px] border-0 p-0 focus-visible:ring-0"
              />
            </div>
          </Card>
          <Card className="p-3">
            <div className="flex items-center gap-2">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <Select value={selectedBranch} onValueChange={setSelectedBranch}>
                <SelectTrigger className="h-8 w-[160px] border-0 p-0 focus:ring-0">
                  <SelectValue placeholder="Semua Unit" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Semua Unit</SelectItem>
                  {branches.map((branch) => (
                    <SelectItem key={branch.id} value={branch.id}>
                      {branch.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </Card>
          <TransactionSheet onSuccess={() => refetch()} />
        </div>
      </div>

      {user?.role === 'admin' && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-2">
          <Card className="border-l-4 border-l-primary">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Cabang</CardTitle>
              <Building2 className="h-4 w-4 text-primary" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">3</div>
              <p className="text-xs text-muted-foreground">Cabang aktif</p>
            </CardContent>
          </Card>
          <Card className="border-l-4 border-l-primary">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Pengguna</CardTitle>
              <Users className="h-4 w-4 text-primary" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">5</div>
              <p className="text-xs text-muted-foreground">Pengguna terdaftar</p>
            </CardContent>
          </Card>
        </div>
      )}

      <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
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
      </div>

      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Ringkasan Keuangan</CardTitle>
        </CardHeader>
        <CardContent className="pl-26">
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={chartData} margin={{ left: 80, right: 20, top: 20, bottom: 20 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis tickFormatter={(value) => formatCurrency(value)} />
              <Tooltip formatter={(value: number) => formatCurrency(value)} />
              <Bar dataKey="pemasukan" name="Pemasukan" fill="#22c55e" label={{ position: 'top', formatter: (value: number) => value > 0 ? formatCurrency(value) : 'Rp 0' }} />
              <Bar dataKey="pengeluaran" name="Pengeluaran" fill="#ef4444" label={{ position: 'top', formatter: (value: number) => value > 0 ? formatCurrency(value) : 'Rp 0' }} />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
    </div>
  )
}
