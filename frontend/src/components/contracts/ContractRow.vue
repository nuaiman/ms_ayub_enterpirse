<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-contract', contract)">

    <div class="col-span-2">
      <span class="text-sm font-medium text-primary">#{{ contract.id }}</span>
    </div>

    <div class="col-span-2" @click.stop>
      <button @click="showCustomerDialog = true"
        class="inline-flex items-center gap-2.5 px-2 py-1.5 rounded-lg hover:bg-surface-alt transition-colors text-left w-full">
        <div
          class="w-8 h-8 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
          <span class="text-xs font-semibold text-muted">#{{ contract.customer_id }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-xs font-medium text-primary truncate">{{ customerName }}</div>
          <div class="text-[10px] text-muted truncate">{{ customer?.phone || '—' }}</div>
          <div v-if="customer?.company" class="text-[10px] text-muted truncate">{{ customer.company }}</div>
        </div>
      </button>
    </div>

    <div class="col-span-2">
      <div class="text-sm text-secondary">
        {{ contract.duration || '—' }} {{ contract.duration_type }}{{ contract.duration && contract.duration > 1 ? 's' :
          '' }}
      </div>
    </div>

    <div class="col-span-2">
      <div class="text-sm text-secondary">
        <span v-if="contract.price">${{ contract.price }}</span>
        <span v-else>—</span>
      </div>
    </div>

    <div class="col-span-2">
      <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
        :class="getStatusClass(contract.status)">
        {{ contract.status }}
      </span>
    </div>

    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <button @click="open = !open" class="w-7 h-7 flex items-center justify-center border border-default rounded-md
               hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100"
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
          class="absolute right-0 top-9 w-52 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider">
            <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
          </div>

          <button @click="openEditDialog"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit Contract
          </button>

          <button v-if="contract.status === 'active'" @click="completeContract"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-success-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            Mark Completed
          </button>

          <button v-if="contract.status === 'active'" @click="cancelContract"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Cancel Contract
          </button>

          <div class="border-t border-divider my-0.5"></div>

          <button @click="openDeleteConfirm"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete Contract
          </button>
        </div>
      </Transition>

      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- Customer Quick View Dialog -->
    <BaseDialog v-model="showCustomerDialog" size="md" @click.stop>
      <div v-if="customer" class="space-y-5">
        <div class="flex items-center gap-4">
          <div class="shrink-0">
            <div class="w-14 h-14 rounded-full bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-sm font-semibold text-muted">#{{ customer.id }}</span>
            </div>
          </div>
          <div>
            <h2 class="text-lg font-bold text-primary">{{ customer.name }}</h2>
            <p class="text-sm text-muted">{{ customer.phone }}</p>
            <p v-if="customer.company" class="text-xs text-secondary">{{ customer.company }}</p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Email</p>
            <p class="text-xs text-primary">{{ customer.email || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Address</p>
            <p class="text-xs text-primary">{{ customer.address || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">ID Type</p>
            <p class="text-xs text-primary capitalize">{{ customer.id_type || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">ID Number</p>
            <p class="text-xs text-primary">{{ customer.id_number || '—' }}</p>
          </div>
          <div class="space-y-1 col-span-2">
            <p class="text-[10px] text-muted uppercase">Notes</p>
            <p class="text-xs text-primary whitespace-pre-wrap">{{ customer.notes || '—' }}</p>
          </div>
        </div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showCustomerDialog = false"
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
            <h2 class="text-lg font-semibold text-primary">Edit Contract #{{ contract.id }}</h2>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Customer <span class="text-warning-text">*</span></label>
            <select v-model.number="editForm.customer_id" required
              class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
              <option v-for="c in customers" :key="c.id" :value="c.id" class="bg-surface text-primary">
                {{ c.name }} — {{ c.phone }}
              </option>
            </select>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Duration Type <span
                class="text-warning-text">*</span></label>
            <select v-model="editForm.duration_type" required
              class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
              <option value="day" class="bg-surface text-primary">Day</option>
              <option value="week" class="bg-surface text-primary">Week</option>
              <option value="month" class="bg-surface text-primary">Month</option>
              <option value="year" class="bg-surface text-primary">Year</option>
            </select>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Duration</label>
            <input v-model.number="editForm.duration" type="number" min="1"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Price</label>
            <input v-model.number="editForm.price" type="number" step="0.01" min="0"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Security Deposit</label>
            <input v-model.number="editForm.security_deposit" type="number" step="0.01" min="0"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Estimated Value</label>
            <input v-model.number="editForm.estimated_value" type="number" step="0.01" min="0"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Status</label>
            <select v-model="editForm.status"
              class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
              <option value="active" class="bg-surface text-primary">Active</option>
              <option value="completed" class="bg-surface text-primary">Completed</option>
              <option value="cancelled" class="bg-surface text-primary">Cancelled</option>
            </select>
          </div>
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Notes</label>
            <textarea v-model="editForm.notes" rows="3"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
          </div>
        </div>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showEditDialog = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="submitEdit"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Save Changes</button>
        </div>
      </div>
    </BaseDialog>

    <!-- Delete Confirm Dialog -->
    <BaseDialog v-model="showDeleteConfirm" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-warning-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Delete Contract</h2>
            <p class="text-xs text-muted">Contract #{{ contract.id }}</p>
          </div>
        </div>
        <p class="text-sm text-secondary">
          Are you sure you want to delete contract <strong>#{{ contract.id }}</strong>? This action cannot be undone.
        </p>
        <div class="flex gap-2 justify-end pt-2">
          <button @click="showDeleteConfirm = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="deleteContract"
            class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg">Delete</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Contract, DurationType, ContractStatus } from '@/types/contract'
import type { Customer } from '@/types/customer'
import { useContractsStore } from '@/stores/contracts'
import { useCustomersStore } from '@/stores/customers'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const props = defineProps<{ contract: Contract }>()
const emit = defineEmits<{
  'contract-updated': []
  'view-contract': [contract: Contract]
}>()

const contractsStore = useContractsStore()
const customersStore = useCustomersStore()

const open = ref(false)
const showEditDialog = ref(false)
const showDeleteConfirm = ref(false)
const showCustomerDialog = ref(false)

const customers = ref(customersStore.customers)

const customer = computed<Customer | null>(() => {
  return customers.value.find(c => c.id === props.contract.customer_id) || null
})

const customerName = computed(() => customer.value?.name || `Customer #${props.contract.customer_id}`)

onMounted(async () => {
  if (customers.value.length === 0) {
    await customersStore.fetchCustomers()
    customers.value = customersStore.customers
  }
})

const editForm = ref({
  customer_id: props.contract.customer_id,
  duration_type: props.contract.duration_type as DurationType,
  duration: props.contract.duration ?? null,
  price: props.contract.price ?? null,
  security_deposit: props.contract.security_deposit ?? null,
  estimated_value: props.contract.estimated_value ?? null,
  status: props.contract.status as ContractStatus,
  notes: props.contract.notes || '',
})

const openEditDialog = () => {
  editForm.value = {
    customer_id: props.contract.customer_id,
    duration_type: props.contract.duration_type as DurationType,
    duration: props.contract.duration ?? null,
    price: props.contract.price ?? null,
    security_deposit: props.contract.security_deposit ?? null,
    estimated_value: props.contract.estimated_value ?? null,
    status: props.contract.status as ContractStatus,
    notes: props.contract.notes || '',
  }
  showEditDialog.value = true
  open.value = false
}

const openDeleteConfirm = () => {
  showDeleteConfirm.value = true
  open.value = false
}

const submitEdit = async () => {
  try {
    const success = await contractsStore.updateContract(props.contract.id, {
      customer_id: editForm.value.customer_id,
      duration_type: editForm.value.duration_type,
      duration: editForm.value.duration,
      price: editForm.value.price,
      security_deposit: editForm.value.security_deposit,
      estimated_value: editForm.value.estimated_value,
      status: editForm.value.status,
      notes: editForm.value.notes || null,
    })
    if (success) {
      showEditDialog.value = false
      emit('contract-updated')
    }
  } catch {
    push.error('Failed to update contract')
  }
}

const completeContract = async () => {
  await contractsStore.changeContractStatus(props.contract.id, 'completed')
  open.value = false
}

const cancelContract = async () => {
  await contractsStore.changeContractStatus(props.contract.id, 'cancelled')
  open.value = false
}

const deleteContract = async () => {
  try {
    const success = await contractsStore.deleteContract(props.contract.id)
    if (success) {
      showDeleteConfirm.value = false
    }
  } catch {
    push.error('Failed to delete contract')
  }
}

const getStatusClass = (status: string): string => {
  switch (status) {
    case 'active': return 'status-success'
    case 'completed': return 'status-info'
    case 'cancelled': return 'status-warning'
    default: return ''
  }
}
</script>
