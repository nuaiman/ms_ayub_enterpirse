<template>
  <div class="max-w-2xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Create New User</h2>
      <p class="text-sm text-secondary">Add a new team member to the system</p>
    </div>

    <!-- Form Card -->
    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Personal Information Section -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Personal Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Name -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Full Name <span class="text-warning-text">*</span>
              </label>
              <input v-model="name" type="text" placeholder="John Doe" required class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Username -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Username <span class="text-warning-text">*</span>
              </label>
              <input v-model="username" type="text" placeholder="johndoe" required class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Email -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Email</label>
              <input v-model="email" type="email" placeholder="john@example.com" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Phone -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Phone</label>
              <input v-model="phone" type="tel" placeholder="+1 234 567 890" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>
          </div>

          <!-- Address -->
          <div class="mt-4 space-y-1.5">
            <label class="text-sm font-medium text-primary">Address</label>
            <textarea v-model="address" rows="2" placeholder="Enter full address" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                     focus:outline-none focus:ring-2 focus:ring-accent
                     focus:border-transparent transition-all duration-200 resize-none"></textarea>
          </div>
        </div>

        <!-- Account Settings Section -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Account Settings
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Password -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Password <span class="text-warning-text">*</span>
              </label>
              <div class="relative">
                <input v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="••••••••" required
                  class="input w-full px-3 py-2 pr-10 rounded-lg placeholder:text-muted
                         focus:outline-none focus:ring-2 focus:ring-accent
                         focus:border-transparent transition-all duration-200" />
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

            <!-- Role -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Role <span class="text-warning-text">*</span>
              </label>
              <select v-model="role" required class="input w-full px-3 py-2 rounded-lg cursor-pointer
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent
                       transition-all duration-200">
                <option value="staff" class="bg-surface text-primary">Staff</option>
                <option value="accounts" class="bg-surface text-primary">Accounts</option>
                <option v-if="canCreateManager" value="manager" class="bg-surface text-primary">Manager</option>
              </select>
            </div>
          </div>
        </div>

        <!-- ID Information Section -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Identification <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">ID Type</label>
              <select v-model="id_type" class="input w-full px-3 py-2 rounded-lg cursor-pointer
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent
                       transition-all duration-200">
                <option value="" class="bg-surface text-primary">Select ID Type</option>
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
              <input v-model="id_number" type="text" placeholder="Enter ID number" :disabled="!id_type" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200
                       disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
          </div>
        </div>

        <!-- Profile Image Section -->
        <div>
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Profile Image <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>

          <div class="space-y-4">
            <div class="flex items-start gap-4">
              <!-- Preview -->
              <div class="shrink-0">
                <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default">
                  <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                  <button type="button" @click="removeImage" class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full
                           flex items-center justify-center text-xs hover:opacity-90 transition-opacity">
                    ×
                  </button>
                </div>
                <div v-else class="w-24 h-24 rounded-lg border-2 border-dashed border-default
                         flex items-center justify-center bg-surface-alt">
                  <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                      d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>

              <!-- Upload Controls -->
              <div class="flex-1 space-y-3">
                <div>
                  <label class="text-sm font-medium text-primary block mb-1.5">Upload Image</label>
                  <div class="flex gap-2">
                    <label class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer
                             hover-surface transition-all duration-200 inline-flex items-center gap-2"
                      :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                      </svg>
                      {{ uploading ? 'Uploading...' : 'Choose File' }}
                      <input type="file" accept="image/*" class="hidden" @change="handleFileSelect"
                        :disabled="uploading" />
                    </label>
                  </div>
                </div>

                <!-- Upload Progress -->
                <div v-if="uploading" class="w-full bg-surface-alt rounded-full h-1.5 overflow-hidden">
                  <div class="bg-accent h-full rounded-full transition-all duration-300"
                    :style="{ width: uploadProgress + '%' }"></div>
                </div>

                <p v-if="uploadError" class="text-xs text-warning-text">{{ uploadError }}</p>
                <p v-else-if="uploadSuccess" class="text-xs text-success-text">{{ uploadSuccess }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">
            Clear Form
          </button>
          <button type="submit" :disabled="submitting" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg
                   shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 transform
                   disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Creating...
            </span>
            <span v-else>Create User</span>
          </button>
        </div>

      </form>
    </div>
  </div>
</template>


<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { IDType } from '@/types/auth'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{
  'user-created': []
}>()

const auth = useAuthStore()
const usersStore = useUsersStore()

const showPassword = ref(false)
const submitting = ref(false)

// Form fields
const name = ref('')
const username = ref('')
const password = ref('')
const email = ref('')
const phone = ref('')
const address = ref('')
const id_type = ref<IDType | ''>('')
const id_number = ref('')

// Image upload state
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)
const uploadSuccess = ref<string | null>(null)

const canCreateManager = computed(() => auth.user?.role === 'admin')
const role = ref<'staff' | 'accounts' | 'manager'>('staff')

watch(id_type, (newVal) => {
  if (!newVal) id_number.value = ''
})

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select an image file'
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'Image size should be less than 5MB'
    return
  }

  uploadError.value = null
  uploadSuccess.value = null
  imageFile.value = file

  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

const removeImage = () => {
  imageFile.value = null
  imagePreview.value = null
  uploadError.value = null
  uploadSuccess.value = null
  uploadProgress.value = 0

  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
  if (fileInput) fileInput.value = ''
}

const resetForm = () => {
  name.value = ''
  username.value = ''
  password.value = ''
  email.value = ''
  phone.value = ''
  address.value = ''
  id_type.value = ''
  id_number.value = ''
  removeImage()
  role.value = 'staff'
  showPassword.value = false
}

const submit = async () => {
  if (role.value === 'manager' && auth.user?.role !== 'admin') return

  submitting.value = true

  try {
    const newUser = await usersStore.createUser({
      name: name.value,
      username: username.value,
      password: password.value,
      role: role.value,
      email: email.value || null,
      phone: phone.value || null,
      address: address.value || null,
      id_type: id_type.value || null,
      id_number: id_number.value || null,
      image_url: null,
    })

    if (!newUser) {
      submitting.value = false
      return
    }

    if (imageFile.value) {
      uploading.value = true
      uploadProgress.value = 0

      const imageUrl = await uploadImage('users', newUser.id, imageFile.value, (progress) => {
        uploadProgress.value = progress
      })

      uploading.value = false

      if (!imageUrl) {
        push.warning('User created but image upload failed')
      }
    }

    push.success('User created successfully!')
    resetForm()
    emit('user-created') // 👈 tells parent to close dialog
  } catch (error) {
    console.error('Error creating user:', error)
    push.error('Failed to create user')
  } finally {
    submitting.value = false
  }
}
</script>
