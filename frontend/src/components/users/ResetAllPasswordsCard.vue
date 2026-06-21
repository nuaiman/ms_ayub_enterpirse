<template>
  <div>
    <h2>Reset All Passwords</h2>

    <form @submit.prevent="submit">
      <input v-model="password" type="password" placeholder="New password" />

      <button type="submit">
        Reset All Passwords
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { useUsersStore } from '@/stores/users'

const usersStore = useUsersStore()

const password = ref('')

const submit = async () => {
  const success =
    await usersStore.resetAllPasswords(
      password.value,
    )

  if (!success) return

  password.value = ''
}
</script>
