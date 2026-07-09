<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-user', user)">

    <!-- USER INFO -->
    <div class="col-span-3">
      <div class="flex items-center gap-3">
        <div class="shrink-0">
          <div v-if="user.image_url" class="w-9 h-9 rounded-full overflow-hidden border border-default">
            <img :src="getImageUrl(user.image_url) || undefined" :alt="user.name" class="w-full h-full object-cover" />
          </div>
          <div v-else
            class="w-9 h-9 rounded-full bg-surface-alt border border-default flex items-center justify-center">
            <span class="text-xs font-semibold text-muted">{{ getInitials(user.name) }}</span>
          </div>
        </div>
        <div class="min-w-0">
          <div class="font-medium text-primary truncate text-sm">{{ user.name }}</div>
          <div class="text-xs text-muted truncate">@{{ user.username }}</div>
        </div>
      </div>
    </div>

    <!-- ROLE -->
    <div class="col-span-2 flex items-center">
      <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize"
        :class="getRoleClass(user.role)">
        <span class="w-1 h-1 rounded-full mr-1" :class="getRoleDotClass(user.role)"></span>
        {{ user.role }}
      </span>
    </div>

    <!-- SALARY -->
    <div class="col-span-1">
      <span class="text-xs text-secondary">{{ user.monthly_salary ? '৳' + user.monthly_salary.toFixed(0) : '—' }}</span>
    </div>

    <!-- STATUS -->
    <div class="col-span-2 flex items-center">
      <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium border"
        :class="user.is_active ? 'status-success' : 'status-warning'">
        <span class="w-1 h-1 rounded-full" :class="user.is_active ? 'bg-success-text' : 'bg-warning-text'"></span>
        {{ user.is_active ? 'Active' : 'Inactive' }}
      </span>
    </div>

    <!-- CREATED -->
    <div class="col-span-2">
      <div class="text-xs text-secondary leading-tight">{{ formatDate(user.created_at) }}</div>
      <div class="text-[10px] text-muted leading-tight">{{ formatTime(user.created_at) }}</div>
    </div>

    <!-- ACTIONS -->
    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <button v-if="!isAdmin" @click="open = !open" class="w-7 h-7 flex items-center justify-center border border-default rounded-md
               hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100"
        :class="{ 'opacity-100': open }">
        <svg class="w-3.5 h-3.5 text-secondary" fill="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" />
          <circle cx="12" cy="12" r="1.5" />
          <circle cx="12" cy="19" r="1.5" />
        </svg>
      </button>

      <Transition v-if="!isAdmin" enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100 translate-y-0"
        leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="open"
          class="absolute right-0 top-9 w-52 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider">
            <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
          </div>

          <button @click="openEditDetails"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit Details
          </button>

          <div class="border-t border-divider my-0.5"></div>

          <button @click="openPasswordDialog"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
            Change Password
          </button>

          <button @click="openRoleDialog"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            Change Role
          </button>

          <button @click="toggleActive"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg v-if="user.is_active" class="w-3.5 h-3.5 shrink-0 text-warning-text" fill="none" stroke="currentColor"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
            <svg v-else class="w-3.5 h-3.5 shrink-0 text-success-text" fill="none" stroke="currentColor"
              viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ user.is_active ? 'Deactivate User' : 'Activate User' }}
          </button>
        </div>
      </Transition>

      <div v-if="open && !isAdmin" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- EDIT DETAILS DIALOG -->
    <BaseDialog v-model="showEditDetails" size="lg" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-success-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Edit User Details</h2>
            <p class="text-xs text-muted">{{ user.name }}</p>
          </div>
        </div>

        <div class="flex items-center gap-4 pb-4 border-b border-divider">
          <div class="shrink-0 relative group/edit-avatar">
            <div v-if="editImagePreview || user.image_url"
              class="w-16 h-16 rounded-full overflow-hidden border border-default">
              <img :src="editImagePreview || getImageUrl(user.image_url) || undefined" :alt="user.name"
                class="w-full h-full object-cover" />
            </div>
            <div v-else
              class="w-16 h-16 rounded-full bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-lg font-semibold text-muted">{{ getInitials(user.name) }}</span>
            </div>
            <button @click="triggerEditImageInput"
              class="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center opacity-0 group-hover/edit-avatar:opacity-100 transition-opacity cursor-pointer">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
            <input ref="editImageInputRef" type="file" accept="image/*" class="hidden"
              @change="handleEditImageSelect" />
          </div>
          <div>
            <p class="text-xs text-muted">Click on avatar to change photo</p>
            <button v-if="user.image_url" @click="removePhoto"
              class="text-xs text-warning-text hover:underline mt-1">Remove photo</button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Full Name</label>
            <input v-model="editForm.name" type="text" placeholder="John Doe"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Username</label>
            <input :value="user.username" type="text" disabled
              class="input w-full px-3 py-2 rounded-lg opacity-50 cursor-not-allowed" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Email</label>
            <input v-model="editForm.email" type="email" placeholder="john@example.com"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Phone</label>
            <input v-model="editForm.phone" type="tel" placeholder="+1 234 567 890"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Address</label>
            <input v-model="editForm.address" type="text" placeholder="Enter address"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Monthly Salary</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
              <input v-model.number="editForm.monthly_salary" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">ID Type</label>
            <select v-model="editForm.id_type"
              class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
              <option value="" class="bg-surface text-primary">None</option>
              <option value="nid" class="bg-surface text-primary">NID</option>
              <option value="passport" class="bg-surface text-primary">Passport</option>
              <option value="driving_license" class="bg-surface text-primary">Driving License</option>
              <option value="birth_certificate" class="bg-surface text-primary">Birth Certificate</option>
              <option value="trade_license" class="bg-surface text-primary">Trade License</option>
              <option value="other" class="bg-surface text-primary">Other</option>
            </select>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">ID Number</label>
            <input v-model="editForm.id_number" type="text" placeholder="Enter ID number" :disabled="!editForm.id_type"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed" />
          </div>
        </div>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showEditDetails = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="submitEditDetails"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Save Changes</button>
        </div>
      </div>
    </BaseDialog>

    <!-- PASSWORD DIALOG -->
    <BaseDialog v-model="showPassword" @click.stop>
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
            <p class="text-xs text-muted">For {{ user.name }}</p>
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-primary">New Password</label>
          <input v-model="password" type="password" placeholder="Enter new password"
            class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent"
            @keyup.enter="updatePassword" />
        </div>
        <div class="flex gap-2 justify-end pt-2">
          <button @click="showPassword = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="updatePassword"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Update Password</button>
        </div>
      </div>
    </BaseDialog>

    <!-- ROLE DIALOG -->
    <BaseDialog v-model="showRole" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-warning-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Change Role</h2>
            <p class="text-xs text-muted">For {{ user.name }}</p>
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="text-sm font-medium text-primary">Role</label>
          <select v-model="role"
            class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
            <option value="manager" class="bg-surface text-primary">Manager</option>
            <option value="accounts" class="bg-surface text-primary">Accounts</option>
            <option value="staff" class="bg-surface text-primary">Staff</option>
          </select>
        </div>
        <div class="flex gap-2 justify-end pt-2">
          <button @click="showRole = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="updateRole"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Update Role</button>
        </div>
      </div>
    </BaseDialog>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { User, IDType } from '@/types/auth'
