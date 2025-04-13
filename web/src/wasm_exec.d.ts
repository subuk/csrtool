declare class Go {
  importObject: WebAssembly.Imports
  run(instance: WebAssembly.Instance): Promise<void>
}

interface Window {
  generateCSR(request: string): Promise<{
    privateKey: string
    csr: string
    error?: string
  }>
}

declare global {
  const Go: {
    new(): Go
  }
}
