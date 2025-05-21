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

    <!-- Komponenten Status -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6">
        <h2 class="text-lg font-medium leading-6 text-gray-900">Komponenten Status</h2>
      </div>
          <div class="border-t border-gray-200">
        <div class="px-4 py-5 sm:p-0">
          <dl class="sm:divide-y sm:divide-gray-200">
            <div v-for="component in filteredComponents" :key="component.id" class="py-4 sm:py-5 sm:px-6">
              <div class="flex justify-between items-start">
                <div>
                  <dt class="text-sm font-medium text-gray-900">{{ component.name }}</dt>
                  <dd class="mt-1 text-sm text-gray-500">{{ component.description }}</dd>
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
                <div>
                  <span :class="getStatusClass(component.status)" class="px-2 py-1 text-xs rounded-full">
                    {{ component.status }}
                  </span>
                </div>
              </div>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <!-- Active Incidents -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6">
        <h2 class="text-lg font-medium leading-6 text-gray-900">Aktive Vorfälle</h2>
      </div>
      <div class="border-t border-gray-200">
        <ul v-if="activeIncidents.length > 0" class="divide-y divide-gray-200">
          <li v-for="incident in activeIncidents" :key="incident.id" class="p-4">
            <div class="flex items-center space-x-4">
              <div class="flex-1">
                <h3 class="text-sm font-medium">{{ incident.title }}</h3>
                <p class="text-sm text-gray-500">{{ incident.description }}</p>
              </div>
              <div>
                <span :class="getIncidentStatusClass(incident.status)" class="px-2 py-1 text-xs rounded-full">
                  {{ incident.status }}
                </span>
              </div>
            </div>
          </li>
        </ul>
        <div v-else class="p-4 text-sm text-gray-500">
          Keine aktiven Vorfälle
        </div>
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

interface Incident {
  id: string
  title: string
  description: string
  status: string
}

const components = ref<Component[]>([])
const activeIncidents = ref<Incident[]>([])
const availableTags = ref<Tag[]>([])
const selectedTags = ref<string[]>([])

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

const getIncidentStatusClass = (status: string) => {
  switch (status) {
    case 'investigating':
      return 'bg-yellow-100 text-yellow-800'
    case 'identified':
      return 'bg-blue-100 text-blue-800'
    case 'monitoring':
      return 'bg-purple-100 text-purple-800'
    case 'resolved':
      return 'bg-green-100 text-green-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

onMounted(async () => {
  try {
    const [componentsRes, incidentsRes, tagsRes] = await Promise.all([
      axios.get('http://localhost:8080/api/components'),
      axios.get('http://localhost:8080/api/incidents'),
      axios.get('http://localhost:8080/api/tags')
    ])
    components.value = componentsRes.data
    activeIncidents.value = incidentsRes.data.filter((i: Incident) => i.status !== 'resolved')
    availableTags.value = tagsRes.data
  } catch (error) {
    console.error('Error fetching data:', error)
  }
})
</script>
