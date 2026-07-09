<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-checkout', checkout)">

    <!-- ID -->
    <div class="col-span-1">
      <span class="text-sm font-medium text-primary">#{{ checkout.id }}</span>
    </div>

    <!-- Type & Status -->
    <div class="col-span-2">
      <div class="flex items-center gap-1.5">
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
          :class="checkout.type === 'pickup' ? 'status-info' : 'status-success'">
          {{ checkout.type }}
        </span>
        <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
          :class="getStatusClass(checkout.status)">
          {{ checkout.status }}
        </span>
      </div>
    </div>

    <!-- Item / Customer -->
    <div class="col-span-2">
      <div class="flex items-center gap-2">
        <div class="shrink-0">
          <div v-if="checkout.image_url" class="w-8 h-8 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(checkout.image_url) || undefined" alt="Checkout" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-8 h-8 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
            <span class="text-[10px] font-semibold text-muted">#{{ checkout.item_id }}</span>
          </div>
        </div>
        <div class="min-w-0">
          <div class="text-xs font-medium text-primary truncate">{{ itemName }}</div>
          <div class="text-[10px] text-muted truncate">{{ item?.customer_phone || '—' }}</div>
        </div>
      </div>
    </div>

    <!-- Quantity -->
    <div class="col-span-1">
      <span class="text-sm font-semibold text-primary">{{ checkout.quantity }}</span>
      <span class="text-xs text-muted ml-1">{{ item?.quantity_unit || '' }}</span>
    </div>

    <!-- Route / Receiver -->
    <div class="col-span-2">
      <div v-if="checkout.type === 'delivery'" class="text-xs text-secondary">
        <div class="truncate">{{ checkout.from_location || '—' }} → {{ checkout.to_location || '—' }}</div>
      </div>
      <div v-else class="text-xs text-muted">—</div>
      <div v-if="checkout.receiver_name" class="text-[10px] text-muted truncate mt-0.5">
        {{ checkout.receiver_name }}
      </div>
    </div>

    <!-- Financial -->
    <div class="col-span-2">
      <div class="text-xs text-secondary">
        <div v-if="checkout.type === 'delivery' && checkout.delivery_charge">Charge: ৳{{ checkout.delivery_charge }}</div>
        <div v-if="checkout.type === 'delivery' && checkout.delivery_cost">Cost: ৳{{ checkout.delivery_cost }}</div>
        <div v-if="checkout.type === 'delivery' && checkout.customer_paid">Paid: ৳{{ checkout.customer_paid }}</div>
        <div v-if="!checkout.delivery_charge && !checkout.delivery_cost && !checkout.customer_paid" class="text-muted">—</div>
      </div>
    </div>

    <!-- Actions -->
    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <button @click="open = !open"
        class="w-7 h-7 flex items-center justify-center border border-default rounded-md hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100"
        :class="{ 'opacity-100': open }">
        <svg class="w-3.5 h-3.5 text-secondary" fill="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" /><circle cx="12" cy="12" r="1.5" /><circle cx="12" cy="19" r="1.5" />
        </svg>
      </button>

      <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95 translate-y-1"
        enter-to-class="opacity-100 scale-100 translate-y-0" leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 scale-100 translate-y-0" leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="open" class="absolute right-0 top-9 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider"><p class="text-[11px] font-semibold text-muted uppercase">Actions</p></div>
          <button @click="openEditDialog" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>Edit Checkout
          </button>
          <div class="border-t border-divider my-0.5"></div>
          <div class="px-3 py-1"><p class="text-[10px] text-muted uppercase">Change Status</p></div>
          <button @click="changeStatus('pending')" :disabled="checkout.status === 'pending'" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-surface-alt transition-colors" :class="checkout.status === 'pending' ? 'text-accent font-semibold' : 'text-secondary'">
            <span class="w-2 h-2 rounded-full bg-yellow-500"></span> Pending <span v-if="checkout.status === 'pending'" class="ml-auto text-accent">✓</span>
          </button>
          <button @click="changeStatus('in_transit')" :disabled="checkout.status === 'in_transit'" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-surface-alt transition-colors" :class="checkout.status === 'in_transit' ? 'text-accent font-semibold' : 'text-secondary'">
            <span class="w-2 h-2 rounded-full bg-blue-500"></span> In Transit <span v-if="checkout.status === 'in_transit'" class="ml-auto text-accent">✓</span>
          </button>
          <button @click="changeStatus('complete')" :disabled="checkout.status === 'complete'" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs hover:bg-surface-alt transition-colors" :class="checkout.status === 'complete' ? 'text-accent font-semibold' : 'text-secondary'">
            <span class="w-2 h-2 rounded-full bg-green-500"></span> Complete <span v-if="checkout.status === 'complete'" class="ml-auto text-accent">✓</span>
          </button>
          <div class="border-t border-divider my-0.5"></div>
          <button @click="openDeleteConfirm" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>Delete
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
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div><h2 class="text-lg font-semibold text-primary">Edit Checkout</h2><p class="text-xs text-muted">#{{ checkout.id }}</p></div>
        </div>

        <!-- Image Edit -->
        <div class="flex items-center gap-4 pb-4 border-b border-divider">
          <div class="shrink-0 relative group/edit-img">
            <div v-if="editImagePreview || checkout.image_url" class="w-16 h-16 rounded-lg overflow-hidden border border-default">
              <img :src="editImagePreview || getImageUrl(checkout.image_url) || undefined" alt="Checkout" class="w-full h-full object-cover" />
            </div>
            <div v-else class="w-16 h-16 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
              <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <button @click="triggerEditImageInput" class="absolute inset-0 rounded-lg bg-black/50 flex items-center justify-center opacity-0 group-hover/edit-img:opacity-100 transition-opacity cursor-pointer">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>
            <input ref="editImageInputRef" type="file" accept="image/*" class="hidden" @change="handleEditImageSelect" />
          </div>
          <div>
            <p class="text-xs text-muted">Click on image to change</p>
            <button v-if="checkout.image_url && !editImageFile" @click="removeExistingImage" type="button" class="text-xs text-warning-text hover:underline mt-1">Remove current image</button>
            <button v-if="editImageFile" @click="clearEditImage" type="button" class="text-xs text-muted hover:underline mt-1">Clear selection</button>
          </div>
        </div>

        <form @submit.prevent="handleUpdate" class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Type</label>
            <div class="flex rounded-lg border border-default overflow-hidden">
              <button type="button" @click="editForm.type = 'pickup'" class="flex-1 py-2 text-sm font-medium transition-all duration-200" :class="editForm.type === 'pickup' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">Pickup</button>
              <button type="button" @click="editForm.type = 'delivery'" class="flex-1 py-2 text-sm font-medium transition-all duration-200 border-l border-default" :class="editForm.type === 'delivery' ? 'bg-accent text-accent-foreground' : 'bg-surface text-secondary hover:bg-surface-alt'">Delivery</button>
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Quantity</label>
            <div class="flex items-center gap-2">
              <button type="button" @click="editForm.quantity > 1 ? editForm.quantity-- : null" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt shrink-0">−</button>
              <input v-model.number="editForm.quantity" type="number" min="1" class="input flex-1 px-3 py-2 rounded-lg text-center focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              <button type="button" @click="editForm.quantity++" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt shrink-0">+</button>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Receiver Name</label><input v-model="editForm.receiver_name" type="text" placeholder="Receiver name" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Receiver Phone</label><input v-model="editForm.receiver_phone" type="tel" placeholder="Receiver phone" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
          </div>

          <template v-if="editForm.type === 'delivery'">
            <div class="border-t border-divider pt-4">
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-3">Delivery Details</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">From Location</label><input v-model="editForm.from_location" type="text" placeholder="Pickup location" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">To Location</label><input v-model="editForm.to_location" type="text" placeholder="Delivery location" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Vehicle</label><input v-model="editForm.vehicle_number" type="text" placeholder="Vehicle number" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Driver</label><input v-model="editForm.driver_name" type="text" placeholder="Driver name" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Driver Phone</label><input v-model="editForm.driver_phone" type="tel" placeholder="Driver phone" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
              </div>
              <div class="grid grid-cols-3 gap-4">
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Delivery Charge</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="editForm.delivery_charge" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Delivery Cost</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="editForm.delivery_cost" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
                <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Customer Paid</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="editForm.customer_paid" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
              </div>
            </div>
          </template>

          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Notes</label><textarea v-model="editForm.notes" rows="2" placeholder="Add notes..." class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea></div>

          <div class="flex justify-end gap-3 pt-4 border-t border-divider">
            <button type="button" @click="showEditDialog = false" class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Cancel</button>
            <button type="submit" :disabled="saving" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
              <span v-if="saving" class="inline-flex items-center gap-2"><svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>Saving...</span>
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
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
          </div>
          <div><h3 class="text-lg font-bold text-primary">Delete Checkout</h3><p class="text-xs text-muted">This action cannot be undone.</p></div>
        </div>
        <p class="text-sm text-secondary">Are you sure you want to delete checkout <span class="font-semibold text-primary">#{{ checkout.id }}</span>?</p>
        <div class="p-4 rounded-xl border border-default bg-surface-alt">
          <label class="flex items-center gap-3 cursor-pointer">
            <input v-model="restoreQuantity" type="checkbox" class="w-4 h-4 rounded border-default text-accent focus:ring-accent" />
            <div><p class="text-sm font-medium text-primary">Restore item quantity</p><p class="text-xs text-muted mt-0.5">Return {{ checkout.quantity }} {{ item?.quantity_unit || 'units' }} back to "{{ itemName }}"</p></div>
          </label>
        </div>
        <div class="flex justify-end gap-3 pt-2">
          <button type="button" @click="showDeleteDialog = false" class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Cancel</button>
          <button type="button" @click="handleDelete" :disabled="deleting" class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2">
            <span v-if="deleting"><svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>Deleting...</span>
            <span v-else>Delete</span>
          </button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import type { Checkout, CheckoutStatus } from '@/types/checkout'
