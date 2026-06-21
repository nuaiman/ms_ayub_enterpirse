<template>
  <div class="min-h-screen bg-app flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <!-- Logo / Brand -->
      <div class="text-center mb-8">
        <div
          class="w-16 h-16 mx-auto rounded-2xl bg-accent flex items-center justify-center mb-4 shadow-lg shadow-accent/20 overflow-hidden">
          <img src="@/assets/msae.png" alt="MS Ayub Enterprise" class="w-12 h-12 object-contain" />
        </div>
        <h1 class="text-2xl font-bold text-primary">MS Ayub Enterprise</h1>
        <p class="text-sm text-muted mt-1">Sign in to your account</p>
      </div>

      <!-- Login Card -->
      <div class="card rounded-2xl shadow-sm p-6 space-y-5">

        <!-- Username -->
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-primary">Username</label>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            <input v-model="username" type="text" placeholder="Enter your username" @keyup.enter="login" class="input w-full pl-9 pr-3 py-2.5 rounded-lg placeholder:text-muted
                     focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent
                     transition-all duration-200" />
          </div>
        </div>

        <!-- Password -->
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-primary">Password</label>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            <input v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="Enter your password"
              @keyup.enter="login" class="input w-full pl-9 pr-10 py-2.5 rounded-lg placeholder:text-muted
                     focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent
                     transition-all duration-200" />
            <button type="button" @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-primary transition-colors">
              <svg v-if="!showPassword" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="flex items-center gap-2 p-3 rounded-lg bg-warning-bg border border-warning-border">
          <svg class="w-4 h-4 text-warning-text shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <p class="text-sm text-warning-text">{{ error }}</p>
        </div>

        <!-- Submit -->
        <button @click="login" :disabled="loading || !username || !password" class="w-full py-2.5 px-4 bg-accent text-accent-foreground rounded-lg text-sm font-semibold
                 shadow-sm hover:shadow-md transition-all duration-200 active:scale-[0.98]
                 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100
                 inline-flex items-center justify-center gap-2">
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
            </path>
          </svg>
          <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
        </button>
      </div>

      <!-- Footer -->
      <p class="text-center text-xs text-muted mt-6">
        &copy; {{ new Date().getFullYear() }} MS Ayub Enterprise. All rights reserved.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const login = async () => {
  if (!username.value || !password.value) return

  loading.value = true
  error.value = ''

  try {
    const success = await auth.login(username.value, password.value)
    if (!success) {
      error.value = 'Invalid username or password'
    }
  } catch {
    error.value = 'An error occurred. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
