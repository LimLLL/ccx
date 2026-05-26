<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { KeyRound, Network, Sparkles, CheckCircle2, Loader2, ExternalLink } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useChannelPresets } from '@/composables/useChannelPresets'
import { openProviderPromotion, openProviderConsole, providerConsoleLinks, providerPromotionLinks } from '@/lib/external-link'
import compshareIcon from '@/assets/compshare.png'
import type { ProviderPreset, ChannelTarget } from '@/types'

const {
  presets,
  keysByProvider,
  loading,
  creating,
  error,
  result,
  loadChannelPresets,
  createChannel,
} = useChannelPresets()

const providerIcons: Record<string, string> = {
  compshare: compshareIcon,
}

const selectedProvider = ref('')
const selectedTarget = ref('')
const selectedPlan = ref('')
const apiKey = ref('')
const channelName = ref('')
const localError = ref('')

onMounted(async () => {
  await loadChannelPresets()
  if (!selectedProvider.value && presets.value.length > 0) {
    selectedProvider.value = presets.value[0].id
  }
})

const currentPreset = computed(() => {
  return presets.value.find((item) => item.id === selectedProvider.value) || null
})

const currentAsset = computed(() => {
  return selectedProvider.value ? keysByProvider.value[selectedProvider.value] : undefined
})

const currentPlan = computed(() => {
  return currentPreset.value?.plans.find((item) => item.id === selectedPlan.value) || currentPreset.value?.plans[0]
})

const targetOptions = computed<ChannelTarget[]>(() => currentPreset.value?.targets || [])

// 切换左侧 provider 时重置表单；仅响应 selectedProvider 变化，
// 避免 target 变化触发 loadChannelPresets 后因 presets 更新而级联重置 selectedTarget
watch(selectedProvider, (id) => {
  const preset = presets.value.find((item) => item.id === id)
  if (!preset) return
  selectedTarget.value = preset.defaultTarget
  selectedPlan.value = bestPlanForTarget(preset, preset.defaultTarget)
  apiKey.value = ''
  channelName.value = `desktop-${preset.id}-${preset.defaultTarget}`
})

// target 变化时重新加载后端过滤后的 plans，尽量保留已选 plan；
// 若当前 plan 被协议过滤掉，尝试切换同区域的协议变体（如 token-cn ↔ token-cn-anthropic）
watch(selectedTarget, async (target) => {
  if (!target) return
  const prevPlan = selectedPlan.value
  await loadChannelPresets(target)
  const preset = currentPreset.value
  if (!preset) return
  if (preset.plans.some((p) => p.id === prevPlan)) {
    selectedPlan.value = prevPlan
  } else {
    const counterpart = prevPlan.endsWith('-anthropic')
      ? prevPlan.replace(/-anthropic$/, '')
      : prevPlan + '-anthropic'
    const match = preset.plans.find((p) => p.id === counterpart)
    selectedPlan.value = match ? match.id : bestPlanForTarget(preset, target)
  }
  channelName.value = `desktop-${preset.id}-${target}`
})

function bestPlanForTarget(preset: ProviderPreset, target: string): string {
  if (preset.plans.length <= 1) return preset.plans[0]?.id || ''
  const wantAnthropic = target === 'messages'
  for (const plan of preset.plans) {
    if (plan.custom) continue
    const isAnthropic = plan.baseUrl?.includes('anthropic')
    if (wantAnthropic && isAnthropic) return plan.id
    if (!wantAnthropic && !isAnthropic) return plan.id
  }
  const recommended = preset.plans.find((p) => p.recommended)
  return recommended?.id || preset.plans[0]?.id || ''
}

const capabilityBadges = computed(() => {
  const preset = currentPreset.value
  if (!preset) return []
  return [
    preset.directAgent && 'Agent 直连',
    preset.nativeMessages && 'Messages 原生',
    preset.chatCompatible && 'OpenAI Chat',
    preset.responsesCompatible && 'Responses',
  ].filter(Boolean) as string[]
})

