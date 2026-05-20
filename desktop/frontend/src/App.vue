<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from '@/components/layout/Sidebar.vue'
import StatusTab from '@/components/status/StatusTab.vue'
import AgentTab from '@/components/agent/AgentTab.vue'
import EnvTab from '@/components/env/EnvTab.vue'
import WebUITab from '@/components/webui/WebUITab.vue'
import UpdateDialog from '@/components/update/UpdateDialog.vue'
import { useStatus } from '@/composables/useStatus'
import { useUpdater } from '@/composables/useUpdater'
import { useWailsEvents } from '@/composables/useWailsEvents'

const activeTab = ref<'status' | 'agent' | 'env' | 'web'>('status')
const { status, actionError, syncStatus } = useStatus()
useUpdater()

useWailsEvents(activeTab, actionError, syncStatus)

const switchToWeb = () => {
  activeTab.value = 'web'
}

// 选项卡标题映射
const tabTitles: Record<typeof activeTab.value, string> = {
  status: '网关状态监控',
  agent: 'Agent 代理配置',
  env: '环境参数管理',
  web: '内置控制台 Web UI'
}
</script>

<template>
  <div class="flex h-screen w-screen bg-[#060a13] text-slate-100 overflow-hidden font-sans">
    <!-- 常驻左侧高级磨砂侧边栏 -->
    <Sidebar v-model="activeTab" />

    <!-- 右侧内容主展区 -->
    <main class="flex-1 flex flex-col min-w-0 h-full relative">
      <!-- 右侧顶部精细页眉 -->
      <header class="h-14 border-b border-slate-900/60 bg-slate-950/25 backdrop-blur-md flex items-center justify-between px-8 shrink-0" data-wails-drag>
        <div class="flex items-center gap-3">
          <span class="text-xs bg-blue-500/10 text-blue-400 font-semibold px-2 py-0.5 rounded border border-blue-500/15">
            CCX CORE
          </span>
          <h2 class="text-sm font-bold text-slate-200 tracking-wide uppercase">
            {{ tabTitles[activeTab] }}
          </h2>
        </div>

        <div class="flex items-center gap-2">
          <!-- 实时网关状态指示微标 -->
          <span
            v-if="status.running"
            class="text-[10px] bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-bold px-2 py-0.5 rounded-full"
          >
            GATEWAY ONLINE
          </span>
          <span
            v-else-if="status.starting"
            class="text-[10px] bg-amber-500/10 text-amber-400 border border-amber-500/20 font-bold px-2 py-0.5 rounded-full animate-pulse"
          >
            CONNECTING...
          </span>
          <span
            v-else
            class="text-[10px] bg-rose-500/10 text-rose-400 border border-rose-500/20 font-bold px-2 py-0.5 rounded-full"
          >
            GATEWAY OFFLINE
          </span>
        </div>
      </header>

      <!-- 独立内容滚动区域 -->
      <div class="flex-1 overflow-y-auto px-8 py-7">
        <div class="max-w-5xl mx-auto h-full">
          <!-- 保持选项卡的过渡结构，不使用 TabSwitcher 从而完全摆脱药丸 Tabs，改用 v-show 优化性能和保持状态 -->
          <div v-show="activeTab === 'status'" class="h-full">
            <StatusTab @switch-to-web="switchToWeb" />
          </div>
          <div v-show="activeTab === 'agent'" class="h-full">
            <AgentTab />
          </div>
          <div v-show="activeTab === 'env'" class="h-full">
            <EnvTab />
          </div>
          <div v-show="activeTab === 'web'" class="h-full">
            <WebUITab :status="status" :loading="false" />
          </div>
        </div>
      </div>
    </main>

    <UpdateDialog />
  </div>
</template>
