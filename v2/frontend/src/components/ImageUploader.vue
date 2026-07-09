<template>
  <div>
    <input type="file" @change="onFileChange" />
    <button @click="upload" :disabled="!file">Upload</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  onUpload: (file: File) => Promise<void>
}>()

const file = ref<File | null>(null)

const onFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  file.value = target.files?.[0] || null
}

const upload = async () => {
  if (!file.value) return
  await props.onUpload(file.value)
  file.value = null
}
</script>