const effectiveBaseUrl = computed(() => {
  return currentPlan.value?.baseUrl || currentAsset.value?.baseUrl || ''
})

const keyPlaceholder = computed(() => {
  return currentAsset.value?.apiKey ? '已保存，留空则复用该 Provider Key' : '输入 API Key，仅保存在本机 Desktop 配置中'
})

const submit = async () => {
  localError.value = ''
  const preset = currentPreset.value
  if (!preset) return
  if (!apiKey.value.trim() && !currentAsset.value?.apiKey) {
    localError.value = '请填写 API Key，或先在 Agent 配置中保存该 Provider 的 key。'
    return
  }
  await createChannel({
    provider: preset.id,
    target: selectedTarget.value || preset.defaultTarget,
    planId: selectedPlan.value,
    baseUrl: effectiveBaseUrl.value,
    apiKey: apiKey.value.trim() || currentAsset.value?.apiKey || '',
    name: channelName.value.trim(),
  })
  apiKey.value = ''
}
</script>

<template>
  <div class="space-y-5">
    <div class="bg-glass border border-white/[0.03] rounded-2xl p-5">
      <div class="flex items-start justify-between gap-4">
        <div>
          <div class="flex items-center gap-2 text-blue-400 mb-2">
            <Network class="w-4 h-4" />
            <span class="text-xs font-bold uppercase tracking-[0.2em]">Channel Preset Center</span>
          </div>
          <h3 class="text-xl font-bold text-slate-100">渠道中心</h3>
          <p class="text-sm text-slate-500 mt-1 max-w-2xl">
            统一把 DeepSeek、MiMo、Kimi、GLM、MiniMax Key 可同时用于 Agent 直连和 CCX 统一渠道池，复杂开关由预设自动处理。
          </p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-[280px_1fr] gap-4">
      <div class="space-y-2">
        <button
          v-for="preset in presets"
          :key="preset.id"
          :class="[
            'w-full p-4 rounded-xl border text-left transition-all duration-200 bg-glass-hover',
            selectedProvider === preset.id
              ? 'border-blue-500/30 bg-blue-500/10 shadow-[0_0_18px_rgba(59,130,246,0.12)]'
              : 'border-white/[0.03] bg-white/[0.01] hover:border-slate-700'
          ]"
          @click="selectedProvider = preset.id"
        >
          <div class="flex items-start gap-3">
            <img
              v-if="providerIcons[preset.id]"
              :src="providerIcons[preset.id]"
              :alt="`${preset.label} icon`"
              class="mt-0.5 h-9 w-9 shrink-0 rounded-xl bg-slate-900/60 object-cover ring-1 ring-white/[0.04]"
            >
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-2">
                <span class="font-semibold text-slate-200">{{ preset.label }}</span>
                <span v-if="keysByProvider[preset.id]" class="text-[10px] text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20">
                  已有 Key
                </span>
              </div>
              <p class="text-xs text-slate-500 mt-1 line-clamp-2">{{ preset.description }}</p>
            </div>
          </div>
        </button>
      </div>

      <div v-if="currentPreset" class="bg-glass border border-white/[0.03] rounded-2xl p-5 space-y-5">
        <div class="space-y-3">
          <div class="flex flex-wrap items-center gap-2">
            <h3 class="text-lg font-semibold text-slate-100">{{ currentPreset.label }}</h3>
            <span
              v-for="badge in capabilityBadges"
              :key="badge"
              class="text-[10px] text-blue-300 bg-blue-500/10 border border-blue-500/20 rounded-full px-2 py-0.5"
            >
              {{ badge }}
            </span>
          </div>
          <p class="text-sm text-slate-500">{{ currentPreset.description }}</p>
          <div class="flex items-center gap-4">
            <button
              v-if="providerPromotionLinks[currentPreset.id]"
              type="button"
              class="inline-flex items-center gap-1.5 text-xs font-medium text-blue-300 hover:text-blue-200"
              @click="openProviderPromotion(currentPreset.id)"
            >
              通过推广链接注册，领取 5 元平台试用金
              <ExternalLink class="h-3 w-3" />
            </button>
            <button
              v-if="providerConsoleLinks[currentPreset.id]"
              type="button"
              class="inline-flex items-center gap-1.5 text-xs font-medium text-slate-400 hover:text-slate-200"
              @click="openProviderConsole(currentPreset.id)"
            >
              访问官方控制台
              <ExternalLink class="h-3 w-3" />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label class="text-xs text-slate-400">添加目标</Label>
            <select
              v-model="selectedTarget"
              class="w-full h-9 rounded-md border border-slate-800 bg-slate-950/70 px-3 text-sm text-slate-200 focus:outline-none focus:ring-1 focus:ring-blue-500"
            >
              <option v-for="target in targetOptions" :key="target.type" :value="target.type">
                {{ target.label }}
              </option>
            </select>
            <p class="text-xs text-slate-500">
              {{ targetOptions.find((item) => item.type === selectedTarget)?.description }}
            </p>
          </div>

          <div class="space-y-2">
            <Label class="text-xs text-slate-400">Token Plan / Base URL</Label>
            <select
              v-model="selectedPlan"
              class="w-full h-9 rounded-md border border-slate-800 bg-slate-950/70 px-3 text-sm text-slate-200 focus:outline-none focus:ring-1 focus:ring-blue-500"
            >
              <option v-for="plan in currentPreset.plans" :key="plan.id" :value="plan.id">
                {{ plan.label }}
              </option>
            </select>
            <p class="text-xs text-slate-500">{{ currentPlan?.description }}</p>
            <p v-if="currentPlan?.baseUrl" class="text-xs text-slate-400 font-mono">{{ currentPlan.baseUrl }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div class="space-y-2">
            <Label class="text-xs text-slate-400">API Key</Label>
            <Input v-model="apiKey" type="password" autocomplete="off" :placeholder="keyPlaceholder" />
            <div v-if="currentAsset?.apiKey" class="flex items-center gap-1.5 text-xs text-emerald-400">
              <KeyRound class="w-3 h-3" />
              将复用本机已保存的 {{ currentPreset.label }} Key
            </div>
          </div>

          <div class="space-y-2">
            <Label class="text-xs text-slate-400">渠道名称</Label>
            <Input v-model="channelName" placeholder="desktop-provider-type" />
            <p class="text-xs text-slate-500">同名渠道会被直接覆盖更新；如需新建独立渠道，请改用不同名称。</p>
          </div>
        </div>

        <div class="rounded-xl bg-slate-950/50 border border-slate-900 p-3 text-xs text-slate-400 space-y-1.5">
          <div class="flex items-center gap-1.5 text-slate-300">
            <Sparkles class="w-3.5 h-3.5 text-blue-400" />
            <span class="font-semibold">预设将自动写入</span>
          </div>
          <p>Base URL: <code class="text-slate-200">{{ effectiveBaseUrl || '—' }}</code></p>
          <p>能力开关：reasoning / vision / model list / 兼容字段会按 Provider 自动配置。</p>
        </div>

        <p v-if="localError || error" class="text-sm text-rose-400">{{ localError || error }}</p>
        <div v-if="result" class="flex items-center gap-2 rounded-xl border border-emerald-500/20 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-300">
          <CheckCircle2 class="w-4 h-4" />
          {{ result.message }}：{{ result.name }} → {{ result.baseUrl }}
        </div>

        <div class="flex justify-end">
          <Button :disabled="creating" @click="submit">
            <Loader2 v-if="creating" class="w-3.5 h-3.5 mr-1.5 animate-spin" />
            添加到 CCX
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
