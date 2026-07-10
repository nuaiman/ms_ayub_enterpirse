<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add Expense</h2>
      <p class="text-sm text-secondary">Record a new expense</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Salary Toggle -->
        <div class="flex items-center gap-3 pb-4 border-b border-divider">
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="is_salary" class="sr-only peer">
            <div
              class="w-11 h-6 bg-surface-alt peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-accent after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all">
            </div>
          </label>
          <span class="text-sm font-medium text-primary">This is a salary payment</span>
        </div>

        <!-- Salary Fields -->
        <div v-if="is_salary" class="space-y-4 pb-4 border-b border-divider">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Employee <span class="text-warning-text">*</span></label>
              <select v-model="salary_user_id"
                class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option :value="null" disabled>Select employee</option>
                <option v-for="user in staffUsers" :key="user.id" :value="user.id">
                  {{ user.name }} ({{ user.role }}) - ৳{{ user.monthly_salary || 0 }}
                </option>
              </select>
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Month/Year <span
                  class="text-warning-text">*</span></label>
              <input v-model="salary_month_year" type="month"
                class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
          <p v-if="selectedUser" class="text-xs text-muted">
            Monthly salary: ৳{{ selectedUser.monthly_salary || 0 }}
          </p>
        </div>

        <!-- Expense Details -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Category with autocomplete -->
          <div class="space-y-1.5 relative">
            <label class="text-sm font-medium text-primary">Category</label>
            <div class="relative">
              <input v-model="category" type="text" placeholder="e.g. office supplies, utilities..."
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
                <button v-for="cat in filteredCategorySuggestions" :key="cat" @click="selectCategory(cat)" type="button"
                  class="w-full flex items-center px-4 py-2.5 text-sm text-secondary hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                  {{ cat }}
                </button>
              </div>
            </div>
          </div>

          <!-- Amount -->
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Amount <span class="text-warning-text">*</span></label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
              <input v-model.number="amount" type="number" step="0.01" min="0" placeholder="0.00" required
                class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>

          <!-- Date -->
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Date <span class="text-warning-text">*</span></label>
            <input v-model="expense_date" type="date" required
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <!-- Notes -->
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Notes</label>
            <textarea v-model="notes" rows="2" placeholder="Add notes..."
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
          </div>
        </div>

        <!-- Image Upload -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Receipt / Image <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="flex items-start gap-4">
            <div class="shrink-0">
              <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default">
                <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                <button type="button" @click="removeImage"
                  class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">×</button>
              </div>
              <div v-else
                class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt">
                <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
            </div>
            <div class="flex-1 space-y-3">
              <label
                class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2"
                :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                {{ uploading ? 'Uploading...' : 'Choose File' }}
                <input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" />
              </label>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit"
            :disabled="submitting || !amount || amount <= 0 || (is_salary && (!salary_user_id || !salary_month_year))"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Creating...
            </span>
            <span v-else>Add Expense</span>
          </button>
        </div>
      </form>
    </div>

    <div v-if="showCategoryDropdown" class="fixed inset-0 z-40" @click="showCategoryDropdown = false"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useExpensesStore } from '@/stores/expenses'
import { useUsersStore } from '@/stores/users'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{ 'expense-created': [] }>()

const expensesStore = useExpensesStore()
const usersStore = useUsersStore()

const submitting = ref(false)
const is_salary = ref(false)
const salary_user_id = ref<number | null>(null)
const salary_month_year = ref('')
const category = ref('')
const amount = ref<number | null>(null)
const expense_date = ref(new Date().toISOString().split('T')[0] as string)
const notes = ref('')

const showCategoryDropdown = ref(false)

const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)

const staffUsers = computed(() => {
  return usersStore.users.filter(u => u.role !== 'admin' && u.is_active)
})

const selectedUser = computed(() => {
  return staffUsers.value.find(u => u.id === salary_user_id.value)
})

const categorySuggestions = computed(() => {
  const cats = new Set<string>()
  expensesStore.expenses.forEach(e => { if (e.category) cats.add(e.category.toLowerCase()) })
  return Array.from(cats).sort()
})

const filteredCategorySuggestions = computed(() => {
  if (!category.value) return categorySuggestions.value
  const q = category.value.toLowerCase()
  return categorySuggestions.value.filter(c => c.toLowerCase().includes(q) && c !== q)
})

const selectCategory = (cat: string) => { category.value = cat; showCategoryDropdown.value = false }

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
  is_salary.value = false; salary_user_id.value = null; salary_month_year.value = ''
  category.value = ''; amount.value = null
  expense_date.value = new Date().toISOString().split('T')[0] as string
  notes.value = ''; removeImage()
}

const submit = async () => {
  if (!amount.value || amount.value <= 0 || !expense_date.value) return
  if (is_salary.value && (!salary_user_id.value || !salary_month_year.value)) {
    push.error('Employee and month are required for salary')
    return
  }

  submitting.value = true
  try {
    // Auto-generate title from date
    const dateObj = new Date(expense_date.value)
    const formattedDate = dateObj.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
    const autoTitle = is_salary.value
      ? `Salary - ${selectedUser.value?.name || 'Unknown'} (${salary_month_year.value})`
      : `Expense - ${formattedDate}`

    const newExpense = await expensesStore.createExpense({
      is_salary: is_salary.value,
      salary_user_id: salary_user_id.value,
      salary_month_year: salary_month_year.value || null,
      title: autoTitle,
      category: category.value.trim().toLowerCase() || null,
      amount: Number(amount.value),
      expense_date: expense_date.value,
      notes: notes.value || null,
    })

    if (!newExpense) { submitting.value = false; return }

    if (imageFile.value) {
      uploading.value = true
      await uploadImage('expenses', newExpense.id, imageFile.value)
      uploading.value = false
    }

    push.success('Expense added successfully!')
    resetForm()
    emit('expense-created')
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
  } catch (error) { push.error('Failed to create expense') }
  finally { submitting.value = false }
}
</script>
