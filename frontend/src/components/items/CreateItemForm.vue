<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add New Item</h2>
      <p class="text-sm text-secondary">Register a new item in inventory</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Customer Information -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Customer Information <span class="text-warning-text">*</span>
          </h3>
          <div class="space-y-4">
            <!-- Name with autocomplete dropdown -->
            <div class="space-y-1.5 relative">
              <label class="text-sm font-medium text-primary">Name <span class="text-warning-text">*</span></label>
              <div class="relative">
                <input v-model="customer_name" type="text" placeholder="Full name / company / organization" required
                  @focus="showNameDropdown = true" @input="showNameDropdown = true"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button v-if="customer_name" @click="customer_name = ''" type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <!-- Name Suggestions Dropdown -->
              <div v-if="showNameDropdown && filteredNameSuggestions.length > 0"
                class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
                <div class="max-h-48 overflow-y-auto">
                  <button v-for="name in filteredNameSuggestions" :key="name" @click="selectName(name)" type="button"
                    class="w-full flex items-center px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                    {{ name }}
                  </button>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Phone <span class="text-warning-text">*</span></label>
                <input v-model="customer_phone" type="tel" placeholder="+880 1XXX-XXXXXX" required
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5">
                <label class="text-sm font-medium text-primary">Email</label>
                <input v-model="customer_email" type="email" placeholder="customer@example.com"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
          </div>
        </div>

        <!-- Item Details -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Item Details</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <!-- Category with autocomplete dropdown -->
            <div class="space-y-1.5 relative">
              <label class="text-sm font-medium text-primary">Category</label>
              <div class="relative">
                <input v-model="category" type="text" placeholder="e.g. electronics, furniture..."
                  @focus="showCategoryDropdown = true" @input="showCategoryDropdown = true"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button v-if="category" @click="category = ''" type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div v-if="showCategoryDropdown && filteredCategorySuggestions.length > 0"
                class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
                <div class="max-h-48 overflow-y-auto">
                  <button v-for="cat in filteredCategorySuggestions" :key="cat" @click="selectCategory(cat)"
                    type="button"
                    class="w-full flex items-center px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                    {{ cat }}
                  </button>
                </div>
              </div>
            </div>

            <!-- Subcategory with autocomplete dropdown -->
            <div class="space-y-1.5 relative">
              <label class="text-sm font-medium text-primary">Subcategory</label>
              <div class="relative">
                <input v-model="subcategory" type="text" placeholder="e.g. laptop, sofa..."
                  @focus="showSubcategoryDropdown = true" @input="showSubcategoryDropdown = true"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button v-if="subcategory" @click="subcategory = ''" type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div v-if="showSubcategoryDropdown && filteredSubcategorySuggestions.length > 0"
                class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">
                <div class="max-h-48 overflow-y-auto">
                  <button v-for="sub in filteredSubcategorySuggestions" :key="sub" @click="selectSubcategory(sub)"
                    type="button"
                    class="w-full flex items-center px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                    {{ sub }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                <option value="bag">Bag</option>
                <option value="bottle">Bottle</option>
                <option value="box">Box</option>
                <option value="can">Can</option>
                <option value="carton">Carton</option>
                <option value="cup">Cup</option>
                <option value="dozen">Dozen</option>
                <option value="gallon">Gallon</option>
                <option value="pack">Pack</option>
                <option value="pair">Pair</option>
                <option value="pcs">Pieces</option>
                <option value="roll">Roll</option>
                <option value="set">Set</option>
                <option value="sheet">Sheet</option>
                <option value="unit">Unit</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Dimensions & Weight -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Dimensions & Weight <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Weight</label><input
                v-model.number="weight" type="number" step="0.001" min="0" placeholder="0.00"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Weight Unit</label><select
                v-model="weight_unit"
                class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option value="kg">kg</option>
                <option value="g">g</option>
                <option value="mg">mg</option>
                <option value="oz">oz</option>
                <option value="lb">lb</option>
                <option value="ton">ton</option>
              </select></div>
          </div>
          <div class="space-y-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Dimension Unit</label><select
                v-model="dimension_unit"
                class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent max-w-xs">
                <option value="in">inches</option>
                <option value="cm">cm</option>
                <option value="mm">mm</option>
                <option value="ft">ft</option>
                <option value="m">m</option>
                <option value="yd">yd</option>
                <option value="km">km</option>
              </select></div>
            <div class="grid grid-cols-3 gap-4">
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Length</label><input
                  v-model.number="length" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Width</label><input
                  v-model.number="width" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Height</label><input
                  v-model.number="height" type="number" step="0.01" min="0" placeholder="0.00"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
          </div>
        </div>

        <!-- Storage Contract -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Storage Contract <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Duration Type</label><select
                  v-model="duration_type"
                  class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                  <option value="">Select Type</option>
                  <option value="day">Day</option>
                  <option value="week">Week</option>
                  <option value="month">Month</option>
                  <option value="year">Year</option>
                </select></div>
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Duration</label><input
                  v-model.number="duration" type="number" min="1" placeholder="Duration"
                  class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Start Date</label><input
                  v-model="start_date" type="date"
                  class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
              </div>
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Amount <span
                    class="text-warning-text" v-if="duration_type">*</span></label>
                <div class="relative"><span
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input
                    v-model.number="amount" type="number" step="0.01" min="0" placeholder="0.00"
                    class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                </div>
              </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Security Deposit</label>
                <div class="relative"><span
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input
                    v-model.number="deposit" type="number" step="0.01" min="0" placeholder="0.00"
                    class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                </div>
              </div>
              <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Customer Paid</label>
                <div class="relative"><span
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input
                    v-model.number="customer_paid" type="number" step="0.01" min="0" placeholder="0.00"
                    class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Image & Notes -->
        <div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Image <span
                  class="normal-case text-xs font-normal">(Optional)</span></h3>
              <div class="flex items-start gap-4">
                <div class="shrink-0">
                  <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default">
                    <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" /><button type="button"
                      @click="removeImage"
                      class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">×</button>
                  </div>
                  <div v-else
                    class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt">
                    <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg></div>
                </div>
                <div class="flex-1 space-y-3">
                  <label
                    class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2"
                    :class="{ 'opacity-50 cursor-not-allowed': uploading }"><svg class="w-4 h-4" fill="none"
                      stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                    </svg>{{ uploading ? 'Uploading...' : 'Choose File' }}<input type="file" accept="image/*"
                      class="hidden" @change="handleFileSelect" :disabled="uploading" /></label>
                </div>
              </div>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Notes <span
                  class="normal-case text-xs font-normal">(Optional)</span></h3>
              <textarea v-model="notes" rows="4" placeholder="Add notes..."
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit" :disabled="isSubmitDisabled"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2"><svg class="animate-spin w-4 h-4" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>Creating...</span>
            <span v-else>Create Item</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Click outside to close dropdowns -->
    <div v-if="showNameDropdown || showCategoryDropdown || showSubcategoryDropdown" class="fixed inset-0 z-40"
      @click="closeAllDropdowns"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { CreateItemPayload } from '@/types/item'
