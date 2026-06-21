<template>
  <nav v-if="auth.isAuthenticated" class="bg-surface border-b border-divider sticky top-0 z-30">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">

        <!-- User Menu -->
        <div class="flex items-center gap-4">
          <div ref="dropdownRef" class="relative">
            <button @click="toggleDropdown"
              class="flex items-center gap-3 py-1.5 px-2 rounded-lg hover:bg-surface-alt transition-colors group">
              <!-- Avatar -->
              <div class="shrink-0">
                <div v-if="auth.user?.image_url" class="w-8 h-8 rounded-full overflow-hidden border border-default">
                  <img :src="getImageUrl(auth.user.image_url)" :alt="auth.user?.name"
                    class="w-full h-full object-cover" />
                </div>
                <div v-else class="w-8 h-8 rounded-full bg-accent flex items-center justify-center">
                  <span class="text-xs font-semibold text-accent-foreground">
                    {{ getInitials(auth.user?.name || '') }}
                  </span>
                </div>
              </div>

              <!-- Name & Role -->
              <div class="hidden sm:block text-left">
                <div class="text-sm font-medium text-primary leading-tight">
                  {{ auth.user?.name }}
                </div>
                <div class="text-xs text-muted capitalize">
                  {{ auth.user?.role }}
                </div>
              </div>

              <!-- Chevron -->
              <svg class="hidden sm:block w-4 h-4 text-muted transition-transform duration-200"
                :class="{ 'rotate-180': open }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            <!-- Dropdown Menu -->
            <Transition enter-active-class="transition ease-out duration-200"
              enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
              leave-active-class="transition ease-in duration-150"
              leave-from-class="opacity-100 scale-100 translate-y-0" leave-to-class="opacity-0 scale-95 translate-y-1">
              <div v-if="open"
                class="absolute right-0 mt-2 w-64 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50">
                <!-- User Info Header -->
                <div class="px-4 py-3 border-b border-divider">
                  <div class="flex items-center gap-3">
                    <div class="shrink-0">
                      <div v-if="auth.user?.image_url"
                        class="w-10 h-10 rounded-full overflow-hidden border border-default">
                        <img :src="getImageUrl(auth.user.image_url)" :alt="auth.user?.name"
                          class="w-full h-full object-cover" />
                      </div>
                      <div v-else class="w-10 h-10 rounded-full bg-accent flex items-center justify-center">
                        <span class="text-sm font-semibold text-accent-foreground">
                          {{ getInitials(auth.user?.name || '') }}
                        </span>
                      </div>
                    </div>
                    <div class="min-w-0">
                      <div class="text-sm font-semibold text-primary truncate">
                        {{ auth.user?.name }}
                      </div>
                      <div class="text-xs text-muted">
                        @{{ auth.user?.username }}
                      </div>
                      <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize mt-1"
                        :class="getRoleBadgeClass(auth.user?.role || '')">
                        {{ auth.user?.role }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Menu Items -->
                <div class="py-1">
                  <!-- Change Password -->
                  <button @click="openChangePassword"
                    class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors group/menu">
                    <svg class="w-4 h-4 shrink-0 group-hover/menu:text-primary transition-colors" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                    </svg>
                    Change Password
                  </button>

                  <!-- Theme Toggle -->
                  <button @click="toggleTheme"
                    class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors group/menu">
                    <svg v-if="theme === 'dark'"
                      class="w-4 h-4 shrink-0 group-hover/menu:text-primary transition-colors" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    <svg v-else class="w-4 h-4 shrink-0 group-hover/menu:text-primary transition-colors" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                    </svg>
                    {{ themeLabel }}
                  </button>

                  <!-- Admin Only Section -->
                  <template v-if="auth.user?.role === 'admin'">
                    <div class="border-t border-divider my-1"></div>

                    <!-- Download Backup -->
                    <button @click="downloadBackup"
                      class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors group/menu">
                      <svg class="w-4 h-4 shrink-0 group-hover/menu:text-primary transition-colors" fill="none"
                        stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                      </svg>
                      Download Backup
                    </button>

                    <!-- Reset All Passwords -->
                    <button @click="openResetAll"
                      class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors group/menu">
                      <svg class="w-4 h-4 shrink-0 group-hover/menu:text-primary transition-colors" fill="none"
                        stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                      </svg>
                      Reset All Passwords
                    </button>
                  </template>

                  <!-- Divider -->
                  <div class="border-t border-divider my-1"></div>

                  <!-- Logout -->
                  <button @click="handleLogout"
                    class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-warning-text hover:bg-surface-alt transition-colors group/menu">
                    <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Sign Out
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </div>

    <!-- CHANGE PASSWORD DIALOG -->
    <BaseDialog v-model="showChangePassword">
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-info-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-info-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Change Password</h2>
            <p class="text-xs text-muted">Update your account password</p>
          </div>
        </div>

        <div class="space-y-3">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Current Password</label>
            <input v-model="changePasswordForm.current" type="password" placeholder="Enter current password"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">New Password</label>
            <input v-model="changePasswordForm.new" type="password" placeholder="Enter new password"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent"
              @keyup.enter="submitChangePassword" />
          </div>
        </div>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showChangePassword = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="submitChangePassword"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Update Password</button>
        </div>
      </div>
    </BaseDialog>

    <!-- RESET ALL PASSWORDS DIALOG -->
    <BaseDialog v-model="showResetAll">
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-warning-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Reset All Passwords</h2>
            <p class="text-xs text-muted">This will reset passwords for all users except root</p>
          </div>
        </div>

        <div class="space-y-1.5">
          <label class="text-sm font-medium text-primary">New Password</label>
          <input v-model="resetAllForm.new" type="password" placeholder="Enter new password for all users"
            class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent"
            @keyup.enter="submitResetAll" />
        </div>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showResetAll = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="submitResetAll"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Reset All</button>
        </div>
      </div>
    </BaseDialog>

  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { getImageUrl } from '@/utils/image'
import api from '@/utils/axios'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const auth = useAuthStore()
const usersStore = useUsersStore()

/* dropdown */
const open = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const toggleDropdown = () => {
  open.value = !open.value
}

const handleClickOutside = (e: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    open.value = false
  }
}

