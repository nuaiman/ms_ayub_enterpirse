<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">New Checkout</h2>
      <p class="text-sm text-secondary">Create a pickup or delivery checkout</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

<!-- Item Selection -->
<div class="pb-6 border-b border-divider">
  <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Item <span class="text-warning-text">*</span></h3>
  
  <!-- Selected Item -->
  <div v-if="selectedItem" class="flex items-start gap-4 p-4 rounded-xl border border-success-border bg-success-bg">
    <div v-if="selectedItem.image_url" class="shrink-0 w-12 h-12 rounded-lg overflow-hidden border border-default">
      <img :src="getImageUrl(selectedItem.image_url) || undefined" :alt="selectedItem.product_name" class="w-full h-full object-cover" />
    </div>
    <div v-else class="shrink-0 w-12 h-12 rounded-lg bg-surface border border-default flex items-center justify-center">
      <span class="text-xs font-semibold text-muted">#{{ selectedItem.id }}</span>
    </div>
    <div class="flex-1 min-w-0">
      <p class="text-sm font-semibold text-primary">{{ selectedItem.product_name }}</p>
      <div class="flex items-center gap-3 mt-1 flex-wrap">
        <span class="text-xs font-medium text-success-text">{{ selectedItem.quantity }} {{ selectedItem.quantity_unit }} available</span>
        <span v-if="selectedItem.customer_phone" class="text-xs text-secondary">{{ selectedItem.customer_phone }}</span>
        <span v-if="selectedItem.lot_number" class="text-xs text-info-text">Lot: {{ selectedItem.lot_number }}</span>
        <span v-if="selectedItem.storage_name" class="text-xs text-muted">{{ selectedItem.storage_name }}</span>
      </div>
    </div>
    <button @click="clearSelectedItem" type="button" class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
    </button>
  </div>

  <!-- Search -->
  <div v-else class="relative">
    <div class="relative">
      <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
      <input v-model="itemSearch" type="text" placeholder="Search by product name, phone, lot, storage..." 
        @focus="showItemDropdown = true" @keydown.escape="showItemDropdown = false"
        class="input w-full pl-9 pr-8 py-2.5 rounded-lg text-sm placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
      <button v-if="itemSearch" @click="itemSearch = ''" type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
      </button>
    </div>

    <!-- Dropdown -->
    <div v-if="showItemDropdown" class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
      <div class="max-h-64 overflow-y-auto">
        <div v-if="filteredItems.length === 0" class="px-4 py-8 text-center">
          <div class="w-12 h-12 mx-auto rounded-full bg-surface-alt flex items-center justify-center mb-3">
            <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
          </div>
          <p class="text-sm font-medium text-primary">No items found</p>
          <p class="text-xs text-muted mt-1">Try a different search term</p>
        </div>
        <button v-for="i in filteredItems" :key="i.id" @click="selectItem(i)" type="button"
          class="w-full flex items-start gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
          <!-- Item Image -->
          <div class="shrink-0 mt-0.5">
            <div v-if="i.image_url" class="w-10 h-10 rounded-lg overflow-hidden border border-default">
              <img :src="getImageUrl(i.image_url) || undefined" :alt="i.product_name" class="w-full h-full object-cover" />
            </div>
            <div v-else class="w-10 h-10 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-[10px] font-semibold text-muted">#{{ i.id }}</span>
            </div>
          </div>
          
          <!-- Item Details -->
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-primary truncate">{{ i.product_name }}</p>
            <div class="flex items-center gap-2 mt-1 flex-wrap">
              <span class="text-xs font-medium text-success-text">{{ i.quantity }} {{ i.quantity_unit }}</span>
              <span v-if="i.customer_phone" class="text-xs text-secondary">{{ i.customer_phone }}</span>
            </div>
            <div class="flex items-center gap-2 mt-0.5 flex-wrap">
              <span v-if="i.lot_number" class="text-[11px] text-info-text font-medium">Lot: {{ i.lot_number }}</span>
              <span v-if="i.storage_name" class="text-[11px] text-muted">{{ i.storage_name }}</span>
              <span v-if="i.account_name" class="text-[11px] text-muted">{{ i.account_name }}</span>
            </div>
            <div v-if="i.category" class="flex items-center gap-1 mt-0.5">
              <span class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded">{{ i.category }}</span>
              <span v-if="i.subcategory" class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded">{{ i.subcategory }}</span>
            </div>
          </div>
          
          <!-- Arrow -->
          <svg class="w-4 h-4 text-muted shrink-0 mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
        </button>
      </div>
      
      <!-- Footer -->
      <div class="px-4 py-2 border-t border-divider bg-surface-alt">
        <p class="text-[11px] text-muted">{{ filteredItems.length }} item{{ filteredItems.length !== 1 ? 's' : '' }} in stock</p>
      </div>
    </div>
    
    <!-- Backdrop -->
    <div v-if="showItemDropdown" class="fixed inset-0 z-40" @click="showItemDropdown = false"></div>
  </div>
