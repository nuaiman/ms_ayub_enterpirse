<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-customer', customer)">

    <!-- CUSTOMER NAME -->
    <div class="col-span-3">
      <div class="flex items-center gap-3">
        <div class="shrink-0">
          <div class="w-9 h-9 rounded-full bg-surface-alt border border-default flex items-center justify-center">
            <span class="text-xs font-semibold text-muted">{{ getInitials(customer.name) }}</span>
          </div>
        </div>
        <div class="min-w-0">
          <div class="font-medium text-primary truncate text-sm">{{ customer.name }}</div>
          <div v-if="customer.email" class="text-xs text-muted truncate flex items-center gap-1 mt-0.5">
            <svg class="w-3 h-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            {{ customer.email }}
          </div>
        </div>
      </div>
    </div>

    <!-- PHONE -->
    <div class="col-span-2">
      <div class="text-sm text-secondary">{{ customer.phone }}</div>
    </div>

    <!-- COMPANY -->
    <div class="col-span-2">
      <div class="text-sm text-secondary truncate">{{ customer.company || '—' }}</div>
    </div>

    <!-- CREATED -->
    <div class="col-span-3">
      <div class="text-xs text-secondary leading-tight">{{ formatDate(customer.created_at) }}</div>
      <div class="text-[10px] text-muted leading-tight">{{ formatTime(customer.created_at) }}</div>
    </div>

    <!-- ACTIONS -->
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
            Edit Details
          </button>

          <div class="border-t border-divider my-0.5"></div>

          <button @click="openDeleteConfirm"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete Customer
          </button>

          <div class="border-t border-divider mt-0.5 pt-0.5">
            <div class="px-3 py-1">
              <p class="text-[10px] text-muted">ID: {{ customer.id }}</p>
            </div>
          </div>
        </div>
      </Transition>

      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- EDIT DIALOG -->
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
            <h2 class="text-lg font-semibold text-primary">Edit Customer</h2>
            <p class="text-xs text-muted">{{ customer.name }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Full Name <span class="text-warning-text">*</span></label>
            <input v-model="editForm.name" type="text" placeholder="John Doe"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Phone <span class="text-warning-text">*</span></label>
            <input v-model="editForm.phone" type="tel" placeholder="+1 234 567 890"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Email</label>
            <input v-model="editForm.email" type="email" placeholder="john@example.com"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Company</label>
            <input v-model="editForm.company" type="text" placeholder="Company name"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Address</label>
            <input v-model="editForm.address" type="text" placeholder="Enter address"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
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
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Notes</label>
            <textarea v-model="editForm.notes" rows="3" placeholder="Add notes about this customer"
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
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

    <!-- DELETE CONFIRM DIALOG -->
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
            <h2 class="text-lg font-semibold text-primary">Delete Customer</h2>
            <p class="text-xs text-muted">{{ customer.name }}</p>
          </div>
        </div>

        <p class="text-sm text-secondary">
          Are you sure you want to delete <strong>{{ customer.name }}</strong>? This action cannot be undone.
        </p>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showDeleteConfirm = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="deleteCustomer"
            class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg">Delete</button>
        </div>
      </div>
    </BaseDialog>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Customer } from '@/types/customer'
import type { IDType } from '@/types/auth'
import { useCustomersStore } from '@/stores/customers'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const props = defineProps<{ customer: Customer }>()
const emit = defineEmits<{
  'customer-updated': []
  'view-customer': [customer: Customer]
}>()

const customersStore = useCustomersStore()

const open = ref(false)
const showEditDialog = ref(false)
const showDeleteConfirm = ref(false)

const editForm = ref({
  name: props.customer.name || '',
  phone: props.customer.phone || '',
  email: props.customer.email || '',
  company: props.customer.company || '',
  address: props.customer.address || '',
  id_type: (props.customer.id_type as IDType | '') || '',
  id_number: props.customer.id_number || '',
  notes: props.customer.notes || '',
})

const openEditDialog = () => {
  editForm.value = {
    name: props.customer.name || '',
    phone: props.customer.phone || '',
    email: props.customer.email || '',
    company: props.customer.company || '',
    address: props.customer.address || '',
    id_type: (props.customer.id_type as IDType | '') || '',
    id_number: props.customer.id_number || '',
    notes: props.customer.notes || '',
  }
  showEditDialog.value = true
  open.value = false
}

const openDeleteConfirm = () => {
  showDeleteConfirm.value = true
  open.value = false
}

const submitEdit = async () => {
  if (!editForm.value.name) { push.error('Name is required'); return }
  if (!editForm.value.phone) { push.error('Phone is required'); return }

  try {
    const success = await customersStore.updateCustomer(props.customer.id, {
      name: editForm.value.name,
      phone: editForm.value.phone,
      email: editForm.value.email || null,
      company: editForm.value.company || null,
      address: editForm.value.address || null,
      id_type: (editForm.value.id_type as IDType) || null,
      id_number: editForm.value.id_number || null,
      notes: editForm.value.notes || null,
    })

    if (success) {
      showEditDialog.value = false
      emit('customer-updated')
    }
  } catch {
    push.error('Failed to update customer')
  }
}

const deleteCustomer = async () => {
  try {
    const success = await customersStore.deleteCustomer(props.customer.id)
    if (success) {
      showDeleteConfirm.value = false
    }
  } catch {
    push.error('Failed to delete customer')
  }
}

const getInitials = (name: string): string => name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
const formatTime = (d: string): string => new Date(d).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
</script>
