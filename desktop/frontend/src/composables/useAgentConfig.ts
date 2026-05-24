import { ref } from 'vue'
import type { AgentPlatform, AgentProvider, AgentConfigStatus, ApplyAgentConfigRequest } from '@/types'
import {
  GetAgentConfigStatus,
  ApplyAgentConfig,
  RestoreAgentConfig,
  GetSavedProviderKeys,
} from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'

const agentLabels: Record<AgentPlatform, string> = {
  claude: 'Claude Code',
  codex: 'Codex',
}

const claudeProviderLabels: Record<AgentProvider | 'custom', string> = {
  ccx: 'CCX',
  deepseek: 'DeepSeek',
  mimo: 'MiMo',
  kimi: 'Kimi',
  glm: 'GLM',
  minimax: 'MiniMax',
  dashscope: 'DashScope',
  openai: 'OpenAI',
  custom: '自定义',
}

const codexProviderLabels: Record<AgentProvider | 'custom', string> = {
  ccx: 'CCX 本地网关',
  openai: 'OpenAI 官方',
  deepseek: 'DeepSeek',
  mimo: 'MiMo',
  kimi: 'Kimi',
  glm: 'GLM',
  minimax: 'MiniMax',
  dashscope: 'DashScope',
  custom: '自定义',
}

const agentPlatforms: AgentPlatform[] = ['claude', 'codex']

// Module-level singletons
const agentStatuses = ref<Record<AgentPlatform, AgentConfigStatus | null>>({
  claude: null,
  codex: null,
})
const configLoading = ref(false)
const selectedClaudeProvider = ref<AgentProvider>('ccx')
const claudeProviderKeys = ref<Record<AgentProvider, string>>({
  ccx: '',
  deepseek: '',
  mimo: '',
  kimi: '',
  glm: '',
  minimax: '',
  dashscope: '',
  openai: '',
})
const savedProviderKeys = ref<Record<string, string>>({})
const codexOpenAIKey = ref('')
const claudeMiMoBaseUrl = ref('https://api.xiaomimimo.com/anthropic')
const selectedMiMoPlan = ref('https://api.xiaomimimo.com/anthropic')
const selectedCodexProvider = ref<AgentProvider>('ccx')

const isClaudeProvider = (value?: string): value is AgentProvider => {
  return value === 'ccx' || value === 'deepseek' || value === 'mimo' || value === 'kimi' || value === 'glm' || value === 'minimax' || value === 'dashscope'
}

const claudeProviderLabel = (value?: string) => {
  if (!value) return '未识别'
  return claudeProviderLabels[value as AgentProvider | 'custom'] || value
}

const codexProviderLabel = (value?: string) => {
  if (!value) return '未识别'
  return codexProviderLabels[value as AgentProvider | 'custom'] || value
}

const claudeTargetBaseUrl = () => {
  switch (selectedClaudeProvider.value) {
    case 'ccx':
      return agentStatuses.value.claude?.targetBaseUrl || '当前 CCX 网关'
    case 'deepseek':
      return 'https://api.deepseek.com/anthropic'
    case 'mimo':
      return claudeMiMoBaseUrl.value || 'https://api.xiaomimimo.com/anthropic'
    case 'kimi':
      return 'https://api.moonshot.cn/anthropic'
    case 'glm':
      return 'https://open.bigmodel.cn/api/anthropic'
    case 'minimax':
      return 'https://api.minimaxi.com/anthropic'
    case 'dashscope':
      return 'https://dashscope.aliyuncs.com/apps/anthropic'
    default:
      return ''
  }
}

const agentStatusText = (item: AgentConfigStatus | null) => {
  if (!item) return '检测中'
  if (item.configured) return '已配置'
  if (item.needsUpdate) return '端口不匹配'
  return '未配置'
}

const agentStatusClass = (item: AgentConfigStatus | null) => {
  if (!item) return 'starting'
  if (item.configured) return 'running'
  if (item.needsUpdate) return 'starting'
  return 'stopped'
}

const resolveMiMoPlan = (url: string): string => {
  const known = [
    'https://api.xiaomimimo.com/anthropic',
    'https://token-plan-cn.xiaomimimo.com/anthropic',
    'https://token-plan-sgp.xiaomimimo.com/anthropic',
    'https://token-plan-ams.xiaomimimo.com/anthropic',
  ]
  return known.includes(url) ? url : ''
}

