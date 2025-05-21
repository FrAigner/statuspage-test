<template>
  <div class="space-y-8">
    <!-- Filter-Bereich -->
    <div class="bg-white shadow rounded-lg p-4">
      <div class="flex flex-wrap gap-2">
        <div v-for="tag in availableTags" :key="tag.id" class="flex items-center">
          <label class="inline-flex items-center">
            <input 
              type="checkbox" 
              :value="tag.name" 
              v-model="selectedTags"
              class="form-checkbox text-blue-600"
            >
            <span class="ml-2">{{ tag.name }}</span>
          </label>
        </div>
      </div>
    </div>

    <!-- Komponenten-Liste -->
    <div class="bg-white shadow rounded-lg">
      <div class="p-4 border-b border-gray-200 flex justify-between items-center">
        <h2 class="text-lg font-medium text-gray-900">Komponenten</h2>
        <button
          @click="openAddModal"
          class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
        >
          Komponente hinzufügen
        </button>
      </div>

      <div class="divide-y divide-gray-200">
        <div v-for="component in filteredComponents" :key="component.id" class="p-4">
          <div class="flex justify-between items-start">
            <div>
              <h3 class="text-lg font-medium text-gray-900">{{ component.name }}</h3>
              <p class="text-sm text-gray-500">{{ component.description }}</p>
              <div class="mt-2 flex flex-wrap gap-2">
                <span
                  v-for="tag in component.tags"
                  :key="tag.id"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  {{ tag.name }}
                </span>
              </div>
            </div>
            <div class="flex items-center space-x-4">
              <span :class="getStatusClass(component.status)" class="px-2 py-1 text-xs rounded-full">
                {{ component.status }}
              </span>
              <button
                @click="editComponent(component)"
                class="text-blue-600 hover:text-blue-800"
              >
                Bearbeiten
              </button>
              <button
                @click="deleteComponent(component.id)"
                class="text-red-600 hover:text-red-800"
              >
                Löschen
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal für Komponente hinzufügen/bearbeiten -->
    <div v-if="showAddModal" class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4">
      <div class="bg-white rounded-lg max-w-lg w-full p-6">
        <h3 class="text-lg font-medium mb-4">
          {{ editingComponent ? 'Komponente bearbeiten' : 'Neue Komponente' }}
        </h3>
        <form @submit.prevent="saveComponent">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Name</label>
              <input
                v-model="componentForm.name"
                type="text"
                required
                :class="[
                  'mt-1 block w-full rounded-md shadow-sm focus:ring-blue-500',
                  (validateForm?.value?.errors?.name)
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                ]"
              >
              <p v-if="validateForm?.value?.errors?.name" class="mt-1 text-sm text-red-600">
                {{ validateForm?.value?.errors?.name }}
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700">Beschreibung</label>
              <textarea
                v-model="componentForm.description"
                rows="3"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700">Status</label>
              <select
                v-model="componentForm.status"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              >
                <option value="operational">Operational</option>
                <option value="degraded">Degraded</option>
                <option value="outage">Outage</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700">Tags</label>
              <div class="mt-1">
                <input
                  v-model="newTag"
                  @keydown.enter.prevent="addTag"
                  type="text"
                  placeholder="Tag eingeben und Enter drücken"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                >
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="tag in componentForm.tags"
                    :key="tag"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                  >
                    {{ tag }}
                    <button
                      @click="removeTag(tag)"
                      class="ml-1 text-blue-400 hover:text-blue-600"
                    >
                      ×
                    </button>
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="showAddModal = false"
              class="px-4 py-2 border rounded-md text-gray-700 hover:bg-gray-50"
            >
              Abbrechen
            </button>
            <button
              type="submit"
              :disabled="!validateForm?.value?.isValid || isSubmitting"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmitting 
                  ? 'Wird gespeichert...' 
                  : editingComponent 
                    ? 'Aktualisieren' 
                    : 'Erstellen' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

