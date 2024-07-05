'use strict';

let wasm;
const go = new Go();

async function fetchPost(postId) {
  try {
    const response = await fetch(`/posts/${postId}/content.md`);
    if (!response.ok) {
      // Your existing JavaScript code for handling non-OK responses goes here
      const postContainer = document.getElementById('post-container');
      if (postContainer) {
        postContainer.innerHTML = '<p>Loading...</p>';
      } else {
        console.error('Failed to fetch posts');
      }
    } else {
      // Handle the successful response here
      const data = await response.text();
      const postContainer = document.getElementById('post-container');
      if (postContainer) {
        postContainer.innerHTML = data;
      } else {
        console.error('No container to display the post.');
      }
    }
  } catch (error) {
    // Handle any errors that occur during the fetch operation or in the try block itself
    console.error('Error fetching the post:', error);
  } finally {
    // This block will run whether the function completes successfully, throws an error, or is rejected
    console.log('Fetch operation completed.');
  }
};

//function setEventListeners() {}

function init(wasmObj) {
  go.run(wasmObj.instance);
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