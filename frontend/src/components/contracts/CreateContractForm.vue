<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">New Contract</h2>
      <p class="text-sm text-secondary">Create a new contract for a customer</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Contract Details
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Customer Search / Select -->
            <div class="space-y-1.5 md:col-span-2">
              <label class="text-sm font-medium text-primary">
                Customer <span class="text-warning-text">*</span>
              </label>

              <!-- Selected customer display -->
              <div v-if="selectedCustomer"
                class="flex items-center gap-3 p-3 rounded-lg border border-success-border bg-success-bg">
                <div
                  class="w-10 h-10 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
                  <span class="text-sm font-semibold text-muted">{{ getInitials(selectedCustomer.name) }}</span>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-primary truncate">{{ selectedCustomer.name }}</p>
                  <p class="text-xs text-secondary truncate">{{ selectedCustomer.phone }}</p>
                  <p v-if="selectedCustomer.company" class="text-xs text-muted truncate">{{ selectedCustomer.company }}
                  </p>
                </div>
                <button @click="clearSelectedCustomer" type="button"
                  class="shrink-0 text-muted hover:text-warning-text transition-colors p-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <!-- Search / Select (when no customer selected) -->
              <div v-else class="relative">
                <div class="relative">
                  <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted pointer-events-none"
                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  <input ref="searchInputRef" v-model="customerSearch" type="text"
                    placeholder="Search customers by name, phone, or company..." @focus="showDropdown = true"
                    @keydown.escape="showDropdown = false" class="input w-full pl-9 pr-8 py-2.5 rounded-lg placeholder:text-muted
                           focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                  <button v-if="customerSearch" @click="customerSearch = ''" type="button"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>

                <!-- Dropdown -->
                <div v-if="showDropdown"
                  class="absolute z-50 left-0 right-0 mt-1 border border-default rounded-xl shadow-lg overflow-hidden bg-surface">

                  <!-- Quick add button at top -->
                  <div class="px-3 py-2 border-b border-divider">
                    <button @click="openQuickAdd" type="button"
                      class="w-full flex items-center gap-2 px-2 py-2 rounded-lg text-sm font-medium text-accent hover:bg-accent/5 transition-colors">
                      <div class="w-8 h-8 rounded-full bg-accent/10 flex items-center justify-center shrink-0">
                        <svg class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
                        </svg>
                      </div>
                      <span>Create New Customer</span>
                    </button>
                  </div>

                  <!-- Customer list -->
                  <div class="max-h-56 overflow-y-auto">
                    <div v-if="filteredCustomers.length === 0" class="px-4 py-6 text-center">
                      <div class="w-10 h-10 mx-auto rounded-full bg-surface-alt flex items-center justify-center mb-2">
                        <svg class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                      </div>
                      <p class="text-sm text-muted">No customers match "{{ customerSearch }}"</p>
                    </div>
                    <button v-for="c in filteredCustomers" :key="c.id" @click="selectCustomer(c)" type="button"
                      class="w-full flex items-center gap-3 px-4 py-3 hover:bg-surface-alt transition-colors text-left border-b border-divider last:border-b-0">
                      <div
                        class="w-9 h-9 rounded-full bg-surface-alt border border-default flex items-center justify-center shrink-0">
                        <span class="text-xs font-semibold text-muted">{{ getInitials(c.name) }}</span>
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium text-primary truncate">{{ c.name }}</p>
                        <div class="flex items-center gap-2 mt-0.5">
                          <span class="text-xs text-muted truncate">{{ c.phone }}</span>
                          <span v-if="c.company"
                            class="text-[10px] text-muted bg-surface-alt px-1.5 py-0.5 rounded truncate max-w-30">{{
                              c.company }}</span>
                        </div>
                      </div>
                      <svg class="w-4 h-4 text-muted shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                      </svg>
                    </button>
                  </div>

                  <!-- Footer with count -->
                  <div class="px-4 py-2 border-t border-divider bg-surface-alt">
                    <p class="text-[11px] text-muted">
                      {{ filteredCustomers.length }} customer{{ filteredCustomers.length !== 1 ? 's' : '' }} found
                    </p>
                  </div>
                </div>

                <!-- Backdrop to close dropdown -->
                <div v-if="showDropdown" class="fixed inset-0 z-40" @click="showDropdown = false"></div>
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Duration Type <span class="text-warning-text">*</span>
              </label>
              <select v-model="duration_type" required class="input w-full px-3 py-2 rounded-lg cursor-pointer
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option value="day" class="bg-surface text-primary">Day</option>
                <option value="week" class="bg-surface text-primary">Week</option>
                <option value="month" class="bg-surface text-primary">Month</option>
                <option value="year" class="bg-surface text-primary">Year</option>
              </select>
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Duration</label>
              <input v-model.number="duration" type="number" min="1" placeholder="e.g. 3" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Price</label>
              <input v-model.number="price" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Security Deposit</label>
              <input v-model.number="security_deposit" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>

            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Estimated Value</label>
              <input v-model.number="estimated_value" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
        </div>

        <div>
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Notes <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <textarea v-model="notes" rows="3" placeholder="Add contract notes..." class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                   focus:outline-none focus:ring-2 focus:ring-accent
                   focus:border-transparent resize-none"></textarea>
        </div>

        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">
            Clear
          </button>
          <button type="submit" :disabled="submitting || !selectedCustomer" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg
                   shadow-sm hover:shadow-md transition-all duration-200 active:scale-95
                   disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Creating...
            </span>
            <span v-else>Create Contract</span>
          </button>
        </div>
      </form>
    </div>

    <!-- Quick Add Customer Dialog -->
    <BaseDialog v-model="showQuickAdd" size="lg">
      <CreateCustomerForm @customer-created="handleCustomerAdded" />
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { DurationType } from '@/types/contract'
import type { Customer } from '@/types/customer'
import { useContractsStore } from '@/stores/contracts'
import { useCustomersStore } from '@/stores/customers'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import CreateCustomerForm from '@/components/customers/CreateCustomerForm.vue'
import { push } from 'notivue'

