<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-shipment', shipment)">

    <!-- ID -->
    <div class="col-span-1">
      <span class="text-sm font-medium text-primary">#{{ shipment.id }}</span>
    </div>

    <!-- Type & Status -->
    <div class="col-span-2">
      <div class="flex items-center gap-1.5">
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="shipment.shipment_type === 'pickup' ? 'status-info' : 'status-success'">
          {{ shipment.shipment_type }}
        </span>
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
          :class="getStatusClass(shipment.status)">
          {{ shipment.status }}
        </span>
      </div>
    </div>

    <!-- Item -->
    <div class="col-span-3">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg bg-surface-alt border border-default flex items-center justify-center shrink-0">
          <span class="text-[10px] font-semibold text-muted">#{{ shipment.item_id }}</span>
        </div>
        <div class="min-w-0">
          <div class="text-xs font-medium text-primary truncate">{{ itemName }}</div>
          <div class="text-[10px] text-muted">{{ itemQuantity }} {{ itemUnit }} available</div>
          <div v-if="itemWeight" class="text-[10px] text-muted">{{ itemWeight }} kg</div>
        </div>
      </div>
    </div>

    <!-- Qty / Route -->
    <div class="col-span-2">
      <div class="text-sm font-medium text-primary">{{ shipment.quantity }} <span
          class="text-xs text-muted font-normal">shipped</span></div>
      <div v-if="shipment.shipment_type === 'delivery'" class="text-[10px] text-muted truncate mt-0.5">
        {{ shipment.from_location || '—' }} → {{ shipment.to_location || '—' }}
      </div>
      <div v-else class="text-[10px] text-muted">—</div>
    </div>

    <!-- Financial -->
    <div class="col-span-2">
      <div class="text-xs text-secondary">
        <div v-if="shipment.customer_charge">C: ৳{{ shipment.customer_charge }}</div>
        <div v-if="shipment.company_cost">E: ৳{{ shipment.company_cost }}</div>
        <div v-if="!shipment.customer_charge && !shipment.company_cost" class="text-muted">—</div>
      </div>
    </div>

    <!-- Actions -->
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
            Edit Shipment
          </button>

          <!-- Change Status -->
          <div class="px-3 py-1">
            <p class="text-[10px] text-muted uppercase">Change Status</p>
          </div>
          <button
            v-for="status in availableStatuses"
            :key="status"
            @click="changeStatus(status)"
            :disabled="status === shipment.status"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-surface-alt transition-colors"
            :class="shipment.status === status ? 'text-accent font-semibold cursor-default' : 'text-secondary'"
          >
            <span class="w-2 h-2 rounded-full" :class="getStatusDotClass(status)"></span>
            {{ formatStatus(status) }}
            <span v-if="shipment.status === status" class="ml-auto text-accent">✓</span>
          </button>

          <div class="border-t border-divider my-0.5"></div>

          <!-- Delete -->
          <button @click="openDeleteConfirm"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete
          </button>
        </div>
      </Transition>
      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- Edit Dialog -->
    <BaseDialog v-model="showEditDialog" size="lg" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-info-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-info-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Edit Shipment</h2>
            <p class="text-xs text-muted">#{{ shipment.id }}</p>
          </div>
        </div>

        <form @submit.prevent="handleUpdate" class="space-y-4">
          <!-- Shipment Type & Quantity -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Type</label>
              <div class="flex rounded-lg border border-default overflow-hidden">
                <button type="button" @click="form.shipment_type = 'delivery'"
                  class="flex-1 py-2 text-sm font-medium transition-all duration-200"
                  :class="form.shipment_type === 'delivery' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">
                  Delivery
                </button>
                <button type="button" @click="form.shipment_type = 'pickup'"
                  class="flex-1 py-2 text-sm font-medium transition-all duration-200 border-l border-default"
                  :class="form.shipment_type === 'pickup' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">
                  Pickup
                </button>
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Quantity</label>
              <div class="flex items-center gap-2">
                <button type="button" @click="form.quantity > 1 ? form.quantity-- : null"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                  </svg>
                </button>
                <input v-model.number="form.quantity" type="number" min="1"
                  class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button type="button" @click="form.quantity++"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Delivery Details -->
          <div v-if="form.shipment_type === 'delivery'" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Vehicle Number</label>
                <input v-model="form.vehicle_number" type="text" placeholder="e.g. DHA-1234"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Driver Name</label>
                <input v-model="form.driver_name" type="text" placeholder="Driver full name"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Driver Phone</label>
                <input v-model="form.driver_phone" type="tel" placeholder="+880 1XXX-XXXXXX"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">From Location</label>
                <input v-model="form.from_location" type="text" placeholder="Pickup location"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">To Location</label>
                <input v-model="form.to_location" type="text" placeholder="Delivery location"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Company Name</label>
                <input v-model="form.company_name" type="text" placeholder="Transport company"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Receiver Name</label>
                <input v-model="form.receiver_name" type="text" placeholder="Receiver full name"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Receiver Phone</label>
                <input v-model="form.receiver_phone" type="tel" placeholder="+880 1XXX-XXXXXX"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
          </div>

          <!-- Financial -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Customer Charge</label>
              <input v-model.number="form.customer_charge" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company Cost</label>
              <input v-model.number="form.company_cost" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company Paid</label>
              <input v-model.number="form.company_paid" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Customer Paid</label>
              <input v-model.number="form.customer_paid" type="number" step="0.01" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>

          <!-- Notes -->
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Notes</label>
            <textarea v-model="form.notes" rows="2" placeholder="Add notes..."
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
          </div>

          <!-- Actions -->
          <div class="flex justify-end gap-3 pt-4 border-t border-divider">
            <button type="button" @click="showEditDialog = false"
              class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">
              Cancel
            </button>
            <button type="submit" :disabled="submitting"
              class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
              <span v-if="submitting" class="inline-flex items-center gap-2">
                <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                Saving...
              </span>
              <span v-else>Save Changes</span>
            </button>
          </div>
        </form>
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
            <h3 class="text-lg font-bold text-primary">Delete Shipment</h3>
            <p class="text-xs text-muted">This action cannot be undone.</p>
          </div>
        </div>
        <p class="text-sm text-secondary">
          Are you sure you want to delete shipment <span class="font-semibold text-primary">#{{ shipment.id }}</span>?
          This will restore the item quantity.
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
import { ref, reactive, computed, onMounted } from 'vue'
import type { Shipment, ShipmentStatus, ShipmentType } from '@/types/shipment'
import { useShipmentsStore } from '@/stores/shipments'
import { useItemsStore } from '@/stores/items'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const props = defineProps<{
  shipment: Shipment
}>()

