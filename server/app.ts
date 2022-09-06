import { createServer } from 'http';
import { WebSocketServer } from 'ws';
import nodePty from 'node-pty';
import os from 'os';

const server = createServer();
const wss = new WebSocketServer({ server });

const shell = os.platform() === "win32" ? "powershell.exe" : "bash";
const pty = nodePty.spawn(shell, [], {
  name: 'xterm-color',
  cols: 80,
  rows: 30,
  cwd: process.env.HOME,
  env: Object.assign(
    {},
    ...Object.keys(process.env)
      .filter((key: string) => process.env[key])
      .map((key: string) => ({ [key]: process.env[key] }))
  )
});

console.log('ddd')
wss.on('connection', function connection(ws) {

  pty.onData(recv => {
    ws.send(recv);
  });

  ws.onmessage()

  ws.on('message', function message(data) {
    console.log('received: %s', data);

    pty.write(data);
  });

});


server.listen(8080);