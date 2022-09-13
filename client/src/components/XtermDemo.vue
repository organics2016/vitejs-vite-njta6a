<script setup lang="ts">

import {onMounted, ref} from 'vue'
import {Terminal} from 'xterm'
import 'xterm/css/xterm.css'
import {AttachAddon} from 'xterm-addon-attach'

const terminal = ref()
const term = new Terminal();

function utf16To8(input: string) {
  const _unescape = function(s: string) {
    function d(x:any, n:string) {
      return String.fromCharCode(parseInt(n, 16));
    }
    return s.replace(/%([0-9A-F]{2})/ig, d);
  };
  try{
    return _unescape(encodeURIComponent(input));
  }catch (URIError) {
    //include invalid character, cannot convert
    return input;
  }
}

onMounted(() => {
  console.log(terminal.value);// <div>

  const socketURL = "ws://127.0.0.1:8080/ping"
  const ws = new WebSocket(socketURL)
  const attachAddon = new AttachAddon(ws)
  term.loadAddon(attachAddon);

  //连接打开时触发
  ws.onopen = function (evt) {
    console.log("Connection open ...");
  };
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
