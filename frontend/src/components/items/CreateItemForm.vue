<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add New Item</h2>
      <p class="text-sm text-secondary">Register a new item in inventory</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Customer Selection (Admin/Manager only) -->
        <div v-if="canManageCustomers" class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Customer <span class="text-xs font-normal text-muted normal-case">(Optional)</span>
          </h3>

          <!-- Selected Customer -->
          <div v-if="selectedCustomer"
            class="flex items-start gap-4 p-4 rounded-xl border border-success-border bg-success-bg">
            <div class="shrink-0 w-14 h-14 rounded-full bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-lg font-semibold text-muted">#{{ selectedCustomer.id }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-primary">{{ selectedCustomer.name }}</p>
              <div class="flex items-center gap-3 mt-1 flex-wrap">
                <span class="text-xs text-secondary">{{ selectedCustomer.phone }}</span>
                <span v-if="selectedCustomer.email" class="text-xs text-muted">{{ selectedCustomer.email }}</span>
                <span v-if="selectedCustomer.company" class="text-xs text-muted">· {{ selectedCustomer.company }}</span>
              </div>
            </div>
            <button @click="clearSelectedCustomer" type="button"
              class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Search (when no customer selected) -->
          <div v-else class="relative">
            <div class="relative">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none" fill="none"
                stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input ref="customerSearchRef" v-model="customerSearch" type="text"
                placeholder="Search customers by name, phone, email, company..." @focus="showCustomerDropdown = true"
                @keydown.escape="showCustomerDropdown = false"
                class="input w-full pl-9 pr-8 py-2.5 rounded-lg text-sm placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button v-if="customerSearch" @click="customerSearch = ''" type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Dropdown -->
            <div v-if="showCustomerDropdown"
              class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
              <div class="max-h-56 overflow-y-auto">
                <div v-if="filteredCustomers.length === 0" class="px-4 py-6 text-center">
                  <div class="w-10 h-10 mx-auto rounded-full bg-surface-alt flex items-center justify-center mb-2">
                    <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  </div>
                  <p class="text-sm text-muted">No customers match "{{ customerSearch }}"</p>
                </div>
                <button v-for="c in filteredCustomers" :key="c.id" @click="selectCustomer(c)" type="button"
                  class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                  <div class="shrink-0">
                    <div class="w-9 h-9 rounded-full bg-surface-alt border border-default flex items-center justify-center">
                      <span class="text-xs font-semibold text-muted">#{{ c.id }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-primary truncate">{{ c.name }}</p>
                    <div class="flex items-center gap-2 mt-0.5">
                      <span class="text-xs text-muted">{{ c.phone }}</span>
                      <span v-if="c.email" class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded">{{ c.email }}</span>
                      <span v-if="c.company" class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded">{{ c.company }}</span>
                    </div>
                  </div>
                  <svg class="w-4 h-4 text-muted shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </div>

              <!-- Footer with count -->
              <div class="px-4 py-2 border-t border-divider bg-surface-alt">
                <p class="text-[11px] text-muted">{{ filteredCustomers.length }} customer{{ filteredCustomers.length !== 1 ? 's' : '' }} found</p>
              </div>
            </div>

            <!-- Backdrop -->
            <div v-if="showCustomerDropdown" class="fixed inset-0 z-40" @click="showCustomerDropdown = false"></div>
          </div>
        </div>

        <!-- Contract Details (Admin/Manager only) -->
        <div v-if="canManageCustomers" class="pb-6 border-b border-divider">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-semibold text-muted uppercase tracking-wider">
              Contract Details <span class="text-xs font-normal normal-case">(Optional)</span>
            </h3>
            <button type="button" @click="showContract = !showContract"
              class="text-sm text-accent hover:text-accent/80 transition-colors">
              {{ showContract ? 'Hide' : 'Add Contract' }}
            </button>
          </div>

          <div v-if="showContract" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Duration Type</label>
              <select v-model="duration_type"
                class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option value="">Select Type</option>
                <option value="day">Day</option>
                <option value="week">Week</option>
                <option value="month">Month</option>
                <option value="year">Year</option>
              </select>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Duration</label>
              <input v-model.number="duration" type="number" min="1" placeholder="Duration"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Price</label>
              <input v-model.number="price" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Security Deposit</label>
              <input v-model.number="security_deposit" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Estimated Value</label>
              <input v-model.number="estimated_value" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
        </div>

        <!-- Basic Info -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Item Details</h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5 md:col-span-2">
              <label class="text-sm font-medium text-primary">Name <span class="text-warning-text">*</span></label>
              <input v-model="name" type="text" placeholder="Item name" required
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Quantity</label>
              <div class="flex items-center gap-2">
                <button type="button" @click="quantity > 1 ? quantity-- : null"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                  </svg>
                </button>
                <input v-model.number="quantity" type="number" min="1" placeholder="1"
                  class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button type="button" @click="quantity++"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                  </svg>
                </button>
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Unit</label>
              <select v-model="quantity_unit"
                class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option value="pcs">pcs</option>
                <option value="g">g</option>
                <option value="kg">kg</option>
                <option value="ton">ton</option>
                <option value="ml">ml</option>
                <option value="liter">liter</option>
              </select>
            </div>

            <div class="space-y-1.5">
  <label class="text-sm font-medium text-primary">Weight per Item (kg)</label>
  <input v-model.number="weight" type="number" step="0.001" min="0" placeholder="0.000"
    class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
  <p v-if="weight && quantity > 1" class="text-xs text-muted mt-1">
    Total weight: <span class="font-medium text-primary">{{ (weight * quantity).toFixed(3) }} kg</span>
  </p>
