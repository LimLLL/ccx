export type HealthInfo = {
  status?: string
  timestamp?: string
  uptime?: number
  mode?: string
  version?: {
    version?: string
    buildTime?: string
    gitCommit?: string
  }
  config?: {
    upstreamCount?: number
  }
}

export type DesktopStatus = {
  running: boolean
  starting: boolean
  attached?: boolean
  port: number
  url: string
  pid: number
  binaryPath: string
  dataDir: string
  health?: HealthInfo
  lastError?: string
  logs: string[]
}

export type AgentPlatform = 'claude' | 'codex'
export type AgentProvider = 'ccx' | 'deepseek' | 'mimo' | 'kimi' | 'glm' | 'minimax' | 'dashscope' | 'openai'

export type AgentConfigStatus = {
  platform: AgentPlatform
  provider?: string
  targetProvider?: string
  configured: boolean
  matchesCurrentPort: boolean
  needsUpdate: boolean
  currentBaseUrl: string
  targetBaseUrl: string
  configPath: string
  authPath?: string
  hasState: boolean
  lastError?: string
}

export type ApplyAgentConfigRequest = {
  platform: AgentPlatform
  provider?: AgentProvider
  apiKey?: string
  baseUrl?: string
}

export type TabValue = 'status' | 'agent' | 'env' | 'channels' | 'web'

export type ProviderPlan = {
  id: string
  label: string
  baseUrl: string
  description: string
  recommended: boolean
  custom: boolean
}

export type ChannelTarget = {
  type: 'messages' | 'chat' | 'responses'
  label: string
  description: string
  recommended: boolean
}

export type ProviderPreset = {
  id: string
  label: string
  description: string
  directAgent: boolean
  nativeMessages: boolean
  chatCompatible: boolean
  responsesCompatible: boolean
  plans: ProviderPlan[]
  targets: ChannelTarget[]
  defaultTarget: 'messages' | 'chat' | 'responses'
}

export type ProviderKeyAsset = {
  provider: string
  apiKey: string
  baseUrl?: string
  planId?: string
  usages?: string[]
}

export type CreateChannelRequest = {
  provider: string
  target: string
  planId?: string
  baseUrl?: string
  apiKey?: string
  name?: string
  description?: string
}

export type CreateChannelResult = {
  provider: string
  target: string
  name: string
  baseUrl: string
  message: string
}
