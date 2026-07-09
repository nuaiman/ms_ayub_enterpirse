<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">

    <!-- HEADER -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Users</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ usersStore.users.length }}
          total</span>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search users..."
            class="input w-full sm:w-64 pl-9 pr-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent placeholder:text-muted" />
          <button v-if="searchQuery" @click="searchQuery = ''" type="button"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="flex items-center gap-2">
          <div class="relative">
            <button @click="toggleSortDropdown"
              class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors"
              title="Sort">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
              </svg>
            </button>
            <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95"
              enter-to-class="opacity-100 scale-100" leave-active-class="transition ease-in duration-150"
              leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
              <div v-if="showSortMenu"
                class="absolute right-0 top-10 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
                <div class="px-3 py-2 border-b border-divider">
                  <p class="text-xs font-semibold text-muted uppercase">Sort By</p>
                </div>
                <button v-for="option in sortOptions" :key="option.key" @click="applySort(option.key)"
                  class="w-full flex items-center justify-between px-3 py-2 text-sm text-secondary hover:bg-surface-alt transition-colors">
                  {{ option.label }}
                  <svg v-if="currentSort === option.key" class="w-4 h-4 text-accent" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
              </div>
            </Transition>
          </div>

          <button @click="copyToClipboard"
            class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors relative"
            title="Copy to clipboard">
            <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>

          <button v-if="canCreate" @click="openCreateDialog"
            class="h-9 px-4 flex items-center gap-2 bg-accent text-accent-foreground rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
            </svg>
            <span class="hidden sm:inline">Create</span>
          </button>
        </div>
      </div>
    </div>

    <div v-if="showSortMenu" class="fixed inset-0 z-40" @click="showSortMenu = false"></div>

    <!-- TABLE -->
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-3xl">
        <div
          class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3 flex items-center gap-1 cursor-pointer hover:text-primary transition-colors"
            @click="toggleSort('name')">Name
            <svg v-if="sortField === 'name'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2">Role</div>
          <div class="col-span-1">Salary</div>
          <div class="col-span-2">Status</div>
          <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
            @click="toggleSort('created_at')">Created
            <svg v-if="sortField === 'created_at'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2 text-right pr-0.5">Actions</div>
        </div>

        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-center space-y-4">
            <svg class="animate-spin w-10 h-10 text-accent mx-auto" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            <p class="text-sm text-muted">Loading users...</p>
          </div>
        </div>

        <div v-else-if="filteredUsers.length === 0" class="flex items-center justify-center py-12">
          <div class="text-center space-y-3">
            <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
              <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-primary">No users found</p>
            <p class="text-xs text-muted mt-1">{{ searchQuery ? 'adjust your search' : 'create new user' }}</p>
          </div>
        </div>

        <div v-else>
          <UserRow v-for="user in filteredUsers" :key="user.id" :user="user" @user-updated="handleUserUpdated"
            @view-user="openDetailDialog" />
        </div>
      </div>
    </div>

    <div v-if="!loading && filteredUsers.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredUsers.length }} of {{ usersStore.users.length }} users</p>
      <button @click="refreshUsers"
        class="button px-3 py-1.5 text-xs rounded-lg hover-surface transition-colors inline-flex items-center gap-1.5"
        :disabled="refreshing">
        <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor"
          viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </button>
    </div>

    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateUserForm @user-created="handleUserCreated" />
    </BaseDialog>

    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedUser" class="space-y-5">
        <div class="flex items-center gap-4">
          <div class="shrink-0">
            <div v-if="selectedUser.image_url" class="w-16 h-16 rounded-full overflow-hidden border border-default">
              <img :src="getImageUrl(selectedUser.image_url) || undefined" :alt="selectedUser.name"
                class="w-full h-full object-cover" />
            </div>
            <div v-else
              class="w-16 h-16 rounded-full bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-xl font-semibold text-muted">{{ getInitials(selectedUser.name) }}</span>
            </div>
          </div>
          <div>
            <h2 class="text-xl font-bold text-primary">{{ selectedUser.name }}</h2>
            <p class="text-sm text-muted">@{{ selectedUser.username }}</p>
            <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize mt-1"
              :class="getRoleBadgeClass(selectedUser.role)">
              <span class="w-1 h-1 rounded-full mr-1" :class="getRoleDotClass(selectedUser.role)"></span>
              {{ selectedUser.role }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-4 pb-4 border-b border-divider">
          <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border"
            :class="selectedUser.is_active ? 'status-success' : 'status-warning'">
            <span class="w-1.5 h-1.5 rounded-full"
              :class="selectedUser.is_active ? 'bg-success-text' : 'bg-warning-text'"></span>
            {{ selectedUser.is_active ? 'Active' : 'Inactive' }}
          </span>
          <span class="text-xs text-muted">ID: {{ selectedUser.id }}</span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Email</p>
            <p class="text-sm text-primary">{{ selectedUser.email || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Phone</p>
            <p class="text-sm text-primary">{{ selectedUser.phone || '—' }}</p>
          </div>
          <div class="space-y-1 md:col-span-2">
            <p class="text-xs text-muted uppercase tracking-wider">Address</p>
            <p class="text-sm text-primary">{{ selectedUser.address || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Monthly Salary</p>
            <p class="text-sm text-primary">{{ selectedUser.monthly_salary ? '৳' +
              selectedUser.monthly_salary.toFixed(2) : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">ID Type</p>
            <p class="text-sm text-primary capitalize">{{ selectedUser.id_type || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">ID Number</p>
            <p class="text-sm text-primary">{{ selectedUser.id_number || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Created</p>
            <p class="text-sm text-primary">{{ formatDate(selectedUser.created_at) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Updated</p>
            <p class="text-sm text-primary">{{ formatDate(selectedUser.updated_at) }}</p>
          </div>
        </div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showDetailDialog = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/types/auth'
import UserRow from './UserRow.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import CreateUserForm from '@/components/users/CreateUserForm.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const usersStore = useUsersStore()
const auth = useAuthStore()

const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const sortField = ref<'name' | 'created_at'>('created_at')
const sortDirection = ref<'asc' | 'desc'>('desc')
const showSortMenu = ref(false)
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedUser = ref<User | null>(null)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'name_asc', label: 'Name A-Z' },
  { key: 'name_desc', label: 'Name Z-A' },
]

const canCreate = computed(() => auth.user?.role === 'admin' || auth.user?.role === 'manager')

onMounted(async () => { await fetchUsers() })

const fetchUsers = async () => { loading.value = true; try { await usersStore.fetchUsers() } finally { loading.value = false } }
const refreshUsers = async () => { refreshing.value = true; try { await usersStore.fetchUsers() } finally { refreshing.value = false } }
const handleUserUpdated = async () => { await usersStore.fetchUsers() }
const handleUserCreated = async () => { showCreateDialog.value = false; await usersStore.fetchUsers() }
const openCreateDialog = () => { showCreateDialog.value = true }

const openDetailDialog = (user: User) => {
  selectedUser.value = user
  showDetailDialog.value = true
}

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }

const applySort = (key: string) => {
  currentSort.value = key; showSortMenu.value = false
  switch (key) {
    case 'newest': sortField.value = 'created_at'; sortDirection.value = 'desc'; break
    case 'oldest': sortField.value = 'created_at'; sortDirection.value = 'asc'; break
    case 'name_asc': sortField.value = 'name'; sortDirection.value = 'asc'; break
    case 'name_desc': sortField.value = 'name'; sortDirection.value = 'desc'; break
  }
}

const copyToClipboard = async () => {
  const users = filteredUsers.value
  if (users.length === 0) return
  const headers = ['ID', 'Name', 'Username', 'Role', 'Status', 'Email', 'Phone', 'Address', 'Salary', 'ID Type', 'ID Number', 'Created']
  const rows = users.map(u => [
    u.id, u.name, u.username, u.role, u.is_active ? 'Active' : 'Inactive',
    u.email || '', u.phone || '', u.address || '', u.monthly_salary || '',
    u.id_type || '', u.id_number || '', formatDate(u.created_at),
  ])
  const tsv = [headers, ...rows].map(row => row.map(cell => String(cell).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')
  try {
    await navigator.clipboard.writeText(tsv)
    copied.value = true
    push.success('Copied to clipboard!')
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    const textarea = document.createElement('textarea')
    textarea.value = tsv
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    push.success('Copied to clipboard!')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

const filteredUsers = computed(() => {
  const users = [...usersStore.users]
  let result = users
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = users.filter(u =>
      u.name.toLowerCase().includes(query) ||
      u.username.toLowerCase().includes(query) ||
      (u.email && u.email.toLowerCase().includes(query)) ||
      u.role.toLowerCase().includes(query) ||
      u.phone?.toLowerCase().includes(query) ||
      u.address?.toLowerCase().includes(query) ||
      u.id_type?.toLowerCase().includes(query) ||
      u.id_number?.toLowerCase().includes(query)
    )
  }
  result.sort((a, b) => {
    let comparison = 0
    if (sortField.value === 'name') comparison = a.name.localeCompare(b.name)
    else if (sortField.value === 'created_at') comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    return sortDirection.value === 'desc' ? -comparison : comparison
  })
  return result
})

const toggleSort = (field: 'name' | 'created_at') => {
  if (sortField.value === field) sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  else { sortField.value = field; sortDirection.value = 'asc' }
}

const getInitials = (name: string): string => name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })

const getRoleBadgeClass = (role: string): string => {
  switch (role) { case 'admin': return 'bg-info-bg text-info-text'; case 'manager': return 'bg-warning-bg text-warning-text'; default: return 'bg-surface-alt text-secondary' }
}
const getRoleDotClass = (role: string): string => {
  switch (role) { case 'admin': return 'bg-info-text'; case 'manager': return 'bg-warning-text'; default: return 'bg-muted' }
}
</script>
