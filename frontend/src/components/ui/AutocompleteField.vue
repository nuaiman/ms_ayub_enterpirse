<template>
  <div class="space-y-1.5 relative">
    <label class="text-sm font-medium text-primary">
      {{ label }} <span v-if="required" class="text-warning-text">*</span>
    </label>
    <div class="relative">
      <input
        :value="modelValue"
        @input="handleInput"
        :type="type"
        :placeholder="placeholder"
        :required="required"
        @focus="showDropdown = true"
        class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent"
      />
      <button v-if="modelValue" @click="clearInput" type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
      </button>
    </div>
    <div v-if="showDropdown && filteredSuggestions.length > 0" class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
      <div class="max-h-48 overflow-y-auto">
        <button v-for="suggestion in filteredSuggestions" :key="suggestion" @click="selectSuggestion(suggestion)" type="button"
          class="w-full flex items-center px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
          {{ suggestion }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  modelValue: string
  label: string
  placeholder?: string
  required?: boolean
  type?: string
  suggestions: string[]
}>()

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const showDropdown = ref(false)

const filteredSuggestions = computed(() => {
  if (!props.modelValue) return props.suggestions
  const q = props.modelValue.toLowerCase()
  return props.suggestions.filter(s => s.toLowerCase().includes(q) && s.toLowerCase() !== q)
})

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
  showDropdown.value = true
}

const clearInput = () => {
  emit('update:modelValue', '')
  showDropdown.value = true
}

const selectSuggestion = (value: string) => {
  emit('update:modelValue', value)
  showDropdown.value = false
}
</script>