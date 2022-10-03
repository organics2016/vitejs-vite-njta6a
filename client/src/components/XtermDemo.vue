<script setup lang="ts">

import {onMounted, onUnmounted, ref} from 'vue'
import {Terminal} from 'xterm'
import 'xterm/css/xterm.css'
import {AttachAddon} from 'xterm-addon-attach'

const props = defineProps({
  host: {type: String, default: "127.0.0.1"},
  port: {type: Number, default: 2233},
  token: {type: String, required: true},
  param: String
})

const terminal = ref()
const term = new Terminal({
  // rows: 100,
  // cols: 200,
});

let ws: WebSocket

onMounted(() => {

  let socketURL = "ws://" + props.host + ":" + props.port + "/doterm?token=" + props.token
  if (props.param) {
    socketURL = socketURL + "&param=" + props.param
  }

  ws = new WebSocket(socketURL)

  ws.onopen = function () {
    console.log("Connection open ...");
  };

  ws.onclose = function () {
    console.log("Connection closed!");
  }

  const attachAddon = new AttachAddon(ws)
  term.loadAddon(attachAddon);
  term.open(terminal.value)
  term.focus()

})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})

</script>

<template>
  <div ref="terminal"></div>
</template>

<style scoped>

</style>
