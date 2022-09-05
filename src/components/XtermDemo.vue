<script setup lang="ts">

import { ref, onMounted } from 'vue'
import { Terminal } from 'xterm'
import 'xterm/css/xterm.css'


const terminal = ref()

onMounted(() => {
  console.log(terminal.value);// <div>

  const term = new Terminal();

  term.open(terminal.value);
  term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
})

</script>
  
<template>
  <div ref="terminal"></div>
</template>

<style scoped>
</style>
