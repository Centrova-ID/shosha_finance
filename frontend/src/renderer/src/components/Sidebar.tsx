import { Link, useLocation } from 'react-router-dom'
import { useAuth, UserRole } from '@/contexts/AuthContext'
import { useSystemStatus } from '@/hooks/useDashboard'
import { cn } from '@/lib/utils'
import {
  LayoutDashboard,
  Receipt,
  Users,
  Settings,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Building2,
  Cloud,
  CloudOff,
  Loader2,
  TrendingUp,
  TrendingDown
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useState, useEffect } from 'react'

interface NavItem {
  path: string
  label: string
  icon: React.ElementType
  roles: UserRole[]
}

const navItems: NavItem[] = [
  { path: '/', label: 'Dashboard', icon: LayoutDashboard, roles: ['admin', 'manager', 'staff'] },
  { path: '/income', label: 'Pemasukan', icon: TrendingUp, roles: ['admin', 'manager', 'staff'] },
  { path: '/expense', label: 'Pengeluaran', icon: TrendingDown, roles: ['admin', 'manager', 'staff'] },
  { path: '/transactions', label: 'Transaksi', icon: Receipt, roles: ['admin', 'manager', 'staff'] },
  { path: '/branches', label: 'Unit', icon: Building2, roles: ['admin', 'manager'] },
  { path: '/users', label: 'Pengguna', icon: Users, roles: ['admin'] },
  { path: '/settings', label: 'Pengaturan', icon: Settings, roles: ['admin', 'manager'] }
]

export default function Sidebar() {
  const location = useLocation()
  const { user, logout } = useAuth()
  const { data: statusData } = useSystemStatus()
  const [collapsed, setCollapsed] = useState(false)
  const [isSyncing, setIsSyncing] = useState(false)
  const [lastSyncTime, setLastSyncTime] = useState<Date | null>(null)

  const isOnline = statusData?.data?.status === 'online'
  const unsyncedCount = statusData?.data?.unsynced_count || 0
  const [prevUnsyncedCount, setPrevUnsyncedCount] = useState(unsyncedCount)

  // Detect when sync happens (unsynced count decreases)
  useEffect(() => {
    if (isOnline && prevUnsyncedCount > 0 && unsyncedCount === 0) {
      setIsSyncing(true)
      setLastSyncTime(new Date())
      setTimeout(() => setIsSyncing(false), 2000)
    }
    setPrevUnsyncedCount(unsyncedCount)
  }, [unsyncedCount, isOnline])

  const filteredNavItems = navItems.filter((item) => user && item.roles.includes(user.role))

  const formatLastSync = () => {
    if (!lastSyncTime) return null
    const now = new Date()
    const diff = Math.floor((now.getTime() - lastSyncTime.getTime()) / 1000)
    if (diff < 60) return `${diff}d yang lalu`
    if (diff < 3600) return `${Math.floor(diff / 60)}m yang lalu`
    return lastSyncTime.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <aside
      className={cn(
        'flex flex-col h-screen bg-card border-r transition-all duration-300',
        collapsed ? 'w-16' : 'w-64'
      )}
    >
      <div className="flex items-center justify-between h-16 px-4 border-b">
        {!collapsed && (
          <h1 className="text-lg font-bold text-foreground">Shosha Finance</h1>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setCollapsed(!collapsed)}
          className={cn('h-8 w-8', collapsed && 'mx-auto')}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </Button>
      </div>

      <nav className="flex-1 p-2 space-y-1">
        {filteredNavItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className={cn(
              'flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors',
              location.pathname === item.path
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
              collapsed && 'justify-center'
            )}
            title={collapsed ? item.label : undefined}
          >
            <item.icon className="h-5 w-5 flex-shrink-0" />
            {!collapsed && <span>{item.label}</span>}
          </Link>
        ))}
      </nav>

      <div className="p-2 border-t space-y-2">
        {!collapsed && unsyncedCount > 0 && (
          <div className="px-3 py-2 text-sm text-amber-600 bg-amber-50 rounded-md">
            {unsyncedCount} belum sync
          </div>
        )}

        {!collapsed && isSyncing && (
          <div className="px-3 py-2 text-sm text-blue-600 bg-blue-50 rounded-md flex items-center gap-2">
            <Loader2 className="h-3 w-3 animate-spin" />
            Syncing...
          </div>
        )}

        {!collapsed && lastSyncTime && !isSyncing && unsyncedCount === 0 && (
          <div className="px-3 py-2 text-xs text-green-600">
            ✓ Sync {formatLastSync()}
          </div>
        )}

        <div
          className={cn(
            'flex items-center gap-2 px-3 py-2 rounded-md text-sm',
            isOnline ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700',
            collapsed && 'justify-center'
          )}
          title={collapsed ? (isOnline ? 'Online' : 'Offline') : undefined}
        >
          {isSyncing ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : isOnline ? (
            <Cloud className="h-4 w-4" />
          ) : (
            <CloudOff className="h-4 w-4" />
          )}
          {!collapsed && <span>{isOnline ? 'Online' : 'Offline'}</span>}
        </div>

        {!collapsed && user && (
          <div className="px-3 py-2 border-t">
            <p className="text-sm font-medium text-foreground">{user.name}</p>
            <p className="text-xs text-muted-foreground capitalize">{user.role}</p>
          </div>
        )}

        <Button
          variant="ghost"
          className={cn(
            'w-full justify-start gap-3 text-muted-foreground hover:text-destructive hover:bg-destructive/10',
            collapsed && 'justify-center'
          )}
          onClick={logout}
          title={collapsed ? 'Logout' : undefined}
        >
          <LogOut className="h-5 w-5" />
          {!collapsed && <span>Logout</span>}
        </Button>
      </div>
    </aside>
  )
}
