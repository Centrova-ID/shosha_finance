import { Link, useLocation } from 'react-router-dom'
import { LayoutDashboard, Receipt, PlusCircle, Cloud, CloudOff } from 'lucide-react'
import { useSystemStatus } from '@/hooks/useDashboard'
import { cn } from '@/lib/utils'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()
  const { data: statusData } = useSystemStatus()

  const isOnline = statusData?.data?.status === 'online'
  const unsyncedCount = statusData?.data?.unsynced_count || 0

  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/transactions', label: 'Transaksi', icon: Receipt },
    { path: '/transactions/new', label: 'Input Baru', icon: PlusCircle }
  ]

  return (
    <div className="min-h-screen bg-background">
      <nav className="border-b bg-card">
        <div className="container mx-auto px-4">
          <div className="flex h-16 items-center justify-between">
            <div className="flex items-center gap-8">
              <h1 className="text-xl font-bold">Shosha Finance</h1>
              <div className="flex gap-1">
                {navItems.map((item) => (
                  <Link
                    key={item.path}
                    to={item.path}
                    className={cn(
                      'flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-colors',
                      location.pathname === item.path
                        ? 'bg-primary text-primary-foreground'
                        : 'hover:bg-accent hover:text-accent-foreground'
                    )}
                  >
                    <item.icon className="h-4 w-4" />
                    {item.label}
                  </Link>
                ))}
              </div>
            </div>

            <div className="flex items-center gap-4">
              {unsyncedCount > 0 && (
                <span className="text-sm text-muted-foreground">
                  {unsyncedCount} belum sync
                </span>
              )}
              <div
                className={cn(
                  'flex items-center gap-2 px-3 py-1 rounded-full text-sm',
                  isOnline ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                )}
              >
                {isOnline ? (
                  <>
                    <Cloud className="h-4 w-4" />
                    Online
                  </>
                ) : (
                  <>
                    <CloudOff className="h-4 w-4" />
                    Offline
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </nav>

      <main className="container mx-auto px-4 py-6">{children}</main>
    </div>
  )
}
