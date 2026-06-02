import { ref } from 'vue'
import {
  IsSetupComplete,
  GenerateProxyAccessKey,
  GetEnvFile,
  SaveEnvFile,
  StartService,
} from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'
import { detectEnvNewline, parseEnvFile, serializeEnvFile } from '@/lib/env-file'
import { useLanguage } from '@/composables/useLanguage'
import type { TabValue } from '@/types'

// 模块级单例：所有调用方共享同一份 setup 状态
const setupChecked = ref(false)
const setupComplete = ref(false)
const setupKey = ref('')
const setupSaving = ref(false)
const setupError = ref('')
const envPath = ref('')
const pendingTab = ref<TabValue | null>(null)
let checkPromise: Promise<void> | null = null

const SKIP_SETUP_PARAM = 'skip_setup'

const checkSetup = async () => {
  if (checkPromise) return checkPromise
  checkPromise = (async () => {
    // 开发调试：URL 带 ?skip_setup 参数时跳过 Setup 引导，直接进入主界面
    if (typeof window !== 'undefined') {
      const params = new URLSearchParams(window.location.search)
      if (params.has(SKIP_SETUP_PARAM)) {
        setupComplete.value = true
        setupChecked.value = true
        return
      }
    }
    setupError.value = ''
    try {
      const [done, env] = await Promise.all([IsSetupComplete(), GetEnvFile()])
      envPath.value = env.path || ''
      if (done) {
        setupComplete.value = true
      } else {
        setupComplete.value = false
        try {
          setupKey.value = await GenerateProxyAccessKey()
        } catch (err) {
          setupError.value = err instanceof Error ? err.message : String(err)
        }
      }
    } catch (err) {
      setupError.value = err instanceof Error ? err.message : String(err)
    } finally {
      setupChecked.value = true
    }
  })()
  return checkPromise
}

const confirmSetup = async (key: string, target: TabValue = 'agent') => {
  const trimmed = key.trim()
  if (!trimmed) {
    setupError.value = useLanguage().t('setup.errorEmptyKey')
    return false
  }
  setupSaving.value = true
  setupError.value = ''
  try {
    // 合并已有 .env 内容：保留注释与现有键值，仅注入/覆盖 PROXY_ACCESS_KEY
    const env = await GetEnvFile()
    const content = env.content || ''
    const newline = detectEnvNewline(content)
    const entries = parseEnvFile(content)
    const serialized = serializeEnvFile(entries, { PROXY_ACCESS_KEY: trimmed }, ['PROXY_ACCESS_KEY'], newline)
    await SaveEnvFile(serialized)
    setupComplete.value = true
    envPath.value = env.path || envPath.value
    // 触发自动启动（失败不阻断进入主界面）
    try {
      await StartService()
    } catch (err) {
      setupError.value = err instanceof Error ? err.message : String(err)
    }
    pendingTab.value = target
    return true
  } catch (err) {
    setupError.value = err instanceof Error ? err.message : String(err)
    return false
  } finally {
    setupSaving.value = false
  }
}

const regenerateKey = async () => {
  setupError.value = ''
  try {
    setupKey.value = await GenerateProxyAccessKey()
  } catch (err) {
    setupError.value = err instanceof Error ? err.message : String(err)
  }
}

export function useSetup() {
  return {
    setupChecked,
    setupComplete,
    setupKey,
    setupSaving,
    setupError,
    envPath,
    pendingTab,
    checkSetup,
    confirmSetup,
    regenerateKey,
  }
}
