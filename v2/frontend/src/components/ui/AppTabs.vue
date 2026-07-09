<template>
  <div class="w-full">

    <!-- 📱 MOBILE: full width dropdown -->
    <div class="lg:hidden w-full">
      <div
        class="relative w-full flex items-center px-3 py-2.5 bg-input border border-default rounded-xl text-sm transition-colors">
        <!-- Icon -->
        <div class="shrink-0 mr-2 text-muted">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
        </div>

        <select
          class="w-full bg-transparent pr-8 py-1 text-sm text-primary focus:outline-none appearance-none cursor-pointer"
          :value="modelValue" @change="onChange">
          <option v-for="tab in visibleTabs" :key="tab.key" :value="tab.key" class="bg-surface text-primary">
            {{ tab.label }}
          </option>
        </select>

        <!-- Dropdown arrow -->
        <div class="absolute right-3 text-muted pointer-events-none">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>
    </div>

    <!-- 🖥️ DESKTOP: pill tabs -->
    <div class="hidden lg:flex justify-start">
      <div class="inline-flex items-center gap-1 border border-default rounded-xl p-1 bg-input">

        <template v-for="tab in tabs" :key="tab.key">
          <!-- Divider -->
          <div v-if="tab.divider" class="w-px h-5 bg-border mx-1 self-center shrink-0"></div>

          <!-- Tab button -->
          <button v-else @click="select(tab.key)"
            class="relative px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200" :class="modelValue === tab.key
              ? 'bg-surface text-primary shadow-sm'
              : 'text-muted hover:text-secondary hover:bg-surface-alt'">
            {{ tab.label }}
          </button>
        </template>

      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  tabs: { key: string; label: string; divider?: boolean }[]
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// Filter out dividers for mobile dropdown
const visibleTabs = computed(() => props.tabs.filter(t => !t.divider))

const select = (key: string) => {
  emit('update:modelValue', key)
}

const onChange = (e: Event) => {
  const value = (e.target as HTMLSelectElement).value
  emit('update:modelValue', value)
}
</script>
