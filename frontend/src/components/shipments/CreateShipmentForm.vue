<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">New Shipment</h2>
      <p class="text-sm text-secondary">Create a pickup or delivery shipment</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Item Selection - TOP -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Item <span class="text-warning-text">*</span>
          </h3>

          <!-- Selected Item -->
          <div v-if="selectedItem"
            class="flex items-start gap-4 p-4 rounded-xl border border-success-border bg-success-bg">
            <div v-if="selectedItem.image_url"
              class="shrink-0 w-14 h-14 rounded-lg overflow-hidden border border-default">
              <img :src="getImageUrl(selectedItem.image_url) || undefined" :alt="selectedItem.name"
                class="w-full h-full object-cover" />
            </div>
            <div v-else
              class="shrink-0 w-14 h-14 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
              <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-primary">{{ selectedItem.name }}</p>
              <div class="flex items-center gap-3 mt-1 flex-wrap">
                <span class="text-xs text-secondary">{{ selectedItem.quantity }} {{ selectedItem.quantity_unit }} available</span>
                <span v-if="selectedItem.weight" class="text-xs text-muted">{{ selectedItem.weight }} kg</span>
                <span v-if="selectedItem.customer_id" class="text-xs text-muted">Customer #{{ selectedItem.customer_id }}</span>
                <span v-if="selectedItem.duration_type" class="text-xs text-accent">· Contract: {{ selectedItem.duration_type }}</span>
              </div>
            </div>
            <button @click="clearSelectedItem" type="button"
              class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Search (when no item selected) -->
          <div v-else class="relative">
            <div class="relative">
              <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none" fill="none"
                stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <input ref="itemSearchRef" v-model="itemSearch" type="text"
                placeholder="Search items by name, customer..." @focus="showItemDropdown = true"
                @keydown.escape="showItemDropdown = false"
                class="input w-full pl-9 pr-8 py-2.5 rounded-lg text-sm placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button v-if="itemSearch" @click="itemSearch = ''" type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Dropdown -->
            <div v-if="showItemDropdown"
              class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
              <div class="max-h-56 overflow-y-auto">
                <div v-if="filteredItems.length === 0" class="px-4 py-6 text-center">
                  <div class="w-10 h-10 mx-auto rounded-full bg-surface-alt flex items-center justify-center mb-2">
                    <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                  </div>
                  <p class="text-sm text-muted">No items match "{{ itemSearch }}"</p>
                </div>
                <button v-for="i in filteredItems" :key="i.id" @click="selectItem(i)" type="button"
                  class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                  <div class="shrink-0">
                    <div v-if="i.image_url" class="w-9 h-9 rounded-lg overflow-hidden border border-default">
                      <img :src="getImageUrl(i.image_url) || undefined" :alt="i.name"
                        class="w-full h-full object-cover" />
                    </div>
                    <div v-else
                      class="w-9 h-9 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
                      <svg class="w-4 h-4 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                      </svg>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-primary truncate">{{ i.name }}</p>
                    <div class="flex items-center gap-2 mt-0.5">
                      <span class="text-xs text-muted">{{ i.quantity }} {{ i.quantity_unit }}</span>
                      <span v-if="i.customer_id" class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded">#{{ i.customer_id }}</span>
                      <span v-if="i.weight" class="text-[10px] text-muted">{{ i.weight }}kg</span>
                      <span v-if="i.duration_type" class="text-[10px] text-accent">contract</span>
                    </div>
                  </div>
                  <svg class="w-4 h-4 text-muted shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </div>

              <!-- Footer with count -->
              <div class="px-4 py-2 border-t border-divider bg-surface-alt">
                <p class="text-[11px] text-muted">{{ filteredItems.length }} item{{ filteredItems.length !== 1 ? 's' :
                  '' }} in stock</p>
              </div>
            </div>

            <!-- Backdrop -->
            <div v-if="showItemDropdown" class="fixed inset-0 z-40" @click="showItemDropdown = false"></div>
          </div>
        </div>

        <!-- Shipment Type & Quantity -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Shipment Details</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Type <span class="text-warning-text">*</span></label>
              <div class="flex rounded-lg border border-default overflow-hidden">
                <button type="button" @click="shipment_type = 'delivery'"
                  class="flex-1 py-2 text-sm font-medium transition-all duration-200"
                  :class="shipment_type === 'delivery' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">
                  Delivery
                </button>
                <button type="button" @click="shipment_type = 'pickup'"
                  class="flex-1 py-2 text-sm font-medium transition-all duration-200 border-l border-default"
                  :class="shipment_type === 'pickup' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">
                  Pickup
                </button>
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Quantity <span class="text-warning-text">*</span></label>
              <div class="flex items-center gap-2">
                <button type="button" @click="quantity > 1 ? quantity-- : null"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                  </svg>
                </button>
                <input v-model.number="quantity" type="number" min="1" required
                  class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button type="button" @click="selectedItem && quantity < selectedItem.quantity ? quantity++ : null"
                  class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors"
                  :class="{ 'opacity-50 cursor-not-allowed': selectedItem && quantity >= selectedItem.quantity }">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                  </svg>
                </button>
              </div>
              <p v-if="selectedItem" class="text-xs"
                :class="quantity > selectedItem.quantity ? 'text-warning-text' : 'text-muted'">
                {{ quantity > selectedItem.quantity ? 'Exceeds available!' : selectedItem.quantity + ' ' +
                  selectedItem.quantity_unit + ' available' }}
              </p>
            </div>
          </div>
        </div>

        <!-- Delivery Fields (shown only for delivery) -->
        <div v-if="shipment_type === 'delivery'" class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Delivery Details</h3>

          <!-- Vehicle & Driver -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Vehicle Number</label>
              <input v-model="vehicle_number" type="text" placeholder="e.g. DHA-1234"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Driver Name</label>
              <input v-model="driver_name" type="text" placeholder="Driver full name"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Driver Phone</label>
              <input v-model="driver_phone" type="tel" placeholder="+880 1XXX-XXXXXX"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>

          <!-- Route -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">From Location</label>
              <input v-model="from_location" type="text" placeholder="Pickup location"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">To Location</label>
              <input v-model="to_location" type="text" placeholder="Delivery location"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>

          <!-- Company & Receiver -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company Name</label>
              <input v-model="company_name" type="text" placeholder="Transport company"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Receiver Name</label>
              <input v-model="receiver_name" type="text" placeholder="Receiver full name"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Receiver Phone</label>
              <input v-model="receiver_phone" type="tel" placeholder="+880 1XXX-XXXXXX"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
        </div>

        <!-- Financial -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Financial <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Customer Charge</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
                <input v-model.number="customer_charge" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company Cost</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
                <input v-model.number="company_cost" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company Paid</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
                <input v-model.number="company_paid" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Customer Paid</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
                <input v-model.number="customer_paid" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
          </div>
        </div>

        <!-- Notes -->
        <div>
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Notes <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <textarea v-model="notes" rows="2" placeholder="Add shipment notes or special instructions..."
            class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">
            Clear Form
          </button>
          <button type="submit"
            :disabled="submitting || !selectedItem || quantity > (selectedItem?.quantity || 0) || quantity <= 0"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Creating...
            </span>
            <span v-else>Create Shipment</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { ShipmentType } from '@/types/shipment'
