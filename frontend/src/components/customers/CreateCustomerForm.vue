<template>
  <div class="max-w-2xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add New Customer</h2>
      <p class="text-sm text-secondary">Register a new customer in the system</p>
    </div>

    <!-- Form Card -->
    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Customer Information Section -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Customer Information
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Name -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Full Name <span class="text-warning-text">*</span>
              </label>
              <input v-model="name" type="text" placeholder="John Doe" required class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Phone -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">
                Phone <span class="text-warning-text">*</span>
              </label>
              <input v-model="phone" type="tel" placeholder="+1 234 567 890" required class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Email -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Email</label>
              <input v-model="email" type="email" placeholder="john@example.com" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>

            <!-- Company -->
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Company</label>
              <input v-model="company" type="text" placeholder="Company name" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200" />
            </div>
          </div>

          <!-- Address -->
          <div class="mt-4 space-y-1.5">
            <label class="text-sm font-medium text-primary">Address</label>
            <textarea v-model="address" rows="2" placeholder="Enter full address" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                     focus:outline-none focus:ring-2 focus:ring-accent
                     focus:border-transparent transition-all duration-200 resize-none"></textarea>
          </div>
        </div>

        <!-- ID Information Section -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Identification <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">ID Type</label>
              <select v-model="id_type" class="input w-full px-3 py-2 rounded-lg cursor-pointer
                       focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent
                       transition-all duration-200">
                <option value="" class="bg-surface text-primary">Select ID Type</option>
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
              <input v-model="id_number" type="text" placeholder="Enter ID number" :disabled="!id_type" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                       focus:outline-none focus:ring-2 focus:ring-accent
                       focus:border-transparent transition-all duration-200
                       disabled:opacity-50 disabled:cursor-not-allowed" />
            </div>
          </div>
        </div>

        <!-- Notes Section -->
        <div>
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Notes <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>

          <div class="space-y-1.5">
            <textarea v-model="notes" rows="3" placeholder="Add notes about this customer" class="input w-full px-3 py-2 rounded-lg placeholder:text-muted
                     focus:outline-none focus:ring-2 focus:ring-accent
                     focus:border-transparent transition-all duration-200 resize-none"></textarea>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">
            Clear Form
          </button>
          <button type="submit" :disabled="submitting" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg
                   shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 transform
                   disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Creating...
            </span>
            <span v-else>Create Customer</span>
          </button>
        </div>

      </form>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, watch } from 'vue'
import type { IDType } from '@/types/auth'
import { useCustomersStore } from '@/stores/customers'
import { push } from 'notivue'

const emit = defineEmits<{
  'customer-created': []
}>()

const customersStore = useCustomersStore()

const submitting = ref(false)

// Form fields
const name = ref('')
const phone = ref('')
const email = ref('')
const company = ref('')
const address = ref('')
const id_type = ref<IDType | ''>('')
const id_number = ref('')
const notes = ref('')

watch(id_type, (newVal) => {
  if (!newVal) id_number.value = ''
})

const resetForm = () => {
  name.value = ''
  phone.value = ''
  email.value = ''
  company.value = ''
  address.value = ''
  id_type.value = ''
  id_number.value = ''
  notes.value = ''
}

const submit = async () => {
  if (!name.value || !phone.value) return

  submitting.value = true

  try {
    const newCustomer = await customersStore.createCustomer({
      name: name.value,
      phone: phone.value,
      email: email.value || null,
      company: company.value || null,
      address: address.value || null,
      id_type: id_type.value || null,
      id_number: id_number.value || null,
      notes: notes.value || null,
    })

    if (newCustomer) {
      push.success('Customer created successfully!')
      resetForm()
      emit('customer-created')  // 👈 tells parent to close dialog + refresh
    }
  } catch (error) {
    console.error('Error creating customer:', error)
    push.error('Failed to create customer')
  } finally {
    submitting.value = false
  }
}
</script>
