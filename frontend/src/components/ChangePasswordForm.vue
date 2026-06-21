<template>
  <div class="relative group">
    <!-- Trigger -->
    <button class="button hover-surface px-3 py-1">
      Password
    </button>

    <!-- Hover Panel -->
    <div class="absolute right-0 mt-2 w-80 card p-4 space-y-3
             opacity-0 invisible group-hover:opacity-100 group-hover:visible
             transition-all duration-150 z-50">
      <h2 class="text-primary font-semibold text-sm">
        Change Password
      </h2>

      <input v-model="currentPassword" type="password" placeholder="Current password" class="input w-full" />

      <input v-model="newPassword" type="password" placeholder="New password" class="input w-full" />

      <input v-model="confirmPassword" type="password" placeholder="Confirm new password" class="input w-full" />

      <p v-if="error" class="text-sm text-warning-text">
        {{ error }}
      </p>

      <button class="button w-full hover-surface" @click="submit">
        Update
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const error = ref('')

const submit = async () => {
  error.value = ''

  if (!currentPassword.value || !newPassword.value) {
    error.value = 'All fields are required'
    return
  }

  if (newPassword.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  const ok = await auth.changePassword(
    currentPassword.value,
    newPassword.value
  )

  if (!ok) {
    error.value = 'Failed to update password'
    return
  }

  currentPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
}
</script>
