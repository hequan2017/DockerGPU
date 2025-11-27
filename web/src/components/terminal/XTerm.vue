<template>
  <div ref="termEl" style="width:100%;height:100%;"></div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

const props = defineProps({ wsUrl: { type: String, required: true } })
const termEl = ref()
let term
let fit
let ws
let buf = ''

onMounted(() => {
  term = new Terminal({ cursorBlink: true, fontFamily: 'monospace', fontSize: 14 })
  fit = new FitAddon()
  term.loadAddon(fit)
  term.open(termEl.value)
  fit.fit()
  ws = new WebSocket(props.wsUrl)
  ws.binaryType = 'arraybuffer'
  ws.onmessage = (e) => {
    let d = ''
    if (e.data instanceof ArrayBuffer) {
      const dec = new TextDecoder()
      d = dec.decode(new Uint8Array(e.data))
    } else if (typeof e.data === 'string') {
      d = e.data
    }
    if (d) term.write(d)
  }
  ws.onerror = () => { term.write('\r\n[ws error]') }
  ws.onclose = () => { term.write('\r\n[closed]') }
  term.onData((d) => {
    try {
      const enc = new TextEncoder()
      ws.send(enc.encode(d))
    } catch {}
  })
  window.addEventListener('resize', () => fit.fit())
})

onBeforeUnmount(() => {
  try { ws && ws.close() } catch {}
  try { term && term.dispose() } catch {}
})
</script>

<style scoped>
</style>
