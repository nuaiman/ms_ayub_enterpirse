<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-shipment', shipment)">

    <div class="col-span-1"><span class="text-sm font-medium text-primary">#{{ shipment.id }}</span></div>

    <div class="col-span-2">
      <div class="flex items-center gap-1.5">
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="shipment.shipment_type === 'pickup' ? 'status-info' : 'status-success'">{{ shipment.shipment_type
          }}</span>
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
          :class="getStatusClass(shipment.status)">{{ shipment.status }}</span>
      </div>
    </div>

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

    <div class="col-span-2">
      <div class="text-sm font-medium text-primary">{{ shipment.quantity }} <span
          class="text-xs text-muted font-normal">shipped</span></div>
      <div v-if="shipment.shipment_type === 'delivery'" class="text-[10px] text-muted truncate mt-0.5">
        {{ shipment.from_location || '—' }} → {{ shipment.to_location || '—' }}
      </div>
      <div v-else class="text-[10px] text-muted">—</div>
    </div>

    <div class="col-span-2">
      <div class="text-xs text-secondary">
        <div v-if="shipment.customer_charge">C: ${{ shipment.customer_charge }}</div>
        <div v-if="shipment.company_cost">E: ${{ shipment.company_cost }}</div>
        <div v-if="!shipment.customer_charge && !shipment.company_cost" class="text-muted">—</div>
      </div>
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
          class="absolute right-0 top-9 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider">
            <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
          </div>
          <button v-if="shipment.status === 'pending'" @click="changeStatus('in_transit')"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-info-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg> Mark In Transit
          </button>
          <button v-if="shipment.status === 'in_transit'" @click="changeStatus('completed')"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-success-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg> Mark Completed
          </button>
          <button v-if="shipment.status === 'pending'" @click="changeStatus('cancelled')"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg> Cancel
          </button>
          <div class="border-t border-divider my-0.5"></div>
          <button @click="showDeleteConfirm = true; open = false"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg> Delete
          </button>
        </div>
      </Transition>
      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

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
            <h2 class="text-lg font-semibold text-primary">Delete Shipment</h2>
            <p class="text-xs text-muted">#{{ shipment.id }}</p>
          </div>
        </div>
        <p class="text-sm text-secondary">Are you sure? This will restore the item quantity.</p>
        <div class="flex gap-2 justify-end pt-2">
          <button @click="showDeleteConfirm = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="deleteShipment"
            class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg">Delete</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Shipment, ShipmentStatus } from '@/types/shipment'
import { useShipmentsStore } from '@/stores/shipments'
import { useItemsStore } from '@/stores/items'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const props = defineProps<{ shipment: Shipment }>()
const emit = defineEmits<{ 'shipment-updated': []; 'view-shipment': [shipment: Shipment] }>()

const shipmentsStore = useShipmentsStore()
const itemsStore = useItemsStore()
const open = ref(false)
const showDeleteConfirm = ref(false)

const item = computed(() => itemsStore.items.find(i => i.id === props.shipment.item_id))
const itemName = computed(() => item.value?.name || `Item #${props.shipment.item_id}`)
const itemQuantity = computed(() => item.value?.quantity ?? '?')
const itemUnit = computed(() => item.value?.quantity_unit ?? '')
const itemWeight = computed(() => item.value?.weight ?? null)

onMounted(async () => {
  if (itemsStore.items.length === 0) await itemsStore.fetchItems()
})

const changeStatus = async (status: ShipmentStatus) => {
  await shipmentsStore.changeShipmentStatus(props.shipment.id, status)
  open.value = false
}

const deleteShipment = async () => {
  try { const s = await shipmentsStore.deleteShipment(props.shipment.id); if (s) showDeleteConfirm.value = false }
  catch { push.error('Failed to delete shipment') }
}

const getStatusClass = (s: string): string => ({ pending: 'status-warning', in_transit: 'status-info', completed: 'status-success', cancelled: 'status-warning' }[s] || '')
</script>