interface Tag {
  id: string
  name: string
}

interface Component {
  id: string
  name: string
  description: string
  status: string
  tags: Tag[]
}

const components = ref<Component[]>([])
const availableTags = ref<Tag[]>([])
const selectedTags = ref<string[]>([])
const showAddModal = ref(false)
const editingComponent = ref<Component | null>(null)
const newTag = ref('')
const isSubmitting = ref(false)

const componentForm = ref({
  name: '',
  description: '',
  status: 'operational',
  tags: [] as string[]
})

const validateForm = computed(() => {
  return {
    isValid: !!componentForm.value.name.trim(),
    errors: {
      name: !componentForm.value.name.trim() ? 'Name ist erforderlich' : null
    }
  }
})

const filteredComponents = computed(() => {
  if (selectedTags.value.length === 0) return components.value
  
  return components.value.filter(component =>
    component.tags.some(tag => selectedTags.value.includes(tag.name))
  )
})

const getStatusClass = (status: string) => {
  switch (status) {
    case 'operational':
      return 'bg-green-100 text-green-800'
    case 'degraded':
      return 'bg-yellow-100 text-yellow-800'
    case 'outage':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const fetchComponents = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/components')
    components.value = response.data
  } catch (error) {
    console.error('Error fetching components:', error)
  }
}

const fetchTags = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/tags')
    availableTags.value = response.data
  } catch (error) {
    console.error('Error fetching tags:', error)
  }
}

const addTag = () => {
  if (newTag.value && !componentForm.value.tags.includes(newTag.value)) {
    componentForm.value.tags.push(newTag.value)
  }
  newTag.value = ''
}

const removeTag = (tag: string) => {
  componentForm.value.tags = componentForm.value.tags.filter(t => t !== tag)
}

const editComponent = (component: Component) => {
  editingComponent.value = component
  componentForm.value.name = component.name
  componentForm.value.description = component.description
  componentForm.value.status = component.status
  componentForm.value.tags = component.tags.map(t => t.name)
  showAddModal.value = true
}

const openAddModal = () => {
  console.log('openAddModal aufgerufen')
  componentForm.value.name = ''
  componentForm.value.description = ''
  componentForm.value.status = 'operational'
  componentForm.value.tags = []
  editingComponent.value = null
  showAddModal.value = true
}

const resetForm = () => {
  componentForm.value.name = ''
  componentForm.value.description = ''
  componentForm.value.status = 'operational'
  componentForm.value.tags = []
  editingComponent.value = null
  newTag.value = ''
  isSubmitting.value = false
}

const saveComponent = async () => {
  if (!validateForm.value.isValid || isSubmitting.value) return;

  try {
    isSubmitting.value = true;

    // Tags von string[] zu {name: string}[] mappen
    const payloadTags = componentForm.value.tags.map(tagName => ({ name: tagName }));

    const payload = {
      name: componentForm.value.name,
      description: componentForm.value.description,
      status: componentForm.value.status,
      tags: payloadTags, // Gemappte Tags verwenden
    };

    if (editingComponent.value) {
      await axios.put(`http://localhost:8080/api/components/${editingComponent.value.id}`, payload);
    } else {
      await axios.post('http://localhost:8080/api/components', payload);
    }
    showAddModal.value = false;
    await fetchComponents(); // Daten neu laden
    resetForm();
  } catch (error) {
    console.error('Error saving component:', error);
    // Hier könnten benutzerfreundliche Fehlermeldungen hinzugefügt werden
  } finally {
    isSubmitting.value = false;
  }
};

const deleteComponent = async (id: string) => {
  if (!confirm('Möchten Sie diese Komponente wirklich löschen?')) return
  
  try {
    await axios.delete(`http://localhost:8080/api/components/${id}`)
    await fetchComponents()
  } catch (error) {
    console.error('Error deleting component:', error)
  }
}

onMounted(async () => {
  fetchComponents()
  fetchTags()
})
</script>