</div>

<!-- Checkout Details -->
<div class="pb-6 border-b border-divider">
  <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Checkout Details</h3>
  
  <!-- Type - Full Width -->
  <div class="space-y-1.5 mb-4">
    <label class="text-sm font-medium text-primary">Type <span class="text-warning-text">*</span></label>
    <div class="flex rounded-lg border border-default overflow-hidden">
      <button type="button" @click="checkoutType = 'pickup'" class="flex-1 py-2 text-sm font-medium transition-all duration-200" :class="checkoutType === 'pickup' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">Pickup</button>
      <button type="button" @click="checkoutType = 'delivery'" class="flex-1 py-2 text-sm font-medium transition-all duration-200 border-l border-default" :class="checkoutType === 'delivery' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">Delivery</button>
    </div>
  </div>

  <!-- Quantity & Weight -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
    <div class="space-y-1.5">
      <label class="text-sm font-medium text-primary">Quantity <span class="text-warning-text">*</span></label>
      <div class="flex items-center gap-2">
        <button type="button" @click="quantity > 1 ? quantity-- : null" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" /></svg></button>
        <input v-model.number="quantity" type="number" min="1" required class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
        <button type="button" @click="selectedItem && quantity < selectedItem.quantity ? quantity++ : null" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" /></svg></button>
      </div>
    </div>
    <div class="space-y-1.5">
      <label class="text-sm font-medium text-primary">Weight (kg) <span class="text-warning-text">*</span></label>
      <input v-model.number="weight" type="number" step="0.001" min="0" placeholder="0.000" required
        class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
    </div>
  </div>

  <!-- Receiver Info -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div class="space-y-1.5">
      <label class="text-sm font-medium text-primary">Receiver Name</label>
      <input v-model="receiver_name" type="text" placeholder="Receiver name"
        class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
    </div>
    <div class="space-y-1.5">
      <label class="text-sm font-medium text-primary">Receiver Phone</label>
      <input v-model="receiver_phone" type="tel" placeholder="Receiver phone"
        class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
    </div>
  </div>
</div>

        <!-- Delivery Details -->
        <div v-if="checkoutType === 'delivery'" class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Delivery Details <span class="text-warning-text">*</span></h3>
          


          <!-- Location -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">From Location</label><input v-model="from_location" type="text" placeholder="Pickup location" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">To Location <span class="text-warning-text">*</span></label><input v-model="to_location" type="text" placeholder="Delivery location" required class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
          </div>

          <!-- Transport -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Vehicle</label><input v-model="vehicle_number" type="text" placeholder="Vehicle number" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Driver</label><input v-model="driver_name" type="text" placeholder="Driver name" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Driver Phone</label><input v-model="driver_phone" type="tel" placeholder="Driver phone" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
          </div>

          <!-- Financial -->
          <div class="grid grid-cols-3 gap-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Delivery Charge</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="delivery_charge" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Delivery Cost</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="delivery_cost" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Customer Paid</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="customer_paid" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
          </div>
        </div>

        <!-- Image & Notes -->
        <div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Image <span class="normal-case text-xs font-normal">(Optional)</span></h3>
              <div class="flex items-start gap-4">
                <div class="shrink-0">
                  <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default"><img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" /><button type="button" @click="removeImage" class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">×</button></div>
                  <div v-else class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt"><svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>
                </div>
                <div class="flex-1 space-y-3">
                  <label class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2" :class="{ 'opacity-50 cursor-not-allowed': uploading }"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>{{ uploading ? 'Uploading...' : 'Choose File' }}<input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" /></label>
                </div>
              </div>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Notes <span class="normal-case text-xs font-normal">(Optional)</span></h3>
              <textarea v-model="notes" rows="4" placeholder="Add notes..." class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm" class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit" :disabled="isSubmitDisabled" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2"><svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>Creating...</span>
            <span v-else>Create Checkout</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { CheckoutType } from '@/types/checkout'