</div>
          </div>
        </div>

        <!-- Dimensions -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Dimensions <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="grid grid-cols-3 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Length (in)</label>
              <input v-model.number="length" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Width (in)</label>
              <input v-model.number="width" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Height (in)</label>
              <input v-model.number="height" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
        </div>

        <!-- Image Upload -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Image <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="flex items-start gap-4">
            <div class="shrink-0">
              <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default">
                <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                <button type="button" @click="removeImage"
                  class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90">×</button>
              </div>
              <div v-else
                class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt">
                <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
            </div>
            <div class="flex-1 space-y-3">
              <label class="text-sm font-medium text-primary block">Upload Image</label>
              <div class="flex gap-2">
                <label
                  class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2"
                  :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                  </svg>
                  {{ uploading ? 'Uploading...' : 'Choose File' }}
                  <input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" />
                </label>
              </div>
              <div v-if="uploading" class="w-full bg-surface-alt rounded-full h-1.5 overflow-hidden">
                <div class="bg-accent h-full rounded-full transition-all duration-300"
                  :style="{ width: uploadProgress + '%' }"></div>
              </div>
              <p v-if="uploadError" class="text-xs text-warning-text">{{ uploadError }}</p>
            </div>
          </div>
        </div>

        <!-- Notes -->
        <div>
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Notes <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <textarea v-model="notes" rows="3" placeholder="Add notes..."
            class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit" :disabled="submitting || !name"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Creating...
            </span>
            <span v-else>Create Item</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { QuantityUnit, DurationType, CreateItemPayload } from '@/types/item'
import type { Customer } from '@/types/customer'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import { useAuthStore } from '@/stores/auth'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{ 'item-created': [] }>()

const itemsStore = useItemsStore()
const customersStore = useCustomersStore()
const authStore = useAuthStore()

// Check if user can manage customers (admin or manager)
const canManageCustomers = computed(() => {
  const role = authStore.user?.role
  return role === 'admin' || role === 'manager'
})

const submitting = ref(false)

// Customer fields
const customerSearchRef = ref<HTMLInputElement | null>(null)
const customerSearch = ref('')
const showCustomerDropdown = ref(false)
const selectedCustomer = ref<Customer | null>(null)

