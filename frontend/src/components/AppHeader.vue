<script setup lang="ts">
import { Activity, History } from 'lucide-vue-next'
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()

const navItems = [
  { path: '/', label: '信号分析', icon: Activity },
  { path: '/history', label: '检测历史', icon: History },
]

const isActive = (path: string) => computed(() => route.path === path)
</script>

<template>
  <header class="bg-osc-panel border-b border-osc-border px-6 py-3 flex items-center justify-between sticky top-0 z-50">
    <div class="flex items-center gap-3">
      <div class="relative w-8 h-8 flex items-center justify-center">
        <div class="absolute inset-0 bg-osc-green/20 rounded-md animate-pulse-slow"></div>
        <Activity class="w-5 h-5 text-osc-green relative z-10" />
      </div>
      <div>
        <h1 class="text-osc-bright font-display font-bold text-sm tracking-wider">ULTRASONIC NDT</h1>
        <p class="text-osc-muted text-[10px] font-mono">SIGNAL ANALYZER v1.0</p>
      </div>
    </div>

    <nav class="flex items-center gap-1">
      <button
        v-for="item in navItems"
        :key="item.path"
        @click="router.push(item.path)"
        :class="[
          'flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-all duration-200',
          isActive(item.path).value
            ? 'bg-osc-green/10 text-osc-green border border-osc-green/30 shadow-[0_0_12px_rgba(0,230,118,0.15)]'
            : 'text-osc-text hover:text-osc-bright hover:bg-osc-card/50',
        ]"
      >
        <component :is="item.icon" class="w-4 h-4" />
        <span>{{ item.label }}</span>
      </button>
    </nav>

    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full bg-osc-green animate-pulse shadow-green-glow"></div>
      <span class="text-osc-muted text-xs font-mono">SYSTEM ONLINE</span>
    </div>
  </header>
</template>
