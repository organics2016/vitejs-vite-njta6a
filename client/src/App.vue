<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import {reactive, ref} from 'vue'
import XtermDemoVue from './components/XtermDemo.vue'


const from = reactive({
  host: undefined,
  port: undefined,
  token: undefined,
  param: undefined,
})

const term = reactive({
  host: undefined,
  port: undefined,
  token: undefined,
  param: undefined,
})

const termKey = ref(new Date())

function refresh() {
  Object.assign(term, from)
  termKey.value = new Date()
}

</script>

<template>
  <div class="from">
    <div class="from-row">
      <span>doterm host:</span><input type="text" v-model="from.host">
    </div>
    <div class="from-row">
      <span>doterm port:</span><input type="number" v-model="from.port">
    </div>
    <div class="from-row">
      <span>token:</span><input type="text" v-model="from.token">
    </div>
    <div class="from-row">
      <span>param:</span><input type="text" v-model="from.param">
    </div>
    <button @click="refresh">Refresh</button>
  </div>
  <XtermDemoVue :host="term.host" :port="term.port" :token="term.token" :param="term.param" v-if="term.token"
                :key="termKey"/>
</template>

<style scoped>
.from {
  display: flex;
  flex-direction: column;
  padding: 1rem;
}

.from-row {
  display: flex;
  flex-direction: column;
  padding: 0.5rem 0;
}
</style>