// Contract fields
const showContract = ref(false)
const duration_type = ref<DurationType | null>(null)
const duration = ref<number | null>(null)
const price = ref<number | null>(null)
const security_deposit = ref<number | null>(null)
const estimated_value = ref<number | null>(null)

// Item fields
const name = ref('')
const quantity = ref(1)
const quantity_unit = ref<QuantityUnit>('pcs')
const weight = ref<number | null>(null)
const length = ref<number | null>(null)
const width = ref<number | null>(null)
const height = ref<number | null>(null)
const notes = ref('')

// Image fields
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)

const filteredCustomers = computed(() => {
  if (!customerSearch.value) return customersStore.customers
  const q = customerSearch.value.toLowerCase()
  return customersStore.customers.filter(c =>
    c.name.toLowerCase().includes(q) ||
    c.phone.includes(q) ||
    (c.email && c.email.toLowerCase().includes(q)) ||
    (c.company && c.company.toLowerCase().includes(q))
  )
})

const selectCustomer = (c: Customer) => {
  selectedCustomer.value = c
  customerSearch.value = ''
  showCustomerDropdown.value = false
}

const clearSelectedCustomer = () => {
  selectedCustomer.value = null
  customerSearch.value = ''
}

onMounted(async () => {
  // Only fetch customers if user can manage them
  if (canManageCustomers.value && customersStore.customers.length === 0) {
    await customersStore.fetchCustomers()
  }
})

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { uploadError.value = 'Please select an image file'; return }
  if (file.size > 5 * 1024 * 1024) { uploadError.value = 'Image size should be less than 5MB'; return }
  uploadError.value = null
  imageFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => { imagePreview.value = e.target?.result as string }
  reader.readAsDataURL(file)
}

const removeImage = () => {
  imageFile.value = null
  imagePreview.value = null
  uploadError.value = null
  uploadProgress.value = 0
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
  if (fileInput) fileInput.value = ''
}

const resetForm = () => {
  name.value = ''
  quantity.value = 1
  quantity_unit.value = 'pcs'
  weight.value = null
  length.value = null
  width.value = null
  height.value = null
  notes.value = ''
  showContract.value = false
  duration_type.value = null
  duration.value = null
  price.value = null
  security_deposit.value = null
  estimated_value.value = null
  clearSelectedCustomer()
  removeImage()
}

const submit = async () => {
  if (!name.value) return
  submitting.value = true

  try {
    const payload: CreateItemPayload = {
      name: name.value,
      quantity_unit: quantity_unit.value,
      quantity: quantity.value,
      weight: weight.value || null,
      length: length.value || null,
      width: width.value || null,
      height: height.value || null,
      notes: notes.value || null,
    }

    // Only include customer_id and contract fields if user has permission
    if (canManageCustomers.value) {
      payload.customer_id = selectedCustomer.value?.id || null

      // Add contract fields if contract section is shown and has values
      if (showContract.value && duration_type.value) {
        payload.duration_type = duration_type.value
        payload.duration = duration.value || null
        payload.price = price.value || null
        payload.security_deposit = security_deposit.value || null
        payload.estimated_value = estimated_value.value || null
      }
    }

    const newItem = await itemsStore.createItem(payload)

    if (!newItem) { submitting.value = false; return }

    if (imageFile.value) {
      uploading.value = true
      uploadProgress.value = 0
      const imageUrl = await uploadImage('items', newItem.id, imageFile.value, (progress) => {
        uploadProgress.value = progress
      })
      uploading.value = false
      if (!imageUrl) push.warning('Item created but image upload failed')
    }

    push.success('Item created successfully!')
    resetForm()
    emit('item-created')
  } catch (error) {
    console.error('Error creating item:', error)
    push.error('Failed to create item')
  } finally {
    submitting.value = false
  }
}
</script>
