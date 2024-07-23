'use strict';

let wasm;
const go = new Go();

function setEventListeners() {
  document.addEventListener("DOMContentLoaded", function () {
    /* let postList = fetchPostList();
    buildNav(postList);

    // Automatically load the home page by default
    fetchPost('home');
    displayPost('home'); */
  });
}

function init(wasmObj) {
  go.run(wasmObj.instance);
  setEventListeners();
}

if ('instantiateStreaming' in WebAssembly) {
  WebAssembly.instantiateStreaming(fetch("go.wasm"), go.importObject).then(wasmObj => {
    init(wasmObj);
  })
} else {
  fetch("go.wasm").then(resp =>
    resp.arrayBuffer()
  ).then(bytes =>
    WebAssembly.instantiate(bytes, go.importObject).then(wasmObj => {
      init(wasmObj);
    })
  )
}