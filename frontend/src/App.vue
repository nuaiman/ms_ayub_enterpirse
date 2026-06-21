<template>
  <div class="h-screen bg-app text-primary flex flex-col overflow-hidden">
    <RouterView class="flex-1 min-h-0" />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>
