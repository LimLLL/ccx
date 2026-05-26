<script setup lang="ts">
import { ExternalLink } from 'lucide-vue-next'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { openProviderPromotion, openProviderConsole, providerConsoleLinks, providerPromotionLinks } from '@/lib/external-link'
import type { AgentProvider } from '@/types'

const props = defineProps<{
  selectedProvider: AgentProvider
  providerKeys: Record<AgentProvider, string>
  savedProviderKeys: Record<string, string>
  mimoBaseUrl: string
  selectedMimoPlan: string
}>()

const emit = defineEmits<{
  'update:selectedProvider': [value: AgentProvider]
  'update:providerKeys': [value: Record<AgentProvider, string>]
  'update:mimoBaseUrl': [value: string]
  'update:selectedMimoPlan': [value: string]
}>()

const mimoPlanOptions = [
  { label: '按量', value: 'https://api.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 中国', value: 'https://token-plan-cn.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 新加坡', value: 'https://token-plan-sgp.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 欧洲', value: 'https://token-plan-ams.xiaomimimo.com/anthropic' },
]

const onProviderChange = (e: Event) => {
  emit('update:selectedProvider', (e.target as HTMLSelectElement).value as AgentProvider)
}

const onKeyChange = (value: string | number) => {
  emit('update:providerKeys', {
    ...props.providerKeys,
    [props.selectedProvider]: String(value),
  })
}

const onMimoPlanChange = (e: Event) => {
  const planValue = (e.target as HTMLSelectElement).value
  emit('update:selectedMimoPlan', planValue)
  emit('update:mimoBaseUrl', planValue)
}

const keyPlaceholder = (provider: AgentProvider) => {
  if (provider === 'mimo' && props.selectedMimoPlan && props.savedProviderKeys[`claude:${provider}:${props.selectedMimoPlan}`]) {
    return '已保存，留空则使用已保存的 key'
  }
  if (props.savedProviderKeys[`claude:${provider}`]) {
    return '已保存，留空则使用已保存的 key'
  }
  if (provider === 'mimo') return '必填：MiMo API Key（tp-xxx 或账号 key）'
  return '必填：输入 API Key'
}
</script>

<template>
  <div class="space-y-3">
    <div class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">Provider</Label>
      <select
        :value="selectedProvider"
        class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
        @change="onProviderChange"
      >
        <option value="ccx">CCX 本地网关</option>
        <option value="deepseek">DeepSeek 直连</option>
        <option value="mimo">MiMo 直连</option>
        <option value="compshare">Compshare 直连</option>
        <option value="kimi">Kimi 直连</option>
        <option value="glm">GLM 直连</option>
        <option value="minimax">MiniMax 直连</option>
        <option value="dashscope">DashScope 直连</option>
        <option value="opencode-zen">OpenCode Zen 直连</option>
        <option value="opencode-go">OpenCode Go 直连</option>
      </select>
    </div>

    <div v-if="selectedProvider !== 'ccx'" class="inline-flex items-center gap-3 flex-wrap">
      <button
        v-if="providerPromotionLinks[selectedProvider]"
        type="button"
        class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-300 hover:text-blue-200"
        @click="openProviderPromotion(selectedProvider)"
      >
        通过推广链接注册，领取 5 元平台试用金
        <ExternalLink class="h-3 w-3" />
      </button>
      <button
        v-if="providerConsoleLinks[selectedProvider]"
        type="button"
        class="inline-flex items-center gap-1.5 text-xs font-medium text-slate-400 hover:text-slate-200"
        @click="openProviderConsole(selectedProvider)"
      >
        访问官方控制台
        <ExternalLink class="h-3 w-3" />
      </button>
    </div>

    <div v-if="selectedProvider === 'mimo'" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">MiMo 计费模式</Label>
      <select
        :value="selectedMimoPlan"
        class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
        @change="onMimoPlanChange"
      >
        <option
          v-for="opt in mimoPlanOptions"
          :key="opt.value || '__custom__'"
          :value="opt.value"
        >
          {{ opt.label }}
        </option>
      </select>
    </div>

    <div v-if="selectedProvider !== 'ccx'" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">API Key <span class="text-destructive">*</span></Label>
      <Input
        type="password"
        autocomplete="off"
        :placeholder="keyPlaceholder(selectedProvider)"
        :model-value="providerKeys[selectedProvider]"
        @update:model-value="onKeyChange"
      />
    </div>

  </div>
</template>
