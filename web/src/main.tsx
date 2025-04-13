import './wasm_exec'
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

// Initialize WASM
const go = new Go()
WebAssembly.instantiateStreaming(fetch('/csrtool.wasm'), go.importObject).then((result) => {
  go.run(result.instance)

  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  )
})
