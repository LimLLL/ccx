<script setup lang="ts">
import { computed } from 'vue'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import ProviderForm from '@/components/agent/ProviderForm.vue'
import type { AgentPlatform, AgentConfigStatus, AgentProvider } from '@/types'

const props = defineProps<{
  platform: AgentPlatform
  agentStatus: AgentConfigStatus | null
  configLoading: boolean
  serviceRunning: boolean
  agentLabel: string
  agentStatusText: string
  agentStatusClass: string
  canApply: boolean
  selectedClaudeProvider?: AgentProvider
  claudeProviderKeys?: Record<AgentProvider, string>
  savedProviderKeys?: Record<string, string>
  claudeMiMoBaseUrl?: string
  selectedMiMoPlan?: string
  claudeProviderLabel?: (value?: string) => string
  claudeTargetBaseUrl?: () => string
  selectedCodexProvider?: AgentProvider
  codexOpenAIKey?: string
  codexProviderLabels?: Record<string, string>
  codexProviderLabel?: (value?: string) => string
}>()

const emit = defineEmits<{
  apply: []
  restore: []
  'update:selectedClaudeProvider': [value: AgentProvider]
  'update:claudeProviderKeys': [value: Record<AgentProvider, string>]
  'update:claudeMiMoBaseUrl': [value: string]
  'update:selectedMiMoPlan': [value: string]
  'update:selectedCodexProvider': [value: AgentProvider]
  'update:codexOpenAIKey': [value: string]
}>()

const badgeClass = computed(() => {
  if (props.agentStatusClass === 'running') return 'bg-accent text-accent-foreground border-0'
  if (props.agentStatusClass === 'starting') return 'bg-warning text-warning-foreground border-0'
  return 'bg-destructive text-destructive-foreground border-0'
})

const applyLabel = computed(() => '应用配置')
</script>

<template>
  <Card>
    <CardHeader class="pb-3">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-xs text-muted-foreground uppercase tracking-widest">{{ platform }}</p>
          <h3 class="text-base font-semibold mt-1">{{ agentLabel }}</h3>
        </div>
        <Badge :class="badgeClass">{{ agentStatusText }}</Badge>
      </div>
    </CardHeader>

    <CardContent class="space-y-4">
      <div class="space-y-2 text-sm">
        <div v-for="detail in [
          { label: '当前 Provider', value: (platform === 'codex' ? codexProviderLabel : claudeProviderLabel)?.(agentStatus?.provider) || agentStatus?.provider || '未设置' },
          { label: '当前 URL', value: agentStatus?.currentBaseUrl || '未设置' },
          { label: '目标 URL', value: platform === 'claude' ? claudeTargetBaseUrl?.() : (platform === 'codex' && selectedCodexProvider === 'openai' ? '-' : (agentStatus?.targetBaseUrl || '-')) },
          { label: '配置文件', value: agentStatus?.configPath || '--' },
        ]" :key="detail.label" class="flex items-center justify-between">
          <span class="text-muted-foreground">{{ detail.label }}</span>
          <code class="text-xs bg-secondary px-2 py-0.5 rounded break-all max-w-[60%] text-right">{{ detail.value }}</code>
        </div>
        <div v-if="agentStatus?.authPath" class="flex items-center justify-between">
          <span class="text-muted-foreground">认证文件</span>
          <code class="text-xs bg-secondary px-2 py-0.5 rounded break-all max-w-[60%] text-right">{{ agentStatus.authPath }}</code>
        </div>
      </div>

      <ProviderForm
        v-if="platform === 'claude'"
        :selected-provider="selectedClaudeProvider!"
        :provider-keys="claudeProviderKeys!"
        :saved-provider-keys="savedProviderKeys || {}"
        :mi-m-o-base-url="claudeMiMoBaseUrl!"
        :selected-mi-mo-plan="selectedMiMoPlan || 'https://api.xiaomimimo.com/anthropic'"
        @update:selected-provider="emit('update:selectedClaudeProvider', $event)"
        @update:provider-keys="emit('update:claudeProviderKeys', $event)"
        @update:mi-m-o-base-url="emit('update:claudeMiMoBaseUrl', $event)"
        @update:selected-mi-mo-plan="emit('update:selectedMiMoPlan', $event)"
      />
      <div v-else-if="platform === 'codex'" class="space-y-3">
        <div class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">Provider</Label>
          <select
            :value="selectedCodexProvider"
            class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            @change="emit('update:selectedCodexProvider', ($event.target as HTMLSelectElement).value as AgentProvider)"
          >
            <option value="ccx">CCX 本地网关</option>
            <option value="openai">OpenAI 官方</option>
          </select>
        </div>
        <div v-if="selectedCodexProvider === 'openai'" class="space-y-1.5">
          <Label class="text-xs text-muted-foreground">OpenAI API Key</Label>
          <Input
            type="password"
            autocomplete="off"
            :placeholder="savedProviderKeys?.['codex:openai'] ? '已保存，留空则使用已保存的 key' : '仅写入 Codex 官方配置'"
            :model-value="codexOpenAIKey || ''"
            @update:model-value="emit('update:codexOpenAIKey', String($event))"
          />
        </div>
      </div>

      <p v-if="agentStatus?.lastError" class="text-sm text-destructive-foreground">{{ agentStatus.lastError }}</p>

      <div class="flex flex-wrap gap-2">
        <Button size="sm" :disabled="!canApply" @click="emit('apply')">
          {{ applyLabel }}
        </Button>
        <Button size="sm" variant="secondary" :disabled="configLoading || !agentStatus?.hasState" @click="emit('restore')">
          恢复原始配置
        </Button>
      </div>
    </CardContent>
  </Card>
</template>