import { useItemsStore } from '@/stores/items'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{ 'item-created': [] }>()
const itemsStore = useItemsStore()

const submitting = ref(false)

const customer_name = ref('')
const customer_phone = ref('')
const customer_email = ref('')

const category = ref('')
const subcategory = ref('')

const quantity = ref(1)
const quantity_unit = ref('pcs')

const weight = ref<number | null>(null)
const weight_unit = ref('kg')
const length = ref<number | null>(null)
const width = ref<number | null>(null)
const height = ref<number | null>(null)
const dimension_unit = ref('in')

const duration_type = ref('')
const duration = ref<number | null>(null)
const start_date = ref('')
const amount = ref<number | null>(null)
const deposit = ref<number | null>(null)
const customer_paid = ref<number | null>(null)

const notes = ref('')

const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)

// Dropdown visibility
const showNameDropdown = ref(false)
const showCategoryDropdown = ref(false)
const showSubcategoryDropdown = ref(false)

// Suggestions - unique, lowercase, sorted
const nameSuggestions = computed(() => {
  const names = new Set<string>()
  itemsStore.items.forEach(i => { if (i.name) names.add(i.name.toLowerCase()) })
  return Array.from(names).sort()
})

const categorySuggestions = computed(() => {
  const cats = new Set<string>()
  itemsStore.items.forEach(i => { if (i.category) cats.add(i.category.toLowerCase()) })
  return Array.from(cats).sort()
})