import { useUsersStore } from '@/stores/users'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const props = defineProps<{ user: User }>()
const emit = defineEmits<{
  'user-updated': []
  'view-user': [user: User]
}>()

const usersStore = useUsersStore()
const isAdmin = computed(() => props.user.role === 'admin')

const open = ref(false)
const showEditDetails = ref(false)
const showPassword = ref(false)
const showRole = ref(false)

const password = ref('')
const role = ref(props.user.role)

const editImageInputRef = ref<HTMLInputElement | null>(null)
const editImageFile = ref<File | null>(null)
const editImagePreview = ref<string | null>(null)

const editForm = ref({
  name: props.user.name || '',
  email: props.user.email || '',
  phone: props.user.phone || '',
  address: props.user.address || '',
  monthly_salary: props.user.monthly_salary || null as number | null,
  id_type: (props.user.id_type as IDType | '') || '',
  id_number: props.user.id_number || '',
})

const openEditDetails = () => {
  editForm.value = {
    name: props.user.name || '',
    email: props.user.email || '',
    phone: props.user.phone || '',
    address: props.user.address || '',
    monthly_salary: props.user.monthly_salary || null,
    id_type: (props.user.id_type as IDType | '') || '',
    id_number: props.user.id_number || '',
  }
  editImageFile.value = null
  editImagePreview.value = null
  showEditDetails.value = true
  open.value = false
}

const openPasswordDialog = () => { showPassword.value = true; open.value = false }
const openRoleDialog = () => { showRole.value = true; open.value = false }

const triggerEditImageInput = () => editImageInputRef.value?.click()

const handleEditImageSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }
  if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }
  editImageFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => { editImagePreview.value = e.target?.result as string }
  reader.readAsDataURL(file)
}

const removePhoto = async () => {
  if (!props.user.image_url) return
  try {
    const success = await deleteImage('users', props.user.id, props.user.image_url)
    if (success) { editImagePreview.value = null; emit('user-updated') }
  } catch { push.error('Failed to remove photo') }
}

const submitEditDetails = async () => {
  if (!editForm.value.name) { push.error('Name is required'); return }
  try {
    if (editImageFile.value) await uploadImage('users', props.user.id, editImageFile.value)
    const success = await usersStore.updateUserProfile(props.user.id, {
      name: editForm.value.name,
      email: editForm.value.email || null,
      phone: editForm.value.phone || null,
      address: editForm.value.address || null,
      monthly_salary: editForm.value.monthly_salary || null,
      id_type: (editForm.value.id_type as IDType) || null,
      id_number: editForm.value.id_number || null,
    })
    if (success) { showEditDetails.value = false; emit('user-updated') }
  } catch { push.error('Failed to update user') }
}

const updatePassword = async () => {
  if (!password.value) return
  await usersStore.changeUserPassword(props.user.id, password.value)
  password.value = ''; showPassword.value = false
}

const updateRole = async () => {
  await usersStore.changeRole(props.user.id, role.value)
  showRole.value = false
}

const toggleActive = async () => {
  await usersStore.toggleActive(props.user.id)
  open.value = false
}

const getInitials = (name: string): string => name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
const formatTime = (d: string): string => new Date(d).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })

const getRoleClass = (role: string): string => {
  switch (role) { case 'admin': return 'bg-info-bg text-info-text'; case 'manager': return 'bg-warning-bg text-warning-text'; default: return 'bg-surface-alt text-secondary' }
}
const getRoleDotClass = (role: string): string => {
  switch (role) { case 'admin': return 'bg-info-text'; case 'manager': return 'bg-warning-text'; default: return 'bg-muted' }
}
</script>