import type { Item } from '@/types/item'
import { useCheckoutsStore } from '@/stores/checkout'
import { useItemsStore } from '@/stores/items'
import { push } from 'notivue'
import { getImageUrl, uploadImage } from '@/utils/image'

const emit = defineEmits<{ 'checkout-created': [] }>()
const checkoutsStore = useCheckoutsStore()
const itemsStore = useItemsStore()

const submitting = ref(false)
const checkoutType = ref<CheckoutType>('pickup')
const quantity = ref(1)
const weight = ref<number | null>(null)

// Receiver (only for delivery)
const receiver_name = ref('')
const receiver_phone = ref('')

// Delivery fields
const from_location = ref('')
const to_location = ref('')
const vehicle_number = ref('')
const driver_name = ref('')
const driver_phone = ref('')

// Financial
const delivery_charge = ref<number | null>(null)
const delivery_cost = ref<number | null>(null)
const customer_paid = ref<number | null>(null)
const notes = ref('')

// Image
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)

// Item search
const itemSearch = ref('')
const showItemDropdown = ref(false)
const selectedItem = ref<Item | null>(null)

const isSubmitDisabled = computed(() => {
  if (submitting.value) return true
  if (!selectedItem.value) return true
  if (!quantity.value || quantity.value <= 0) return true
  if (!weight.value || weight.value <= 0) return true
  if (quantity.value > (selectedItem.value?.quantity || 0)) return true
  if (checkoutType.value === 'delivery' && !to_location.value) return true
  return false
})

const filteredItems = computed(() => {
  const inStock = itemsStore.items.filter(i => i.quantity > 0)
  if (!itemSearch.value) return inStock
  const q = itemSearch.value.toLowerCase()
  return inStock.filter(i =>
    i.product_name.toLowerCase().includes(q) ||
    (i.customer_phone && i.customer_phone.includes(q)) ||
    (i.lot_number && i.lot_number.toLowerCase().includes(q)) ||
    (i.storage_name && i.storage_name.toLowerCase().includes(q)) ||
    (i.category && i.category.toLowerCase().includes(q))
  )
})

onMounted(async () => { await itemsStore.fetchItems() })

const selectItem = (item: Item) => { selectedItem.value = item; itemSearch.value = ''; showItemDropdown.value = false }
const clearSelectedItem = () => { selectedItem.value = null; itemSearch.value = ''; quantity.value = 1 }

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement; const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }
  if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }
  imageFile.value = file
  const reader = new FileReader(); reader.onload = (e) => { imagePreview.value = e.target?.result as string }; reader.readAsDataURL(file)
}

const removeImage = () => { imageFile.value = null; imagePreview.value = null }

const resetForm = () => {
  checkoutType.value = 'pickup'; quantity.value = 1; weight.value = null; clearSelectedItem()
  receiver_name.value = ''; receiver_phone.value = ''
  from_location.value = ''; to_location.value = ''
  vehicle_number.value = ''; driver_name.value = ''; driver_phone.value = ''
  delivery_charge.value = null; delivery_cost.value = null; customer_paid.value = null
  notes.value = ''; removeImage()
}

const submit = async () => {
  if (!selectedItem.value || !quantity.value || !weight.value) return
  if (quantity.value <= 0 || weight.value <= 0) { push.error('Quantity and weight are required'); return }
  if (quantity.value > selectedItem.value.quantity) { push.error('Insufficient quantity'); return }
  if (checkoutType.value === 'delivery' && !to_location.value) { push.error('Delivery location is required'); return }
  submitting.value = true
  try {
    const newCheckout = await checkoutsStore.createCheckout({
      type: checkoutType.value, item_id: selectedItem.value.id, quantity: quantity.value,
      weight: weight.value,
      receiver_name: receiver_name.value || null, receiver_phone: receiver_phone.value || null,
      from_location: from_location.value || null, to_location: to_location.value || null,
      vehicle_number: vehicle_number.value || null, driver_name: driver_name.value || null, driver_phone: driver_phone.value || null,
      delivery_charge: delivery_charge.value, delivery_cost: delivery_cost.value, customer_paid: customer_paid.value,
      notes: notes.value || null,
    })
    if (!newCheckout) { submitting.value = false; return }
    if (imageFile.value) { uploading.value = true; await uploadImage('checkouts', newCheckout.id, imageFile.value); uploading.value = false }
    await itemsStore.fetchItems()
    push.success('Checkout created successfully!'); resetForm(); emit('checkout-created')
  } catch { push.error('Failed to create checkout') }
  finally { submitting.value = false }
}
</script>