import { useCheckoutsStore } from '@/stores/checkout'
import { useItemsStore } from '@/stores/items'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'
import { getImageUrl, uploadImage, deleteImage } from '@/utils/image'

const props = defineProps<{ checkout: Checkout }>()
const emit = defineEmits<{ 'checkout-updated': []; 'checkout-deleted': []; 'view-checkout': [checkout: Checkout] }>()

const checkoutsStore = useCheckoutsStore()
const itemsStore = useItemsStore()

const open = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const saving = ref(false)
const deleting = ref(false)
const restoreQuantity = ref(true)

const editImageInputRef = ref<HTMLInputElement | null>(null)
const editImageFile = ref<File | null>(null)
const editImagePreview = ref<string | null>(null)

const item = computed(() => itemsStore.items.find(i => i.id === props.checkout.item_id))
const itemName = computed(() => item.value?.name || `Item #${props.checkout.item_id}`)

const editForm = reactive({
  type: props.checkout.type as 'pickup' | 'delivery',
  quantity: props.checkout.quantity,
  receiver_name: props.checkout.receiver_name || '',
  receiver_phone: props.checkout.receiver_phone || '',
  from_location: props.checkout.from_location || '',
  to_location: props.checkout.to_location || '',
  vehicle_number: props.checkout.vehicle_number || '',
  driver_name: props.checkout.driver_name || '',
  driver_phone: props.checkout.driver_phone || '',
  delivery_charge: props.checkout.delivery_charge ?? null as number | null,
  delivery_cost: props.checkout.delivery_cost ?? null as number | null,
  customer_paid: props.checkout.customer_paid ?? null as number | null,
  notes: props.checkout.notes || '',
})