const emit = defineEmits<{ 'contract-created': [] }>()

const contractsStore = useContractsStore()
const customersStore = useCustomersStore()

const submitting = ref(false)
const duration_type = ref<DurationType>('month')
const duration = ref<number | null>(null)
const price = ref<number | null>(null)
const security_deposit = ref<number | null>(null)
const estimated_value = ref<number | null>(null)
const notes = ref('')

// Customer search/select
const searchInputRef = ref<HTMLInputElement | null>(null)
const customerSearch = ref('')
const showDropdown = ref(false)
const selectedCustomer = ref<Customer | null>(null)
const showQuickAdd = ref(false)

const customers = ref<Customer[]>([])

onMounted(async () => {
  await customersStore.fetchCustomers()
  customers.value = customersStore.customers
})

const filteredCustomers = computed(() => {
  if (!customerSearch.value) return customers.value
  const q = customerSearch.value.toLowerCase()
  return customers.value.filter(c =>
    c.name.toLowerCase().includes(q) ||
    c.phone.includes(q) ||
    (c.company && c.company.toLowerCase().includes(q))
  )
})

const selectCustomer = (customer: Customer) => {
  selectedCustomer.value = customer
  customerSearch.value = ''
  showDropdown.value = false
}

const clearSelectedCustomer = () => {
  selectedCustomer.value = null
  customerSearch.value = ''
  showDropdown.value = false
}

const openQuickAdd = () => {
  showDropdown.value = false
  customerSearch.value = ''
  showQuickAdd.value = true
}

const handleCustomerAdded = async () => {
  showQuickAdd.value = false
  await customersStore.fetchCustomers()
  customers.value = customersStore.customers

  if (customers.value.length > 0) {
    const lastCustomer = customers.value[0]
    if (lastCustomer) {
      selectedCustomer.value = lastCustomer
      push.success('Customer added and selected')
    }
  }
}

const resetForm = () => {
  clearSelectedCustomer()
  duration_type.value = 'month'
  duration.value = null
  price.value = null
  security_deposit.value = null
  estimated_value.value = null
  notes.value = ''
}

const submit = async () => {
  if (!selectedCustomer.value) return

  submitting.value = true

  try {
    const newContract = await contractsStore.createContract({
      customer_id: selectedCustomer.value.id,
      duration_type: duration_type.value,
      duration: duration.value,
      price: price.value,
      security_deposit: security_deposit.value,
      estimated_value: estimated_value.value,
      notes: notes.value || null,
    })

    if (newContract) {
      push.success('Contract created successfully!')
      resetForm()
      emit('contract-created')
    }
  } catch (error) {
    console.error('Error creating contract:', error)
    push.error('Failed to create contract')
  } finally {
    submitting.value = false
  }
}

const getInitials = (name: string): string =>
  name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
</script>