/* dialogs */
const showChangePassword = ref(false)
const showResetAll = ref(false)

/* form state */
const changePasswordForm = ref({
  current: '',
  new: ''
})

const resetAllForm = ref({
  new: ''
})

/* actions */
const openChangePassword = () => {
  changePasswordForm.value = { current: '', new: '' }
  showChangePassword.value = true
  open.value = false
}

const openResetAll = () => {
  resetAllForm.value = { new: '' }
  showResetAll.value = true
  open.value = false
}

const submitChangePassword = async () => {
  if (!changePasswordForm.value.current || !changePasswordForm.value.new) {
    push.error('Please fill in all fields')
    return
  }

  const success = await auth.changePassword(
    changePasswordForm.value.current,
    changePasswordForm.value.new
  )

  if (success) {
    changePasswordForm.value = { current: '', new: '' }
    showChangePassword.value = false
  }
}

const submitResetAll = async () => {
  if (!resetAllForm.value.new) return

  await usersStore.resetAllPasswords(resetAllForm.value.new)
  resetAllForm.value = { new: '' }
  showResetAll.value = false
}

/* theme */
const theme = ref<'light' | 'dark'>('light')

const applyTheme = (value: 'light' | 'dark') => {
  theme.value = value
  document.documentElement.setAttribute('data-theme', value)
  localStorage.setItem('theme', value)
}

const toggleTheme = () => {
  applyTheme(theme.value === 'dark' ? 'light' : 'dark')
  open.value = false
}

const themeLabel = computed(() =>
  theme.value === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode'
)

const downloadBackup = async () => {
  open.value = false

  try {
    const res = await api.get('/backup', {
      responseType: 'blob',
    })

    let filename = 'backup.zip'

    const disposition = res.headers['content-disposition']

    if (disposition) {
      const match = disposition.match(/filename="?([^"]+)"?/)
      if (match) filename = match[1]
    }

    const blob = new Blob([res.data])

    const url = window.URL.createObjectURL(blob)

    const a = document.createElement('a')
    a.href = url
    a.download = filename

    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)

    window.URL.revokeObjectURL(url)

    push.success('Backup downloaded successfully!')
  } catch (err) {
    console.error(err)
    push.error('Failed to download backup')
  }
}

const handleLogout = () => {
  open.value = false
  auth.logout()
}

/* helpers */
const getInitials = (name: string): string => {
  return name
    .split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getRoleBadgeClass = (role: string): string => {
  switch (role) {
    case 'admin':
      return 'bg-info-bg text-info-text'
    case 'manager':
      return 'bg-warning-bg text-warning-text'
    default:
      return 'bg-surface-alt text-secondary'
  }
}

/* lifecycle */
onMounted(() => {
  const saved = (localStorage.getItem('theme') as 'light' | 'dark') || 'light'
  applyTheme(saved)
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