const openEditDialog = () => {
  editForm.type = props.checkout.type as 'pickup' | 'delivery'
  editForm.quantity = props.checkout.quantity
  editForm.receiver_name = props.checkout.receiver_name || ''
  editForm.receiver_phone = props.checkout.receiver_phone || ''
  editForm.from_location = props.checkout.from_location || ''
  editForm.to_location = props.checkout.to_location || ''
  editForm.vehicle_number = props.checkout.vehicle_number || ''
  editForm.driver_name = props.checkout.driver_name || ''
  editForm.driver_phone = props.checkout.driver_phone || ''
  editForm.delivery_charge = props.checkout.delivery_charge ?? null
  editForm.delivery_cost = props.checkout.delivery_cost ?? null
  editForm.customer_paid = props.checkout.customer_paid ?? null
  editForm.notes = props.checkout.notes || ''
  editImageFile.value = null
  editImagePreview.value = null
  showEditDialog.value = true
  open.value = false
}

const triggerEditImageInput = () => editImageInputRef.value?.click()

const handleEditImageSelect = (event: Event) => {
  const input = event.target as HTMLInputElement; const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }
  if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }
  editImageFile.value = file
  const reader = new FileReader(); reader.onload = (e) => { editImagePreview.value = e.target?.result as string }; reader.readAsDataURL(file)
}

