import { useState, useEffect } from 'react'
import { useBranches } from '@/hooks/useBranches'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table'
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
import { ExpenseEntryRequest } from '@/types/income'
import { Branch } from '@/types'
import { AlertTriangle, Save, RefreshCw } from 'lucide-react'
import { apiClient } from '@/api/client'

interface ExpenseRowData {
  branch_id: string
  branch_name: string
  omzet: string
  pengeluaran_toru: string
  pengeluaran_cash: number
  qris_bca: string
  qris_bni: string
  qris_bri: string
  transfer_bca: string
  transfer_bni: string
  transfer_bri: string
  difference?: number // positive = sisa, negative = kelebihan
  status?: 'over' | 'under'
  warning?: string
}

export default function ExpensePage() {
  const { data: branchesData, isLoading } = useBranches()
  const branches = branchesData?.data || []

  const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0])
  const [rows, setRows] = useState<ExpenseRowData[]>([])
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (branches.length > 0) {
      initializeRows(branches)
    }
  }, [branches])

  const initializeRows = (branchList: Branch[]) => {
    const initialRows: ExpenseRowData[] = branchList.map((branch) => ({
      branch_id: branch.id,
      branch_name: branch.name,
      omzet: '',
      pengeluaran_toru: '',
      pengeluaran_cash: 0,
      qris_bca: '',
      qris_bni: '',
      qris_bri: '',
      transfer_bca: '',
      transfer_bni: '',
      transfer_bri: '',
      warning: undefined
    }))
    setRows(initialRows)
  }

  const parseNumber = (value: string): number => {
    const parsed = parseInt(value.replace(/\D/g, ''), 10)
    return isNaN(parsed) ? 0 : parsed
  }

  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('id-ID').format(value)
  }

  const handleInputChange = (index: number, field: keyof ExpenseRowData, value: string) => {
    setRows((prev) => {
      const updated = [...prev]
      const row = { ...updated[index] }

      // Format the input value to currency format while allowing editing
      const formattedValue = value.replace(/\D/g, '') // Keep only digits
      const numericValue = parseInt(formattedValue, 10) || 0
      const displayValue = formattedValue ? formatCurrency(numericValue) : ''

      if (field === 'omzet' || field === 'pengeluaran_toru') {
        row[field] = displayValue
        const omzet = parseNumber(row.omzet)
        const toruInput = parseNumber(row.pengeluaran_toru)
        row.pengeluaran_cash = omzet - toruInput
      } else if (
        field === 'qris_bca' ||
        field === 'qris_bni' ||
        field === 'qris_bri' ||
        field === 'transfer_bca' ||
        field === 'transfer_bni' ||
        field === 'transfer_bri'
      ) {
        row[field] = displayValue
      }

      // Calculate total payments
      const totalPayments =
        parseNumber(row.qris_bca) +
        parseNumber(row.qris_bni) +
        parseNumber(row.qris_bri) +
        parseNumber(row.transfer_bca) +
        parseNumber(row.transfer_bni) +
        parseNumber(row.transfer_bri)

      // Validate with delta (positive = sisa, negative = kelebihan)
      const delta = row.pengeluaran_cash - totalPayments
      row.difference = delta

      if (delta < 0) {
        const over = Math.abs(delta)
        row.status = 'over'
        row.warning = `Cabang ${row.branch_name}: kelebihan ${formatCurrency(over)} (bayar ${formatCurrency(totalPayments)} > cash ${formatCurrency(row.pengeluaran_cash)})`
      } else if (delta > 0 && row.pengeluaran_cash > 0) {
        row.status = 'under'
        row.warning = `Cabang ${row.branch_name}: sisa ${formatCurrency(delta)} (cash ${formatCurrency(row.pengeluaran_cash)} > bayar ${formatCurrency(totalPayments)})`
      } else {
        row.status = undefined
        row.warning = undefined
      }

      updated[index] = row
      return updated
    })
  }

  const warnings = rows.filter((row) => row.warning)
  const hasWarnings = warnings.length > 0

  const handleSave = async () => {
    if (hasWarnings) {
      toast({
        title: 'Peringatan',
        description: 'Ada cabang yang belum balance (kelebihan atau sisa pembayaran)',
        variant: 'destructive'
      })
      return
    }

    const dataToSave = rows.filter((row) => parseNumber(row.omzet) > 0)

    if (dataToSave.length === 0) {
      toast({
        title: 'Peringatan',
        description: 'Tidak ada data untuk disimpan',
        variant: 'destructive'
      })
      return
    }

    setConfirmOpen(true)
  }

  const confirmSave = async () => {
    setSaving(true)
    setConfirmOpen(false)

    try {
      const requests: ExpenseEntryRequest[] = rows
        .filter((row) => parseNumber(row.omzet) > 0)
        .map((row) => ({
          branch_id: row.branch_id,
          date: selectedDate,
          omzet: parseNumber(row.omzet),
          pengeluaran_toru: parseNumber(row.pengeluaran_toru),
          qris_bca: parseNumber(row.qris_bca),
          qris_bni: parseNumber(row.qris_bni),
          qris_bri: parseNumber(row.qris_bri),
          transfer_bca: parseNumber(row.transfer_bca),
          transfer_bni: parseNumber(row.transfer_bni),
          transfer_bri: parseNumber(row.transfer_bri)
        }))

      const results = await Promise.allSettled(
        requests.map((req) => apiClient.post('/expense-entries', req))
      )

      const successful = results.filter((r) => r.status === 'fulfilled').length
      const failed = results.filter((r) => r.status === 'rejected').length

      if (failed > 0) {
        toast({
          title: 'Sebagian berhasil disimpan',
          description: `${successful} berhasil, ${failed} gagal`,
          variant: 'default'
        })
      } else {
        toast({
          title: 'Berhasil',
          description: `${successful} data pengeluaran berhasil disimpan`
        })
        initializeRows(branches)
      }
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.message || 'Gagal menyimpan data',
        variant: 'destructive'
      })
    } finally {
      setSaving(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight">Input Pengeluaran</h2>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Label htmlFor="date">Tanggal:</Label>
            <Input
              id="date"
              type="date"
              value={selectedDate}
              onChange={(e) => setSelectedDate(e.target.value)}
              className="w-40"
            />
          </div>
          <Button onClick={handleSave} disabled={saving || hasWarnings}>
            <Save className="mr-2 h-4 w-4" />
            {saving ? 'Menyimpan...' : 'Simpan Semua'}
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Tabel Input Pengeluaran per Cabang</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[150px]">Nama Cabang</TableHead>
                  <TableHead className="w-[120px]">Omzet</TableHead>
                  <TableHead className="w-[120px]">Pengeluaran Toru</TableHead>
                  <TableHead className="w-[120px] bg-muted/50">Pengeluaran Cash</TableHead>
                  <TableHead className="w-[120px]">QRIS BCA</TableHead>
                  <TableHead className="w-[120px]">QRIS BNI</TableHead>
                  <TableHead className="w-[120px]">QRIS BRI</TableHead>
                  <TableHead className="w-[120px]">TF BCA</TableHead>
                  <TableHead className="w-[120px]">TF BNI</TableHead>
                  <TableHead className="w-[120px]">TF BRI</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((row, index) => (
                  <TableRow key={row.branch_id} className={row.warning ? 'bg-red-50' : ''}>
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        {row.branch_name}
                        {row.warning && (
                          <AlertTriangle className="h-4 w-4 text-destructive" />
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.omzet}
                        onChange={(e) => handleInputChange(index, 'omzet', e.target.value)}
                        placeholder="0"
                        className="w-full"
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.pengeluaran_toru}
                        onChange={(e) => handleInputChange(index, 'pengeluaran_toru', e.target.value)}
                        placeholder="0"
                        className="w-full"
                      />
                    </TableCell>
                    <TableCell className="bg-muted/50">
                      <div className="font-medium text-sm">
                        {formatCurrency(row.pengeluaran_cash)}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.qris_bca}
                        onChange={(e) => handleInputChange(index, 'qris_bca', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.qris_bni}
                        onChange={(e) => handleInputChange(index, 'qris_bni', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.qris_bri}
                        onChange={(e) => handleInputChange(index, 'qris_bri', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.transfer_bca}
                        onChange={(e) => handleInputChange(index, 'transfer_bca', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.transfer_bni}
                        onChange={(e) => handleInputChange(index, 'transfer_bni', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                    <TableCell>
                      <Input
                        type="text"
                        value={row.transfer_bri}
                        onChange={(e) => handleInputChange(index, 'transfer_bri', e.target.value)}
                        placeholder="0"
                        className="w-full"
                        disabled={row.pengeluaran_cash <= 0}
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          {hasWarnings && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md space-y-2">
              <p className="text-sm text-red-800 font-medium flex items-center gap-2">
                <AlertTriangle className="h-4 w-4" />
                Peringatan: Ada cabang yang belum balance. Detail:
              </p>
              <ul className="text-sm text-red-800 list-disc list-inside space-y-1">
                {warnings.map((w) => (
                  <li key={w.branch_id}>{w.warning}</li>
                ))}
              </ul>
            </div>
          )}
        </CardContent>
      </Card>

      <AlertDialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Konfirmasi Simpan</AlertDialogTitle>
            <AlertDialogDescription>
              Simpan data pengeluaran untuk tanggal {selectedDate}?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Batal</AlertDialogCancel>
            <AlertDialogAction onClick={confirmSave}>Simpan</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
