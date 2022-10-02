<script setup lang="ts">

import {onMounted, ref} from 'vue'
import {Terminal} from 'xterm'
import 'xterm/css/xterm.css'
import {AttachAddon} from 'xterm-addon-attach'

const terminal = ref()
const term = new Terminal({
  rows: 100,
  cols: 200,
});

onMounted(() => {
  console.log(terminal.value);// <div>

  const socketURL = "ws://127.0.0.1:8080/doterm?token=2&param=testtest1"
  const ws = new WebSocket(socketURL)

  //连接打开时触发
  ws.onopen = function (evt) {
    console.log("Connection open ...");

    // ws.send("aaa")
  };

  ws.onclose = function () {
    console.log("Connection closed!");
  }

  const attachAddon = new AttachAddon(ws)
  term.loadAddon(attachAddon);
  term.open(terminal.value)
  term.focus()

  // term.onData(send => {
  //   console.log("received: %s", send)
  //   term.write(utf16To8(send))
  // });
})

</script>

<template>
  <div ref="terminal"></div>
</template>

<style scoped>

</style>
