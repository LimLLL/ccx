import { reactive, ref, onMounted, onBeforeUnmount } from 'vue'
import { Events } from '@wailsio/runtime'
import {
  CheckUpdate,
  DownloadAndInstall,
  CancelUpdate,
  GetVersion,
} from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'
import type {
  UpdateInfo,
  VersionInfo,
} from '@bindings/github.com/BenedictKing/ccx/desktop/models'

// Phase 与 Go 端 updater.Phase 字符串保持一致
export type UpdaterPhase =
  | ''
  | 'checking'
  | 'downloading'
  | 'verifying'
  | 'installing'
  | 'done'
  | 'error'

interface UpdaterState {
  checking: boolean
  downloading: boolean
  available: boolean
  dialogOpen: boolean
  version: VersionInfo | null
  info: UpdateInfo | null
  phase: UpdaterPhase
  percent: number
  downloaded: number
  total: number
  error: string
}

const state = reactive<UpdaterState>({
  checking: false,
  downloading: false,
  available: false,
  dialogOpen: false,
  version: null,
  info: null,
  phase: '',
  percent: 0,
  downloaded: 0,
  total: 0,
  error: '',
})

let unsubProgress: (() => void) | null = null
let unsubAvailable: (() => void) | null = null
let listenersInitialized = false

const installEventListeners = () => {
  if (listenersInitialized) return
  listenersInitialized = true

  unsubProgress = Events.On('update:progress', (event: unknown) => {
    const data = (event as { data?: unknown[] })?.data?.[0] as
      | {
          phase?: UpdaterPhase
          percent?: number
          downloaded?: number
          total?: number
          error?: string
        }
      | undefined
    if (!data) return
    state.phase = data.phase ?? ''
    state.percent = data.percent ?? 0
    state.downloaded = data.downloaded ?? 0
    state.total = data.total ?? 0
    state.error = data.error ?? ''
    if (data.phase === 'done') {
      state.downloading = false
    } else if (data.phase === 'error') {
      state.downloading = false
    } else if (data.phase === 'downloading' || data.phase === 'verifying' || data.phase === 'installing') {
      state.downloading = true
    }
  })

  unsubAvailable = Events.On('update:available', (event: unknown) => {
    const data = (event as { data?: unknown[] })?.data?.[0] as UpdateInfo | undefined
    if (!data) return
    state.info = data
    state.available = !!data.available
    if (state.available && !state.downloading) {
      state.dialogOpen = true
    }
  })
}

const teardownEventListeners = () => {
  if (unsubProgress) {
    unsubProgress()
    unsubProgress = null
  }
  if (unsubAvailable) {
    unsubAvailable()
    unsubAvailable = null
  }
  listenersInitialized = false
}

const initialized = ref(false)

const syncVersion = async () => {
  try {
    state.version = await GetVersion()
  } catch {
    // 版本未注入时静默
  }
}

const check = async () => {
  state.checking = true
  state.error = ''
  try {
    const info = await CheckUpdate()
    state.info = info
    state.available = !!info?.available
    if (state.available) {
      state.dialogOpen = true
    }
    return info
  } catch (err) {
    state.error = err instanceof Error ? err.message : String(err)
    throw err
  } finally {
    state.checking = false
  }
}

const downloadAndInstall = async () => {
  if (!state.info) return
  state.downloading = true
  state.error = ''
  state.percent = 0
  state.downloaded = 0
  state.total = state.info.size ?? 0
  try {
    await DownloadAndInstall(state.info)
  } catch (err) {
    state.error = err instanceof Error ? err.message : String(err)
    state.downloading = false
    throw err
  }
}

const cancel = async () => {
  try {
    await CancelUpdate()
  } finally {
    state.downloading = false
    state.phase = ''
  }
}

const closeDialog = () => {
  state.dialogOpen = false
}

export function useUpdater() {
  onMounted(async () => {
    if (!initialized.value) {
      initialized.value = true
      installEventListeners()
      await syncVersion()
    }
  })

  onBeforeUnmount(() => {
    // 单例事件订阅在整个应用生命周期内保留，组件卸载不取消
  })

  return {
    state,
    check,
    downloadAndInstall,
    cancel,
    closeDialog,
    syncVersion,
    teardownEventListeners,
  }
}
