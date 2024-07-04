'use strict';

let wasm;
const go = new Go();

//function setEventListeners() {}

function renderPostContent(content) {
  // Parse and render markdown content using a library like marked.js
  const renderedContent = marked.parse(content);
  document.getElementById('post-content').innerHTML = renderedContent;
}

async function fetchPostContent(postId) {
  const response = await fetch(`/posts/${postId}/content.md`);
  return await response.text();
}

async function fetchMedia(postId, imageName) {
  const response = await fetch(`/media/${postId}/${imageName}`);
  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  return url;
}

function init(wasmObj) {
  go.run(wasmObj.instance);
  //!fetchPostContent(postId)
  //!fetchMedia(postId, imageName)
  //! renderPostContent(content)
  //setEventListeners()
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