const loadAgentStatuses = async () => {
  configLoading.value = true
  try {
    const [claude, codex, keys] = await Promise.all([
      GetAgentConfigStatus('claude') as Promise<AgentConfigStatus>,
      GetAgentConfigStatus('codex') as Promise<AgentConfigStatus>,
      GetSavedProviderKeys(),
    ])
    agentStatuses.value = { claude, codex }
    savedProviderKeys.value = Object.fromEntries(
      Object.entries(keys).filter((entry): entry is [string, string] => typeof entry[1] === 'string')
    )
    if (isClaudeProvider(claude.provider)) {
      selectedClaudeProvider.value = claude.provider
    }
    if (claude.provider === 'mimo' && claude.currentBaseUrl) {
      claudeMiMoBaseUrl.value = claude.currentBaseUrl
      selectedMiMoPlan.value = resolveMiMoPlan(claude.currentBaseUrl)
    }
    if (codex.provider === 'openai') {
      selectedCodexProvider.value = 'openai'
    } else {
      selectedCodexProvider.value = 'ccx'
    }
  } catch (error) {
    // error is handled by caller
  } finally {
    configLoading.value = false
  }
}

const findSavedKey = (provider: string, planID?: string): string => {
  if (planID) {
    const planKey = savedProviderKeys.value[`claude:${provider}:${planID}`]
    if (planKey) return planKey
  }
  return savedProviderKeys.value[`claude:${provider}`] || ''
}

const canApplyAgent = (platform: AgentPlatform, serviceRunning: boolean) => {
  if (configLoading.value) return false
  if (platform === 'codex') {
    // 切换到 OpenAI 时始终允许，后端会尝试使用 auth.json 中现有的 key 或保存的 key
    if (selectedCodexProvider.value === 'openai') {
      return true
    }
    return true
  }
  if (selectedClaudeProvider.value === 'ccx') return true
  const provider = selectedClaudeProvider.value
  const inputKey = claudeProviderKeys.value[provider].trim()
  const hasSaved = !!findSavedKey(provider, selectedMiMoPlan.value)
  return inputKey !== '' || hasSaved
}

const applyAgent = async (platform: AgentPlatform) => {
  configLoading.value = true
  try {
    const request: ApplyAgentConfigRequest = { platform }
    if (platform === 'claude') {
      request.provider = selectedClaudeProvider.value
      if (selectedClaudeProvider.value !== 'ccx') {
        const inputKey = claudeProviderKeys.value[selectedClaudeProvider.value].trim()
        request.apiKey = inputKey || findSavedKey(selectedClaudeProvider.value, selectedMiMoPlan.value)
      }
      if (selectedClaudeProvider.value === 'mimo') {
        request.baseUrl = claudeMiMoBaseUrl.value.trim()
      }
    }
    if (platform === 'codex') {
      request.provider = selectedCodexProvider.value
      if (selectedCodexProvider.value === 'openai') {
        request.apiKey = codexOpenAIKey.value.trim() || savedProviderKeys.value['codex:openai'] || ''
      }
    }
    await ApplyAgentConfig(request)
    await loadAgentStatuses()
  } finally {
    configLoading.value = false
  }
}

const restoreAgent = async (platform: AgentPlatform) => {
  configLoading.value = true
  try {
    await RestoreAgentConfig(platform)
  } finally {
    configLoading.value = false
  }
}

export function useAgentConfig() {
  return {
    agentStatuses,
    configLoading,
    selectedClaudeProvider,
    claudeProviderKeys,
    savedProviderKeys,
    codexOpenAIKey,
    claudeMiMoBaseUrl,
    selectedMiMoPlan,
    agentLabels,
    claudeProviderLabels,
    codexProviderLabels,
    agentPlatforms,
    isClaudeProvider,
    claudeProviderLabel,
    claudeTargetBaseUrl,
    agentStatusText,
    agentStatusClass,
    loadAgentStatuses,
    canApplyAgent,
    applyAgent,
    restoreAgent,
    selectedCodexProvider,
    codexProviderLabel,
  }
}
