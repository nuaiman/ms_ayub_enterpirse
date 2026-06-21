<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-item', item)">

    <div class="col-span-3">
      <div class="flex items-center gap-3">
        <div class="shrink-0">
          <div v-if="item.image_url" class="w-9 h-9 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(item.image_url) || undefined" :alt="item.name" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-9 h-9 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
            <svg class="w-4 h-4 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
        </div>
        <div class="min-w-0">
          <div class="font-medium text-primary truncate">{{ item.name }}</div>
          <div class="text-xs text-muted">{{ item.quantity }} {{ item.quantity_unit }}</div>
        </div>
      </div>
    </div>

    <div class="col-span-2">
      <span class="text-xs text-secondary">{{ item.weight ? item.weight + ' kg' : '—' }}</span>
    </div>

    <div class="col-span-2">
      <div v-if="item.contract_id" @click.stop>
        <button @click="showContractDialog = true"
          class="inline-flex items-center gap-1.5 px-2 py-1 rounded-lg hover:bg-surface-alt transition-colors text-left">
          <span class="text-sm font-medium text-primary">#{{ item.contract_id }}</span>
          <span v-if="getContract(item.contract_id)" class="text-[10px] px-1.5 py-0.5 rounded-full border"
            :class="getContract(item.contract_id)?.status === 'active' ? 'status-success' : getContract(item.contract_id)?.status === 'completed' ? 'status-info' : 'status-warning'">
            {{ getContract(item.contract_id)?.status }}
          </span>
        </button>
      </div>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <div class="col-span-2">
      <span v-if="item.length || item.width || item.height" class="text-xs text-muted">
        {{ [item.length, item.width, item.height].filter(Boolean).join('×') }}
      </span>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <div class="col-span-1">
      <span class="text-xs text-muted">{{ formatDate(item.created_at) }}</span>
    </div>

    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <button @click="open = !open"
        class="w-7 h-7 flex items-center justify-center border border-default rounded-md hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100"
        :class="{ 'opacity-100': open }">
        <svg class="w-3.5 h-3.5 text-secondary" fill="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" />
          <circle cx="12" cy="12" r="1.5" />
          <circle cx="12" cy="19" r="1.5" />
        </svg>
      </button>

      <Transition enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100 translate-y-0"
        leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="open"
          class="absolute right-0 top-9 w-44 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider">
            <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
          </div>
          <button @click="openEditDialog"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit Item
          </button>
          <div class="border-t border-divider my-0.5"></div>
          <button @click="openDeleteConfirm"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete Item
          </button>
        </div>
      </Transition>
      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- Contract Quick View Dialog -->
    <BaseDialog v-model="showContractDialog" size="md" @click.stop>
      <div v-if="item.contract_id && getContract(item.contract_id)" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-lg font-bold text-primary">Contract #{{ item.contract_id }}</h2>
            <p class="text-sm text-muted capitalize">{{ getContract(item.contract_id)?.status }}</p>
          </div>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border"
            :class="getContract(item.contract_id)?.status === 'active' ? 'status-success' : getContract(item.contract_id)?.status === 'completed' ? 'status-info' : 'status-warning'">
            {{ getContract(item.contract_id)?.status }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Duration</p>
            <p class="text-xs text-primary">{{ getContract(item.contract_id)?.duration || '—' }} {{
              getContract(item.contract_id)?.duration_type }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Price</p>
            <p class="text-xs text-primary">{{ getContract(item.contract_id)?.price ? '$' +
              getContract(item.contract_id)?.price : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Sec. Deposit</p>
            <p class="text-xs text-primary">{{ getContract(item.contract_id)?.security_deposit ? '$' +
              getContract(item.contract_id)?.security_deposit : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Est. Value</p>
            <p class="text-xs text-primary">{{ getContract(item.contract_id)?.estimated_value ? '$' +
              getContract(item.contract_id)?.estimated_value : '—' }}</p>
          </div>
        </div>

        <!-- Customer from contract -->
        <div v-if="getCustomerFromContract(item.contract_id)"
          class="p-3 rounded-lg border border-default bg-surface-alt">
          <p class="text-[10px] text-muted uppercase mb-2">Customer</p>
          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 rounded-full bg-surface border border-default flex items-center justify-center shrink-0">
              <span class="text-xs font-semibold text-muted">#{{ getContract(item.contract_id)?.customer_id }}</span>
            </div>
            <div>
              <p class="text-sm font-medium text-primary">{{ getCustomerFromContract(item.contract_id)?.name }}</p>
              <p class="text-xs text-muted">{{ getCustomerFromContract(item.contract_id)?.phone }}</p>
            </div>
          </div>
        </div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showContractDialog = false"
            class="button px-4 py-2 text-xs rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>

    <!-- Edit Dialog -->
    <BaseDialog v-model="showEditDialog" size="lg" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-success-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Edit Item</h2>
            <p class="text-xs text-muted">{{ item.name }}</p>
          </div>
        </div>

        <!-- Contract Assignment (Admin/Manager only) -->
        <div v-if="canAssignContract" class="pb-4 mb-4 border-b border-divider">
          <label class="text-sm font-medium text-primary block mb-2">Assigned Contract</label>

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

          <div v-else class="relative">
            <div class="relative">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none" fill="none"
                stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input v-model="contractSearch" type="text" placeholder="Search contracts to assign..."
                @focus="showContractDropdown = true" @keydown.escape="showContractDropdown = false"
                class="input w-full pl-9 pr-8 py-2.5 rounded-lg text-sm placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button v-if="contractSearch" @click="contractSearch = ''" type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div v-if="showContractDropdown"
              class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
              <div class="max-h-48 overflow-y-auto">
                <div v-if="filteredContracts.length === 0" class="px-4 py-6 text-center">
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
                </button>
              </div>
            </div>
            <div v-if="showContractDropdown" class="fixed inset-0 z-40" @click="showContractDropdown = false"></div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Name <span class="text-warning-text">*</span></label>
            <input v-model="editForm.name" type="text"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Quantity</label>
            <input v-model.number="editForm.quantity" type="number" min="1"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Unit</label>
            <select v-model="editForm.quantity_unit"
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
            <input v-model.number="editForm.weight" type="number" step="0.001" min="0" placeholder="0.000"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 pt-2">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Length</label>
            <input v-model.number="editForm.length" type="number" step="0.01" min="0" placeholder="0.00"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Width</label>
            <input v-model.number="editForm.width" type="number" step="0.01" min="0" placeholder="0.00"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Height</label>
            <input v-model.number="editForm.height" type="number" step="0.01" min="0" placeholder="0.00"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
        </div>

        <div class="space-y-1.5 pt-2">
          <label class="text-sm font-medium text-primary">Notes</label>
          <textarea v-model="editForm.notes" rows="2" placeholder="Add notes..."
            class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t border-divider mt-2">
          <button type="button" @click="showEditDialog = false"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Cancel</button>
          <button type="button" @click="handleUpdate" :disabled="saving || !editForm.name"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="saving" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Saving...
            </span>
            <span v-else>Save Changes</span>
          </button>
        </div>
      </div>
    </BaseDialog>

    <!-- Delete Dialog -->
    <BaseDialog v-model="showDeleteDialog" size="sm" @click.stop>
      <div class="space-y-5">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-warning-text/10 flex items-center justify-center shrink-0">
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold text-primary">Delete Item</h3>
            <p class="text-xs text-muted">This action cannot be undone.</p>
          </div>
        </div>
        <p class="text-sm text-secondary">
          Are you sure you want to delete <span class="font-semibold text-primary">"{{ item.name }}"</span>?
        </p>
        <div class="flex justify-end gap-3 pt-2">
          <button type="button" @click="showDeleteDialog = false"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Cancel</button>
          <button type="button" @click="handleDelete" :disabled="deleting"
            class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2">
            <span v-if="deleting">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Deleting...
            </span>
            <span v-else>Delete</span>
          </button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import type { Item, QuantityUnit } from '@/types/item'
import type { Contract } from '@/types/contract'
import { useItemsStore } from '@/stores/items'
import { useContractsStore } from '@/stores/contracts'
import { useCustomersStore } from '@/stores/customers'
import { useAuthStore } from '@/stores/auth'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'
import { getImageUrl } from '@/utils/image'

const props = defineProps<{
  item: Item
}>()

const emit = defineEmits<{
  'view-item': [item: Item]
  'item-updated': []
  'item-deleted': []
}>()

const itemsStore = useItemsStore()
const contractsStore = useContractsStore()
const customersStore = useCustomersStore()
const auth = useAuthStore()

const canAssignContract = computed(() => auth.user?.role === 'admin' || auth.user?.role === 'manager')

const open = ref(false)
const showContractDialog = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const saving = ref(false)
const deleting = ref(false)

const contractSearch = ref('')
const showContractDropdown = ref(false)
const selectedContract = ref<Contract | null>(null)

const editForm = reactive({
  name: '',
  quantity: 1,
  quantity_unit: 'pcs' as QuantityUnit,
  weight: null as number | null,
  length: null as number | null,
  width: null as number | null,
  height: null as number | null,
  notes: '' as string | null
})

const filteredContracts = computed(() => {
  if (!contractSearch.value) return contractsStore.contracts
  const q = contractSearch.value.toLowerCase()
  return contractsStore.contracts.filter(c =>
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

const formatDate = (dateStr: string) => {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const getContract = (id: number) => {
  return contractsStore.contracts.find(c => c.id === id)
}

const getCustomerFromContract = (contractId: number) => {
  const contract = getContract(contractId)
  if (!contract) return null
  return customersStore.customers.find(c => c.id === contract.customer_id)
}

const openEditDialog = () => {
  open.value = false

  // Set current form values
  editForm.name = props.item.name
  editForm.quantity = props.item.quantity
  editForm.quantity_unit = props.item.quantity_unit
  editForm.weight = props.item.weight ?? null
  editForm.length = props.item.length ?? null
  editForm.width = props.item.width ?? null
  editForm.height = props.item.height ?? null
  editForm.notes = props.item.notes ?? ''

  // Load existing contract if any
  if (props.item.contract_id) {
    const existingContract = getContract(props.item.contract_id)
    selectedContract.value = existingContract || null
  } else {
    selectedContract.value = null
  }

  contractSearch.value = ''
  showEditDialog.value = true
}

const handleUpdate = async () => {
  if (!editForm.name) return
  saving.value = true

  try {
    await itemsStore.updateItem(props.item.id, {
      name: editForm.name,
      quantity: editForm.quantity,
      quantity_unit: editForm.quantity_unit,
      weight: editForm.weight,
      length: editForm.length,
      width: editForm.width,
      height: editForm.height,
      notes: editForm.notes || null,
      contract_id: selectedContract.value?.id ?? null // Included the selected contract
    })

    push.success('Item updated successfully')
    showEditDialog.value = false
    emit('item-updated')
  } catch (error) {
    console.error('Error updating item:', error)
    push.error('Failed to update item')
  } finally {
    saving.value = false
  }
}

const openDeleteConfirm = () => {
  open.value = false
  showDeleteDialog.value = true
}

const handleDelete = async () => {
  deleting.value = true

  try {
    await itemsStore.deleteItem(props.item.id)
    push.success('Item deleted successfully')
    showDeleteDialog.value = false
    emit('item-deleted')
  } catch (error) {
    console.error('Error deleting item:', error)
    push.error('Failed to delete item')
  } finally {
    deleting.value = false
  }
}
</script>
