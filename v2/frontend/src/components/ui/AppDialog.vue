<template>
  <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">

    <!-- backdrop -->
    <div class="absolute inset-0 bg-black/40" @click="close" />

    <!-- dialog -->
    <div class="relative bg-surface p-6 rounded-lg w-full max-w-lg card">
      <slot />
    </div>

  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits(['update:modelValue'])

const close = () => {
  emit('update:modelValue', false)
}
</script>
