<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import type { AgentProvider } from '@/types'

const props = defineProps<{
  selectedProvider: AgentProvider
  providerKeys: Record<AgentProvider, string>
  savedProviderKeys: Record<string, string>
  miMOBaseUrl: string
  selectedMiMoPlan: string
}>()

const emit = defineEmits<{
  'update:selectedProvider': [value: AgentProvider]
  'update:providerKeys': [value: Record<AgentProvider, string>]
  'update:miMOBaseUrl': [value: string]
  'update:selectedMiMoPlan': [value: string]
}>()

const mimoPlanOptions = [
  { label: '按量计费（默认）', value: 'https://api.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 中国', value: 'https://token-plan-cn.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 新加坡', value: 'https://token-plan-sgp.xiaomimimo.com/anthropic' },
  { label: '订阅套餐 - 欧洲', value: 'https://token-plan-ams.xiaomimimo.com/anthropic' },
  { label: '自定义', value: '' },
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

const onMiMoPlanChange = (e: Event) => {
  const planValue = (e.target as HTMLSelectElement).value
  emit('update:selectedMiMoPlan', planValue)
  if (planValue !== '') {
    emit('update:miMOBaseUrl', planValue)
  }
}

const keyPlaceholder = (provider: AgentProvider) => {
  if (provider === 'mimo' && props.selectedMiMoPlan && props.savedProviderKeys[`claude:${provider}:${props.selectedMiMoPlan}`]) {
    return '已保存，留空则使用已保存的 key'
  }
  if (props.savedProviderKeys[`claude:${provider}`]) {
    return '已保存，留空则使用已保存的 key'
  }
  if (provider === 'mimo') return 'MiMo API Key（tp-xxx 或账号 key）'
  return '仅写入 Claude Code 配置'
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
        <option value="kimi">Kimi 直连</option>
        <option value="glm">GLM 直连</option>
        <option value="minimax">MiniMax 直连</option>
        <option value="dashscope">DashScope 直连</option>
      </select>
    </div>

    <div v-if="selectedProvider === 'mimo'" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">MiMo 计费模式</Label>
      <select
        :value="selectedMiMoPlan"
        class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
        @change="onMiMoPlanChange"
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
      <Label class="text-xs text-muted-foreground">API Key</Label>
      <Input
        type="password"
        autocomplete="off"
        :placeholder="keyPlaceholder(selectedProvider)"
        :model-value="providerKeys[selectedProvider]"
        @update:model-value="onKeyChange"
      />
    </div>

    <div v-if="selectedProvider === 'mimo' && selectedMiMoPlan === ''" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">Base URL</Label>
      <Input
        type="url"
        placeholder="https://api.xiaomimimo.com/anthropic"
        :model-value="miMOBaseUrl"
        @update:model-value="emit('update:miMOBaseUrl', String($event))"
      />
    </div>
  </div>
</template>
