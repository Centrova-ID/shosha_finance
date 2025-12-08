import { app, shell, BrowserWindow } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import { spawn, ChildProcess } from 'child_process'

let backendProcess: ChildProcess | null = null

function startBackend(): void {
  const appPath = app.getAppPath()
  const backendDir = join(appPath, '..', 'backend')
  const backendPath = is.dev
    ? join(backendDir, 'cmd/local/main.go')
    : join(process.resourcesPath, 'backend', 'local-backend')

  const commonEnv = {
    ...process.env,
    PORT: '8080',
    SQLITE_PATH: join(app.getPath('userData'), 'shosha_finance.db'),
    CLOUD_API_URL: process.env.CLOUD_API_URL || 'http://localhost:3000',
    SYNC_INTERVAL: process.env.SYNC_INTERVAL || '30'
  }

  if (is.dev) {
    backendProcess = spawn('go', ['run', backendPath], {
      cwd: backendDir,
      env: commonEnv
    })
  } else {
    backendProcess = spawn(backendPath, [], {
      env: commonEnv
    })
  }

  backendProcess.stdout?.on('data', (data) => {
    console.log(`Backend: ${data}`)
  })

  backendProcess.stderr?.on('data', (data) => {
    console.error(`Backend Error: ${data}`)
  })

  backendProcess.on('close', (code) => {
    console.log(`Backend process exited with code ${code}`)
  })
}

function createWindow(): void {
  const mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    show: false,
    autoHideMenuBar: true,
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })

  mainWindow.on('ready-to-show', () => {
    mainWindow.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

app.whenReady().then(() => {
  electronApp.setAppUserModelId('com.shosha.finance')

  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  startBackend()

  setTimeout(() => {
    createWindow()
  }, 2000)

  app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (backendProcess) {
    backendProcess.kill()
  }
  if (process.platform !== 'darwin') {
    app.quit()
  }
})
