'use strict';

let wasm;
const go = new Go();

// Function to check for featured image or video and display it at the top of the post content
function displayPost(content, postId) {
  const postContainer = document.getElementById('post-container');
  if (!postContainer) {
    console.error('No container to display the post.');
    return;
  }

  const hasFeaturedMedia = /featured\.(jpg|jpeg|png|gif|webp|mp4|avi|mov|webm)/i.test(postId);
  let displayedContent;
  if (hasFeaturedMedia) {
    const featuredImage = postId.match(/featured\.(jpg|jpeg|png|gif|webp|mp4|avi|mov|webm)/i)[0];
    displayedContent = `<div id="post-media"><img src="/posts/${postId}/${featuredImage}" alt="Featured Media"></div><div id="post-content">${content}</div>`;
  }
  displayedContent + `<div id="post-content">${content}</div>`;
  postContainer.innerHTML = displayedContent;
}

//! function buildNav();

function setEventListeners() {
  document.addEventListener("DOMContentLoaded", function () {
    let postList = fetchPostList();
    buildNav(postList);

    // Automatically load the home page by default
    fetchPost('home');
    displayPost('home');
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