import type { Item } from '@/types/item'
import { useShipmentsStore } from '@/stores/shipments'
import { useItemsStore } from '@/stores/items'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{ 'shipment-created': [] }>()

const shipmentsStore = useShipmentsStore()
const itemsStore = useItemsStore()

const submitting = ref(false)
const shipment_type = ref<ShipmentType>('delivery')
const quantity = ref(1)

// Delivery fields
const vehicle_number = ref('')
const driver_name = ref('')
const driver_phone = ref('')
const from_location = ref('')
const company_name = ref('')
const to_location = ref('')
const receiver_name = ref('')
const receiver_phone = ref('')

// Financial
const customer_charge = ref<number | null>(null)
const company_cost = ref<number | null>(null)
const company_paid = ref<number | null>(null)
const customer_paid = ref<number | null>(null)
const notes = ref('')

// Item search
const itemSearchRef = ref<HTMLInputElement | null>(null)
const itemSearch = ref('')
const showItemDropdown = ref(false)
const selectedItem = ref<Item | null>(null)
const items = ref<Item[]>([])

onMounted(async () => {
  await itemsStore.fetchItems()
  items.value = itemsStore.items
})

const filteredItems = computed(() => {
  const inStock = items.value.filter(i => i.quantity > 0)
  if (!itemSearch.value) return inStock
  const q = itemSearch.value.toLowerCase()
  return inStock.filter(i =>
    i.name.toLowerCase().includes(q) ||
    (i.notes && i.notes.toLowerCase().includes(q)) ||
    (i.customer_id && String(i.customer_id).includes(q)) ||
    (i.duration_type && i.duration_type.toLowerCase().includes(q))
  )
})

