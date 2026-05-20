<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { useUpdater } from '@/composables/useUpdater'

const { state, downloadAndInstall, cancel, closeDialog } = useUpdater()

const formatBytes = (n: number): string => {
  if (!n) return '—'
  const units = ['B', 'KB', 'MB', 'GB']
  let v = n
  let i = 0
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(1)} ${units[i]}`
}

const phaseLabel = computed(() => {
  switch (state.phase) {
    case 'checking':
      return '检查中…'
    case 'downloading':
      return `下载中 · ${formatBytes(state.downloaded)} / ${formatBytes(state.total)}`
    case 'verifying':
      return '校验中…'
    case 'installing':
      return '准备安装…'
    case 'done':
      return '准备就绪'
    case 'error':
      return state.error || '出错'
    default:
      return ''
  }
})

const close = () => {
  if (!state.downloading) closeDialog()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="state.dialogOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="close"
      >
        <div
          class="w-[min(560px,90vw)] overflow-hidden rounded-2xl border border-white/[0.06] bg-[#0f1421] shadow-2xl"
        >
          <div class="border-b border-white/[0.05] px-6 py-4">
            <div class="flex items-baseline justify-between">
              <h2 class="text-lg font-semibold text-white">发现新版本</h2>
              <div class="text-xs text-white/40">
                v{{ state.version?.version || '—' }} → v{{ state.info?.latestVersion || '—' }}
              </div>
            </div>
          </div>

          <div class="max-h-[40vh] overflow-y-auto px-6 py-5">
            <div v-if="state.info?.notes" class="space-y-2 text-sm leading-relaxed text-white/75">
              <pre class="whitespace-pre-wrap break-words font-sans">{{ state.info.notes }}</pre>
            </div>
            <div v-else class="text-sm text-white/40">本次更新无发布说明</div>
          </div>

          <div v-if="state.downloading || state.phase === 'done' || state.phase === 'error'" class="border-t border-white/[0.05] px-6 py-4">
            <div class="mb-2 flex items-center justify-between text-xs">
              <span class="text-white/60">{{ phaseLabel }}</span>
              <span v-if="state.phase === 'downloading'" class="font-mono text-white/80">
                {{ state.percent.toFixed(1) }}%
              </span>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-white/[0.06]">
              <div
                class="h-full rounded-full bg-emerald-500 transition-all duration-200"
                :class="{ 'bg-red-500': state.phase === 'error' }"
                :style="{ width: `${state.phase === 'error' || state.phase === 'done' ? 100 : state.percent}%` }"
              />
            </div>
            <div v-if="state.error" class="mt-2 text-xs text-red-400">{{ state.error }}</div>
          </div>

          <div class="flex justify-end gap-2 border-t border-white/[0.05] px-6 py-4">
            <Button
              v-if="state.downloading"
              variant="outline"
              size="sm"
              @click="cancel"
            >
              取消
            </Button>
            <template v-else>
              <Button variant="ghost" size="sm" @click="close">
                稍后再说
              </Button>
              <Button size="sm" @click="downloadAndInstall">
                {{ state.phase === 'error' ? '重试' : '立即更新' }}
              </Button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
