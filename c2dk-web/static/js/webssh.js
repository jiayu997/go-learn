Terminal.applyAddon(fit);

const addr = prompt("ssh", "127.0.0.1:22")

const term = new Terminal({
    useStyle:false,
    convertEol: true,
    cursorBlink:false,
    rendererType: "canvas",
    theme:{
        foreground: 'yellow',
        background: 'rgba(6,1,1,0.55)',
    }
});
term.open(document.getElementById("main"));

const ws = new WebSocket(`ws://${location.host}/webssh/data?addr=${addr}`);
ws.onopen = () => {
  ws.send(JSON.stringify({ type: "login", data: utoa("login-test") }));
  term.on("data", data => {
    const msg = { type: "stdin", data: btoa(data) };
    ws.send(JSON.stringify(msg));
  });
  term.on("resize", e => {
    const msg = { type: "resize", ...e };
    ws.send(JSON.stringify(msg));
  });
  term.fit();
  window.addEventListener("resize", () => term.fit());
};
ws.onmessage = e => {
  const msg = JSON.parse(e.data);
  switch (msg.type) {
    case "stdout":
    case "stderr":
      term.write(atou(msg.data));
  }
};
ws.onerror = console.error;

function atou(encodeString) {
  return decodeURIComponent(escape(atob(encodeString)));
}
function utoa(rawString) {
  return btoa(encodeURIComponent(rawString));
}