watch(shipment_type, () => {
  if (shipment_type.value === 'pickup') {
    vehicle_number.value = ''
    driver_name.value = ''
    driver_phone.value = ''
    from_location.value = ''
    company_name.value = ''
    to_location.value = ''
    receiver_name.value = ''
    receiver_phone.value = ''
  }
})

const selectItem = (item: Item) => {
  selectedItem.value = item
  itemSearch.value = ''
  showItemDropdown.value = false
}

const clearSelectedItem = () => {
  selectedItem.value = null
  itemSearch.value = ''
  quantity.value = 1
}

const resetForm = () => {
  shipment_type.value = 'delivery'
  quantity.value = 1
  clearSelectedItem()
  vehicle_number.value = ''
  driver_name.value = ''
  driver_phone.value = ''
  from_location.value = ''
  company_name.value = ''
  to_location.value = ''
  receiver_name.value = ''
  receiver_phone.value = ''
  customer_charge.value = null
  company_cost.value = null
  company_paid.value = null
  customer_paid.value = null
  notes.value = ''
}

const submit = async () => {
  if (!selectedItem.value || quantity.value <= 0) return
  if (quantity.value > selectedItem.value.quantity) {
    push.error('Insufficient quantity')
    return
  }

  submitting.value = true
  try {
    const newShipment = await shipmentsStore.createShipment({
      shipment_type: shipment_type.value,
      item_id: selectedItem.value.id,
      quantity: quantity.value,
      vehicle_number: shipment_type.value === 'delivery' ? (vehicle_number.value || null) : null,
      driver_name: shipment_type.value === 'delivery' ? (driver_name.value || null) : null,
      driver_phone: shipment_type.value === 'delivery' ? (driver_phone.value || null) : null,
      from_location: shipment_type.value === 'delivery' ? (from_location.value || null) : null,
      company_name: shipment_type.value === 'delivery' ? (company_name.value || null) : null,
      to_location: shipment_type.value === 'delivery' ? (to_location.value || null) : null,
      receiver_name: shipment_type.value === 'delivery' ? (receiver_name.value || null) : null,
      receiver_phone: shipment_type.value === 'delivery' ? (receiver_phone.value || null) : null,
      customer_charge: customer_charge.value,
      company_cost: company_cost.value,
      company_paid: company_paid.value,
      customer_paid: customer_paid.value,
      notes: notes.value || null,
    })

    if (newShipment) {
      await itemsStore.fetchItems()
      items.value = itemsStore.items
      push.success('Shipment created successfully!')
      resetForm()
      emit('shipment-created')
    }
  } catch (error) {
    console.error('Error creating shipment:', error)
    push.error('Failed to create shipment')
  } finally {
    submitting.value = false
  }
}
</script>