const subcategorySuggestions = computed(() => {
  const subs = new Set<string>()
  itemsStore.items.forEach(i => { if (i.subcategory) subs.add(i.subcategory.toLowerCase()) })
  return Array.from(subs).sort()
})

// Filtered suggestions based on input
const filteredNameSuggestions = computed(() => {
  if (!customer_name.value) return nameSuggestions.value
  const q = customer_name.value.toLowerCase()
  return nameSuggestions.value.filter(n => n.toLowerCase().includes(q) && n !== q)
})

const filteredCategorySuggestions = computed(() => {
  if (!category.value) return categorySuggestions.value
  const q = category.value.toLowerCase()
  return categorySuggestions.value.filter(c => c.toLowerCase().includes(q) && c !== q)
})

const filteredSubcategorySuggestions = computed(() => {
  if (!subcategory.value) return subcategorySuggestions.value
  const q = subcategory.value.toLowerCase()
  return subcategorySuggestions.value.filter(s => s.toLowerCase().includes(q) && s !== q)
})

const isSubmitDisabled = computed(() => submitting.value || !customer_name.value || !customer_phone.value || (!!duration_type.value && !amount.value))

const selectName = (name: string) => { customer_name.value = name; showNameDropdown.value = false }
const selectCategory = (cat: string) => { category.value = cat; showCategoryDropdown.value = false }
const selectSubcategory = (sub: string) => { subcategory.value = sub; showSubcategoryDropdown.value = false }
const closeAllDropdowns = () => { showNameDropdown.value = false; showCategoryDropdown.value = false; showSubcategoryDropdown.value = false }

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
  customer_name.value = ''; customer_phone.value = ''; customer_email.value = ''
  category.value = ''; subcategory.value = ''
  quantity.value = 1; quantity_unit.value = 'pcs'
  weight.value = null; weight_unit.value = 'kg'
  length.value = null; width.value = null; height.value = null; dimension_unit.value = 'in'
  duration_type.value = ''; duration.value = null; start_date.value = ''
  amount.value = null; deposit.value = null; customer_paid.value = null
  notes.value = ''; removeImage()
  closeAllDropdowns()
}

const submit = async () => {
  if (!customer_name.value || !customer_phone.value) return
  if (duration_type.value && !amount.value) { push.error('Amount is required for storage contracts'); return }
  submitting.value = true
  try {
    const payload: CreateItemPayload = {
      name: customer_name.value.trim(), customer_phone: customer_phone.value.trim(), customer_email: customer_email.value?.trim() || null,
      category: category.value.trim().toLowerCase() || null, subcategory: subcategory.value.trim().toLowerCase() || null,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      quantity_unit: quantity_unit.value as any, quantity: quantity.value,
      weight: weight.value || null,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      weight_unit: weight_unit.value as any,
      length: length.value || null, width: width.value || null, height: height.value || null,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      dimension_unit: dimension_unit.value as any, notes: notes.value || null,
    }
    if (duration_type.value) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      payload.duration_type = duration_type.value as any; payload.duration = duration.value || null
      payload.start_date = start_date.value || null; payload.amount = amount.value || 0
      payload.deposit = deposit.value || 0; payload.customer_paid = customer_paid.value || 0
    }
    const newItem = await itemsStore.createItem(payload)
    if (!newItem) { submitting.value = false; return }
    if (imageFile.value) { uploading.value = true; await uploadImage('items', newItem.id, imageFile.value); uploading.value = false }
    push.success('Item created successfully!'); resetForm(); emit('item-created')
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (error) { push.error('Failed to create item') }
  finally { submitting.value = false }
}
</script>