const clearEditImage = () => { editImageFile.value = null; editImagePreview.value = null; if (editImageInputRef.value) editImageInputRef.value.value = '' }

const removeExistingImage = async () => {
  if (!props.checkout.image_url) return
  try { await deleteImage('checkouts', props.checkout.id, props.checkout.image_url); editImagePreview.value = null; push.success('Image removed'); emit('checkout-updated') }
  catch { push.error('Failed to remove image') }
}

const changeStatus = async (status: CheckoutStatus) => {
  if (status === props.checkout.status) return
  try { const result = await checkoutsStore.changeCheckoutStatus(props.checkout.id, status); if (result) { push.success(`Status changed to ${status}`); open.value = false; emit('checkout-updated') } }
  catch { push.error('Failed to update status') }
}

const handleUpdate = async () => {
  saving.value = true
  try {
    if (editImageFile.value) { await uploadImage('checkouts', props.checkout.id, editImageFile.value) }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {
      quantity: editForm.quantity, receiver_name: editForm.receiver_name || null, receiver_phone: editForm.receiver_phone || null, notes: editForm.notes || null,
    }
    payload.from_location = editForm.type === 'delivery' ? (editForm.from_location || null) : null
    payload.to_location = editForm.type === 'delivery' ? (editForm.to_location || null) : null
    payload.vehicle_number = editForm.type === 'delivery' ? (editForm.vehicle_number || null) : null
    payload.driver_name = editForm.type === 'delivery' ? (editForm.driver_name || null) : null
    payload.driver_phone = editForm.type === 'delivery' ? (editForm.driver_phone || null) : null
    payload.delivery_charge = editForm.type === 'delivery' ? editForm.delivery_charge : null
    payload.delivery_cost = editForm.type === 'delivery' ? editForm.delivery_cost : null
    payload.customer_paid = editForm.type === 'delivery' ? editForm.customer_paid : null
    const result = await checkoutsStore.updateCheckout(props.checkout.id, payload)
    if (result) { push.success('Checkout updated'); showEditDialog.value = false; emit('checkout-updated') }
  } catch { push.error('Failed to update checkout') }
  finally { saving.value = false }
}

const openDeleteConfirm = () => { open.value = false; restoreQuantity.value = true; showDeleteDialog.value = true }

const handleDelete = async () => {
  deleting.value = true
  try {
    const result = await checkoutsStore.deleteCheckout(props.checkout.id, restoreQuantity.value)
    if (result) { push.success('Checkout deleted' + (restoreQuantity.value ? ' - quantity restored' : '')); showDeleteDialog.value = false; emit('checkout-deleted') }
  } catch { push.error('Failed to delete checkout') }
  finally { deleting.value = false }
}

const getStatusClass = (s: string): string => ({
  pending: 'border-yellow-400 bg-yellow-50 text-yellow-700',
  in_transit: 'border-blue-400 bg-blue-50 text-blue-700',
  complete: 'border-green-400 bg-green-50 text-green-700'
}[s] || '')
</script>
