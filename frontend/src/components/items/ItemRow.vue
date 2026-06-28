<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="!isStaff && emit('view-item', item)">

    <!-- Name -->
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

    <!-- Customer -->
    <div class="col-span-2">
      <span v-if="item.customer_id" class="text-sm text-secondary truncate block">
        {{ getCustomerName(item.customer_id) }}
      </span>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Contract - Staff only see "—" -->
    <div class="col-span-2">
      <div v-if="!isStaff && hasContract(item)" @click.stop>
        <button @click="showContractDialog = true"
          class="inline-flex items-center gap-1.5 px-2 py-1 rounded-lg hover:bg-surface-alt transition-colors text-left">
          <span class="text-xs font-medium text-primary">Contract</span>
          <span class="text-[10px] px-1.5 py-0.5 rounded-full border" :class="getStatusClass(item.status || 'active')">
            {{ item.status || 'active' }}
          </span>
        </button>
      </div>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Dimensions -->
    <div class="col-span-2">
      <span v-if="item.length || item.width || item.height" class="text-xs text-muted">
        {{ [item.length, item.width, item.height].filter(Boolean).join('×') }}
      </span>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Created -->
    <div class="col-span-1">
      <span class="text-xs text-muted">{{ formatDate(item.created_at) }}</span>
    </div>

    <!-- Actions - Staff see nothing -->
    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <!-- Only show actions dropdown for non-staff users -->
      <template v-if="!isStaff">
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
            class="absolute right-0 top-9 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
            <div class="px-3 py-1.5 border-b border-divider">
              <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
            </div>

            <!-- Edit -->
            <button @click="openEditDialog"
              class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
              <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Edit Item
            </button>

            <!-- Change Status (only for items with contract) -->
            <div v-if="hasContract(item)">
              <div class="px-3 py-1">
                <p class="text-[10px] text-muted uppercase">Change Status</p>
              </div>
              <button
                v-for="status in availableStatuses"
                :key="status"
                @click="handleStatusChange(status)"
                class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-surface-alt transition-colors"
                :class="item.status === status ? 'text-accent font-semibold' : 'text-secondary'"
              >
                <span class="w-2 h-2 rounded-full" :class="getStatusDotClass(status)"></span>
                {{ status.charAt(0).toUpperCase() + status.slice(1) }}
                <span v-if="item.status === status" class="ml-auto text-accent">✓</span>
              </button>
            </div>

            <div class="border-t border-divider my-0.5"></div>

            <!-- Delete -->
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
      </template>

      <!-- For staff users, show nothing -->
    </div>

    <!-- Contract Quick View Dialog (only for non-staff) -->
    <BaseDialog v-if="!isStaff" v-model="showContractDialog" size="md" @click.stop>
      <div v-if="hasContract(item)" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-lg font-bold text-primary">Contract Details</h2>
            <p class="text-sm text-muted">{{ item.name }}</p>
          </div>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border"
            :class="getStatusClass(item.status || 'active')">
            {{ item.status || 'active' }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Duration</p>
            <p class="text-xs text-primary">{{ item.duration || '—' }} {{ item.duration_type }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Price</p>
            <p class="text-xs text-primary">{{ item.price ? '৳' + item.price : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Sec. Deposit</p>
            <p class="text-xs text-primary">{{ item.security_deposit ? '৳' + item.security_deposit : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Est. Value</p>
            <p class="text-xs text-primary">{{ item.estimated_value ? '৳' + item.estimated_value : '—' }}</p>
          </div>
        </div>

        <!-- Customer from item -->
        <div v-if="item.customer_id" class="p-3 rounded-lg border border-default bg-surface-alt">
          <p class="text-[10px] text-muted uppercase mb-2">Customer</p>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-surface border border-default flex items-center justify-center shrink-0">
              <span class="text-xs font-semibold text-muted">#{{ item.customer_id }}</span>
            </div>
            <div>
              <p class="text-sm font-medium text-primary">{{ getCustomerName(item.customer_id) }}</p>
              <p class="text-xs text-muted">{{ getCustomerPhone(item.customer_id) }}</p>
            </div>
          </div>
        </div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showContractDialog = false" class="button px-4 py-2 text-xs rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>

    <!-- Edit Dialog (only for non-staff) -->
    <BaseDialog v-if="!isStaff" v-model="showEditDialog" size="lg" @click.stop>
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

        <!-- Customer Assignment -->
        <div class="pb-4 mb-4 border-b border-divider">
          <label class="text-sm font-medium text-primary block mb-2">Customer</label>
          <div class="relative">
            <div v-if="editCustomer" class="flex items-center gap-3 p-3 rounded-lg border border-success-border bg-success-bg">
              <div class="w-10 h-10 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
                <span class="text-sm font-semibold text-muted">#{{ editCustomer.id }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-primary">{{ editCustomer.name }}</p>
                <p class="text-xs text-secondary">{{ editCustomer.phone }}</p>
              </div>
              <button @click="clearEditCustomer" type="button"
                class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div v-else>
              <input v-model="editCustomerSearch" type="text" placeholder="Search customers..."
                @focus="showEditCustomerDropdown = true"
                class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <div v-if="showEditCustomerDropdown"
                class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface max-h-48 overflow-y-auto">
                <button v-for="c in editFilteredCustomers" :key="c.id" @click="selectEditCustomer(c)" type="button"
                  class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-primary">{{ c.name }}</p>
                    <p class="text-xs text-muted">{{ c.phone }}</p>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Contract Details (Optional) - MATCHES CREATE FORM -->
        <div class="pb-4 mb-4 border-b border-divider">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-semibold text-muted uppercase tracking-wider">
              Contract Details <span class="text-xs font-normal normal-case">(Optional)</span>
            </h3>
            <button type="button" @click="showEditContract = !showEditContract"
              class="text-sm text-accent hover:text-accent/80 transition-colors">
              {{ showEditContract ? 'Hide' : 'Add Contract' }}
            </button>
          </div>

          <div v-if="showEditContract" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Duration Type</label>
              <select v-model="editContract.duration_type"
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
              <input v-model.number="editContract.duration" type="number" min="1" placeholder="Duration"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Price</label>
              <input v-model.number="editContract.price" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Security Deposit</label>
              <input v-model.number="editContract.security_deposit" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Estimated Value</label>
              <input v-model.number="editContract.estimated_value" type="number" step="0.01" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
        </div>

        <!-- Form Fields -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Name <span class="text-warning-text">*</span></label>
            <input v-model="editForm.name" type="text"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Quantity</label>
            <div class="flex items-center gap-2">
              <button type="button" @click="editForm.quantity > 1 ? editForm.quantity-- : null"
                class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                </svg>
              </button>
              <input v-model.number="editForm.quantity" type="number" min="1" placeholder="1"
                class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button type="button" @click="editForm.quantity++"
                class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                </svg>
              </button>
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Unit</label>
            <select v-model="editForm.quantity_unit"
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
            <label class="text-sm font-medium text-primary">Weight (kg)</label>
            <input v-model.number="editForm.weight" type="number" step="0.001" min="0" placeholder="0.000"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
        </div>

        <!-- Dimensions -->
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
          <textarea v-model="editForm.notes" rows="3" placeholder="Add notes..."
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

    <!-- Delete Dialog (only for non-staff) -->
    <BaseDialog v-if="!isStaff" v-model="showDeleteDialog" size="sm" @click.stop>
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
import type { Item, QuantityUnit, DurationType, ItemStatus } from '@/types/item'
import { hasContract } from '@/types/item'
import type { Customer } from '@/types/customer'
import { useItemsStore } from '@/stores/items'
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
const customersStore = useCustomersStore()
const authStore = useAuthStore()

// Check if user is staff (only view permission)
const isStaff = computed(() => authStore.user?.role === 'staff')

const open = ref(false)
const showContractDialog = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const saving = ref(false)
const deleting = ref(false)

// Edit contract visibility - default to false, user must click "Add Contract"
const showEditContract = ref(false)

// Available statuses for contract items
const availableStatuses: ItemStatus[] = ['active', 'completed', 'cancelled']

// Edit form state
const editForm = reactive({
  name: '',
  quantity: 1,
  quantity_unit: 'pcs' as QuantityUnit,
  weight: null as number | null,
  length: null as number | null,
  width: null as number | null,
  height: null as number | null,
  notes: null as string | null,
})

// Edit contract state
const editContract = reactive({
  duration_type: null as DurationType | null,
  duration: null as number | null,
  price: null as number | null,
  security_deposit: null as number | null,
  estimated_value: null as number | null,
  status: null as ItemStatus | null,
})

// Edit customer state
const editCustomerSearch = ref('')
const showEditCustomerDropdown = ref(false)
const editCustomer = ref<Customer | null>(null)

const editFilteredCustomers = computed(() => {
  if (!editCustomerSearch.value) return customersStore.customers
  const q = editCustomerSearch.value.toLowerCase()
  return customersStore.customers.filter(c =>
    c.name.toLowerCase().includes(q) ||
    c.phone.includes(q)
  )
})

const getCustomerName = (id: number): string => {
  const customer = customersStore.customers.find(c => c.id === id)
  return customer?.name || `#${id}`
}

const getCustomerPhone = (id: number): string => {
  const customer = customersStore.customers.find(c => c.id === id)
  return customer?.phone || ''
}

const getStatusClass = (status: string) => {
  const classes = {
    active: 'border-green-400 bg-green-50 text-green-700',
    completed: 'border-blue-400 bg-blue-50 text-blue-700',
    cancelled: 'border-red-400 bg-red-50 text-red-700'
  }
  return classes[status as keyof typeof classes] || classes.active
}

const getStatusDotClass = (status: string) => {
  const classes = {
    active: 'bg-green-500',
    completed: 'bg-blue-500',
    cancelled: 'bg-red-500'
  }
  return classes[status as keyof typeof classes] || 'bg-gray-500'
}

const selectEditCustomer = (c: Customer) => {
  editCustomer.value = c
  editCustomerSearch.value = ''
  showEditCustomerDropdown.value = false
}

const clearEditCustomer = () => {
  editCustomer.value = null
  editCustomerSearch.value = ''
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
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
  editForm.notes = props.item.notes ?? null

  // Load customer
  if (props.item.customer_id) {
    const customer = customersStore.customers.find(c => c.id === props.item.customer_id)
    editCustomer.value = customer || null
  } else {
    editCustomer.value = null
  }

  // Load contract fields if exists
  if (hasContract(props.item)) {
    editContract.duration_type = props.item.duration_type || null
    editContract.duration = props.item.duration || null
    editContract.price = props.item.price || null
    editContract.security_deposit = props.item.security_deposit || null
    editContract.estimated_value = props.item.estimated_value || null
    editContract.status = props.item.status || null
    showEditContract.value = true // Show contract section if item has contract
  } else {
    // Reset contract fields if no contract
    editContract.duration_type = null
    editContract.duration = null
    editContract.price = null
    editContract.security_deposit = null
    editContract.estimated_value = null
    editContract.status = null
    showEditContract.value = false // Hide contract section if no contract
  }

  showEditDialog.value = true
}

const handleStatusChange = async (status: string) => {
  if (status === props.item.status) return // Don't change if same status

  try {
    await itemsStore.updateItemStatus(props.item.id, status)
    push.success(`Contract status updated to ${status}`)
    emit('item-updated')
    open.value = false // Close the actions dropdown
  } catch (error) {
    console.error('Error updating status:', error)
    push.error('Failed to update contract status')
  }
}

const handleUpdate = async () => {
  if (!editForm.name) return
  saving.value = true

  try {
    // Build base payload with all item fields
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {
      name: editForm.name,
      quantity: editForm.quantity,
      quantity_unit: editForm.quantity_unit,
    }

    // Only include optional fields if they have values
    if (editForm.weight !== null && editForm.weight !== undefined) {
      payload.weight = editForm.weight
    }
    if (editForm.length !== null && editForm.length !== undefined) {
      payload.length = editForm.length
    }
    if (editForm.width !== null && editForm.width !== undefined) {
      payload.width = editForm.width
    }
    if (editForm.height !== null && editForm.height !== undefined) {
      payload.height = editForm.height
    }
    if (editForm.notes) {
      payload.notes = editForm.notes
    }
    if (editCustomer.value?.id) {
      payload.customer_id = editCustomer.value.id
    } else {
      payload.customer_id = null
    }

    // Handle contract fields
    if (showEditContract.value) {
      // User wants to add/edit contract
      if (editContract.duration_type) {
        // Only include contract fields if duration_type is selected
        payload.duration_type = editContract.duration_type

        if (editContract.duration !== null && editContract.duration !== undefined) {
          payload.duration = editContract.duration
        }
        if (editContract.price !== null && editContract.price !== undefined) {
          payload.price = editContract.price
        }
        if (editContract.security_deposit !== null && editContract.security_deposit !== undefined) {
          payload.security_deposit = editContract.security_deposit
        }
        if (editContract.estimated_value !== null && editContract.estimated_value !== undefined) {
          payload.estimated_value = editContract.estimated_value
        }
        // If item already had a contract, keep the status or update it
        if (editContract.status) {
          payload.status = editContract.status
        } else if (hasContract(props.item)) {
          // Keep existing status if not changed
          payload.status = props.item.status
        } else {
          // New contract, default to active
          payload.status = 'active'
        }
      } else {
        // User opened contract section but didn't select duration_type
        // Keep existing contract if any, or don't add new one
        if (hasContract(props.item)) {
          // Keep existing contract fields
          payload.duration_type = props.item.duration_type
          payload.duration = props.item.duration
          payload.price = props.item.price
          payload.security_deposit = props.item.security_deposit
          payload.estimated_value = props.item.estimated_value
          payload.status = props.item.status
        }
        // If no existing contract and no duration_type selected, don't add contract
      }
    } else {
      // User hid the contract section
      if (hasContract(props.item)) {
        // If item had a contract and user hid the section, remove the contract
        payload.duration_type = null
        payload.duration = null
        payload.price = null
        payload.security_deposit = null
        payload.estimated_value = null
        payload.status = null
      }
      // If item didn't have a contract and section is hidden, do nothing
    }

    // Log what we're sending for debugging
    console.log('Sending payload:', payload)

    await itemsStore.updateItem(props.item.id, payload)
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