const emit = defineEmits<{
  'shipment-updated': []
  'shipment-deleted': []
  'view-shipment': [shipment: Shipment]
}>()

const shipmentsStore = useShipmentsStore()
const itemsStore = useItemsStore()

const open = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const submitting = ref(false)
const deleting = ref(false)

// Available statuses for shipments
const availableStatuses: ShipmentStatus[] = ['pending', 'in_transit', 'completed', 'cancelled']

const item = computed(() => itemsStore.items.find(i => i.id === props.shipment.item_id))
const itemName = computed(() => item.value?.name || `Item #${props.shipment.item_id}`)
const itemQuantity = computed(() => item.value?.quantity ?? '?')
const itemUnit = computed(() => item.value?.quantity_unit ?? '')
const itemWeight = computed(() => item.value?.weight ?? null)

// Form state
const form = reactive({
  shipment_type: props.shipment.shipment_type as ShipmentType,
  quantity: props.shipment.quantity,
  vehicle_number: props.shipment.vehicle_number || null,
  driver_name: props.shipment.driver_name || null,
  driver_phone: props.shipment.driver_phone || null,
  from_location: props.shipment.from_location || null,
  company_name: props.shipment.company_name || null,
  to_location: props.shipment.to_location || null,
  receiver_name: props.shipment.receiver_name || null,
  receiver_phone: props.shipment.receiver_phone || null,
  customer_charge: props.shipment.customer_charge ?? null,
  company_cost: props.shipment.company_cost ?? null,
  company_paid: props.shipment.company_paid ?? null,
  customer_paid: props.shipment.customer_paid ?? null,
  notes: props.shipment.notes || null,
})

onMounted(async () => {
  if (itemsStore.items.length === 0) await itemsStore.fetchItems()
})

