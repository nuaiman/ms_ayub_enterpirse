<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add New Item</h2>
      <p class="text-sm text-secondary">Register a new item in inventory</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Contract Selection - TOP (only for admin/manager) -->
        <div v-if="canAssignContract" class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Contract <span class="text-xs font-normal text-muted normal-case">(Optional)</span>
          </h3>

          <!-- Selected Contract -->
          <div v-if="selectedContract"
            class="flex items-center gap-3 p-3 rounded-lg border border-success-border bg-success-bg">
            <div
              class="w-10 h-10 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
              <span class="text-sm font-semibold text-muted">#{{ selectedContract.id }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-primary">Contract #{{ selectedContract.id }}</p>
              <p class="text-xs text-secondary">Customer #{{ selectedContract.customer_id }} · <span
                  class="capitalize">{{ selectedContract.status }}</span></p>
            </div>
            <button @click="clearSelectedContract" type="button"
              class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Search (when no contract selected) -->
          <div v-else class="relative">
            <div class="relative">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none" fill="none"
                stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input v-model="contractSearch" type="text" placeholder="Search contracts by ID, customer, or status..."
                @focus="showContractDropdown = true" @keydown.escape="showContractDropdown = false"
                class="input w-full pl-9 pr-8 py-2.5 rounded-lg text-sm placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button v-if="contractSearch" @click="contractSearch = ''" type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Dropdown -->
            <div v-if="showContractDropdown"
              class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
              <div class="px-3 py-2 border-b border-divider">
                <button @click="openQuickAddContract" type="button"
                  class="w-full flex items-center gap-2 px-2 py-2 rounded-lg text-sm font-medium text-accent hover:bg-accent/5 transition-colors">
                  <div class="w-8 h-8 rounded-full bg-accent/10 flex items-center justify-center shrink-0">
                    <svg class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                    </svg>
                  </div>
                  <span>Create New Contract</span>
                </button>
              </div>
              <div class="max-h-48 overflow-y-auto">
                <div v-if="filteredContracts.length === 0" class="px-4 py-6 text-center">
                  <div class="w-10 h-10 mx-auto rounded-full bg-surface-alt flex items-center justify-center mb-2">
                    <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <p class="text-sm text-muted">No contracts match "{{ contractSearch }}"</p>
                </div>
                <button v-for="c in filteredContracts" :key="c.id" @click="selectContract(c)" type="button"
                  class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                  <div
                    class="w-9 h-9 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
                    <span class="text-xs font-semibold text-muted">#{{ c.id }}</span>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-primary">Contract #{{ c.id }}</p>
                    <div class="flex items-center gap-2 mt-0.5">
                      <span class="text-xs text-muted">Customer #{{ c.customer_id }}</span>
                      <span class="text-[10px] px-1.5 py-0.5 rounded-full border"
                        :class="c.status === 'active' ? 'status-success' : c.status === 'completed' ? 'status-info' : 'status-warning'">{{
                        c.status }}</span>
                    </div>
                  </div>
                  <svg class="w-4 h-4 text-muted shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </div>
              <div class="px-4 py-2 border-t border-divider bg-surface-alt">
                <p class="text-[11px] text-muted">{{ filteredContracts.length }} contract{{ filteredContracts.length !==
                  1 ? 's' : '' }} found</p>
              </div>
            </div>
            <div v-if="showContractDropdown" class="fixed inset-0 z-40" @click="showContractDropdown = false"></div>
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
                <option value="pcs" class="bg-surface text-primary">pcs</option>
                <option value="g" class="bg-surface text-primary">g</option>
                <option value="kg" class="bg-surface text-primary">kg</option>
                <option value="ton" class="bg-surface text-primary">ton</option>
                <option value="ml" class="bg-surface text-primary">ml</option>
                <option value="liter" class="bg-surface text-primary">liter</option>
              </select>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Weight (kg)</label>
              <input v-model.number="weight" type="number" step="0.001" min="0" placeholder="0.000"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
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
              <label class="text-sm font-medium text-primary">Length</label>
              <input v-model.number="length" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Width</label>
              <input v-model.number="width" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Height</label>
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

    <BaseDialog v-if="canAssignContract" v-model="showQuickAddContract" size="lg">
      <CreateContractForm @contract-created="handleContractAdded" />
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { QuantityUnit } from '@/types/item'
import type { Contract } from '@/types/contract'
import { useItemsStore } from '@/stores/items'
import { useContractsStore } from '@/stores/contracts'
import { useAuthStore } from '@/stores/auth'
import { uploadImage } from '@/utils/image'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import CreateContractForm from '@/components/contracts/CreateContractForm.vue'
import { push } from 'notivue'

const emit = defineEmits<{ 'item-created': [] }>()

const itemsStore = useItemsStore()
const contractsStore = useContractsStore()
const auth = useAuthStore()

const canAssignContract = computed(() => auth.user?.role === 'admin' || auth.user?.role === 'manager')

const submitting = ref(false)
const name = ref('')
const quantity = ref(1)
const quantity_unit = ref<QuantityUnit>('pcs')
const weight = ref<number | null>(null)
const length = ref<number | null>(null)
const width = ref<number | null>(null)
const height = ref<number | null>(null)
const notes = ref('')

const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)

const contractSearch = ref('')
const showContractDropdown = ref(false)
const selectedContract = ref<Contract | null>(null)
const showQuickAddContract = ref(false)
const contracts = ref<Contract[]>([])

onMounted(async () => {
  if (canAssignContract.value) {
    await contractsStore.fetchContracts()
    contracts.value = contractsStore.contracts
  }
})

const filteredContracts = computed(() => {
  if (!contractSearch.value) return contracts.value
  const q = contractSearch.value.toLowerCase()
  return contracts.value.filter(c =>
    String(c.id).includes(q) ||
    String(c.customer_id).includes(q) ||
    c.status.toLowerCase().includes(q)
  )
})

const selectContract = (c: Contract) => {
  selectedContract.value = c
  contractSearch.value = ''
  showContractDropdown.value = false
}

const clearSelectedContract = () => {
  selectedContract.value = null
  contractSearch.value = ''
}

const openQuickAddContract = () => {
  showContractDropdown.value = false
  showQuickAddContract.value = true
}

const handleContractAdded = async () => {
  showQuickAddContract.value = false
  await contractsStore.fetchContracts()
  contracts.value = contractsStore.contracts
  if (contracts.value.length > 0) {
    const last = contracts.value[0]
    if (last) {
      selectedContract.value = last
      push.success('Contract added and selected')
    }
  }
}

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
  clearSelectedContract()
  removeImage()
}

const submit = async () => {
  if (!name.value) return
  submitting.value = true

  try {
    const newItem = await itemsStore.createItem({
      name: name.value,
      quantity_unit: quantity_unit.value,
      quantity: quantity.value,
      weight: weight.value,
      contract_id: selectedContract.value?.id ?? null,
      length: length.value,
      width: width.value,
      height: height.value,
      notes: notes.value || null,
    })

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
