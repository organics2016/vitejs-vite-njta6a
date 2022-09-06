<script setup lang="ts">

import { ref, onMounted } from 'vue'
import { Terminal } from 'xterm'
import 'xterm/css/xterm.css'
import { AttachAddon } from 'xterm-addon-attach'

const terminal = ref()
const term = new Terminal();

onMounted(() => {
  console.log(terminal.value);// <div>

  const socketURL = "ws://127.0.0.1:8080/"
  const ws = new WebSocket(socketURL)
  const attachAddon = new AttachAddon(ws)

  // term.open(terminal.value);
  term.loadAddon(attachAddon);
})

</script>
  
<template>
  <div ref="terminal"></div>
</template>

<style scoped>
</style>
