import { useState } from 'react'

export default function Home() {
  const [step, setStep] = useState(1)

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold tracking-tight text-gray-900 dark:text-gray-100">
          Generate Certificate Signing Request
        </h1>
        <p className="mt-2 text-lg text-gray-600 dark:text-gray-400">
          Create a new private key and CSR in your browser
        </p>
      </div>

      <div className="mx-auto max-w-2xl">
        <div className="mb-8">
          <div className="flex items-center">
            <div className="flex items-center">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary-600">
                <span className="text-sm font-medium text-white">1</span>
              </div>
              <div className="ml-3 text-sm font-medium text-gray-900 dark:text-gray-100">
                Key Configuration
              </div>
            </div>
            <div className="ml-3 h-0.5 w-8 bg-gray-300 dark:bg-gray-700" />
            <div className="flex items-center">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-300 dark:bg-gray-700">
                <span className="text-sm font-medium text-gray-500 dark:text-gray-400">2</span>
              </div>
              <div className="ml-3 text-sm font-medium text-gray-500 dark:text-gray-400">
                Certificate Details
              </div>
            </div>
            <div className="ml-3 h-0.5 w-8 bg-gray-300 dark:bg-gray-700" />
            <div className="flex items-center">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-300 dark:bg-gray-700">
                <span className="text-sm font-medium text-gray-500 dark:text-gray-400">3</span>
              </div>
              <div className="ml-3 text-sm font-medium text-gray-500 dark:text-gray-400">
                Review & Generate
              </div>
            </div>
          </div>
        </div>

        <div className="rounded-lg bg-white p-6 shadow dark:bg-gray-800">
          {step === 1 && (
            <div className="space-y-6">
              <div>
                <label className="label">Key Type</label>
                <select className="input">
                  <option value="rsa2048">RSA 2048</option>
                  <option value="rsa4096">RSA 4096</option>
                  <option value="ec256">ECDSA P-256</option>
                  <option value="ec384">ECDSA P-384</option>
                </select>
              </div>
              <div className="flex justify-end">
                <button
                  className="btn btn-primary"
                  onClick={() => setStep(2)}
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