const openEditDialog = () => {
  open.value = false
  // Reset form with current shipment data
  form.shipment_type = props.shipment.shipment_type as ShipmentType
  form.quantity = props.shipment.quantity
  form.vehicle_number = props.shipment.vehicle_number || null
  form.driver_name = props.shipment.driver_name || null
  form.driver_phone = props.shipment.driver_phone || null
  form.from_location = props.shipment.from_location || null
  form.company_name = props.shipment.company_name || null
  form.to_location = props.shipment.to_location || null
  form.receiver_name = props.shipment.receiver_name || null
  form.receiver_phone = props.shipment.receiver_phone || null
  form.customer_charge = props.shipment.customer_charge ?? null
  form.company_cost = props.shipment.company_cost ?? null
  form.company_paid = props.shipment.company_paid ?? null
  form.customer_paid = props.shipment.customer_paid ?? null
  form.notes = props.shipment.notes || null
  showEditDialog.value = true
}

const openDeleteConfirm = () => {
  open.value = false
  showDeleteDialog.value = true
}

const changeStatus = async (status: ShipmentStatus) => {
  if (status === props.shipment.status) return
  try {
    const result = await shipmentsStore.changeShipmentStatus(props.shipment.id, status)
    if (result) {
      push.success(`Status changed to ${formatStatus(status)}`)
      open.value = false
      emit('shipment-updated')
    }
  } catch (error) {
    console.error('Error updating status:', error)
    push.error('Failed to update status')
  }
}

const handleUpdate = async () => {
  // Validate
  if (!form.quantity || form.quantity < 1) {
    push.error('Quantity must be at least 1')
    return
  }

  submitting.value = true
  try {
    // Build the payload
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {
      shipment_type: form.shipment_type,
      quantity: form.quantity,
    }

    // Handle delivery fields based on type
    if (form.shipment_type === 'pickup') {
      // Clear all delivery fields
      payload.vehicle_number = null
      payload.driver_name = null
      payload.driver_phone = null
      payload.from_location = null
      payload.company_name = null
      payload.to_location = null
      payload.receiver_name = null
      payload.receiver_phone = null
    } else {
      // Only include delivery fields if they have values
      payload.vehicle_number = form.vehicle_number || null
      payload.driver_name = form.driver_name || null
      payload.driver_phone = form.driver_phone || null
      payload.from_location = form.from_location || null
      payload.company_name = form.company_name || null
      payload.to_location = form.to_location || null
      payload.receiver_name = form.receiver_name || null
      payload.receiver_phone = form.receiver_phone || null
    }

    // Financial fields - always send as is, allowing null to clear
    payload.customer_charge = form.customer_charge
    payload.company_cost = form.company_cost
    payload.company_paid = form.company_paid
    payload.customer_paid = form.customer_paid

    // Notes
    payload.notes = form.notes || null

    console.log('Sending payload:', payload)

    const result = await shipmentsStore.updateShipment(props.shipment.id, payload)
    if (result) {
      push.success('Shipment updated successfully')
      showEditDialog.value = false
      emit('shipment-updated')
    }
  } catch (error) {
    console.error('Error updating shipment:', error)
    push.error('Failed to update shipment')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async () => {
  deleting.value = true
  try {
    const result = await shipmentsStore.deleteShipment(props.shipment.id)
    if (result) {
      push.success('Shipment deleted successfully')
      showDeleteDialog.value = false
      emit('shipment-deleted')
    }
  } catch (error) {
    console.error('Error deleting shipment:', error)
    push.error('Failed to delete shipment')
  } finally {
    deleting.value = false
  }
}

const getStatusClass = (s: string): string => {
  const classes = {
    pending: 'border-yellow-400 bg-yellow-50 text-yellow-700',
    in_transit: 'border-blue-400 bg-blue-50 text-blue-700',
    completed: 'border-green-400 bg-green-50 text-green-700',
    cancelled: 'border-red-400 bg-red-50 text-red-700'
  }
  return classes[s as keyof typeof classes] || ''
}

const getStatusDotClass = (status: string) => {
  const classes = {
    pending: 'bg-yellow-500',
    in_transit: 'bg-blue-500',
    completed: 'bg-green-500',
    cancelled: 'bg-red-500'
  }
  return classes[status as keyof typeof classes] || 'bg-gray-500'
}

const formatStatus = (status: string): string => {
  const map: Record<string, string> = {
    pending: 'Pending',
    in_transit: 'In Transit',
    completed: 'Completed',
    cancelled: 'Cancelled'
  }
  return map[status] || status
}
</script>
