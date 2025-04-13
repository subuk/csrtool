import { useState, FormEvent } from 'react'

interface CSRFormData {
  commonName: string
  keyType: string
  country: string
  state: string
  locality: string
  org: string
  orgUnit: string
  email: string
  dnsNames: string[]
  challengePassword: string
}

interface CSRResponse {
  privateKey: string
  csr: string
  error?: string
}

export default function CSRForm() {
  const [formData, setFormData] = useState<CSRFormData>({
    commonName: '',
    keyType: 'rsa2048',
    country: 'US',
    state: 'California',
    locality: 'San Francisco',
    org: 'Example Inc',
    orgUnit: 'IT',
    email: '',
    dnsNames: [],
    challengePassword: '',
  })

  const [result, setResult] = useState<CSRResponse | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      // @ts-ignore - generateCSR is added by wasm_exec.js
      const response = await window.generateCSR(JSON.stringify(formData))
      setResult(response)
    } catch (error) {
      setResult({ privateKey: '', csr: '', error: String(error) })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-white shadow sm:rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="commonName" className="block text-sm font-medium text-gray-700">
              Common Name (CN)
            </label>
            <input
              type="text"
              id="commonName"
              value={formData.commonName}
              onChange={(e) => setFormData({ ...formData, commonName: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              required
            />
          </div>

          <div>
            <label htmlFor="keyType" className="block text-sm font-medium text-gray-700">
              Key Type
            </label>
            <select
              id="keyType"
              value={formData.keyType}
              onChange={(e) => setFormData({ ...formData, keyType: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            >
              <option value="rsa2048">RSA 2048</option>
              <option value="rsa4096">RSA 4096</option>
              <option value="ec256">EC P-256</option>
              <option value="ec384">EC P-384</option>
            </select>
          </div>

          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div>
              <label htmlFor="country" className="block text-sm font-medium text-gray-700">
                Country
              </label>
              <input
                type="text"
                id="country"
                value={formData.country}
                onChange={(e) => setFormData({ ...formData, country: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                maxLength={2}
                required
              />
            </div>

            <div>
              <label htmlFor="state" className="block text-sm font-medium text-gray-700">
                State/Province
              </label>
              <input
                type="text"
                id="state"
                value={formData.state}
                onChange={(e) => setFormData({ ...formData, state: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                required
              />
            </div>

            <div>
              <label htmlFor="locality" className="block text-sm font-medium text-gray-700">
                Locality/City
              </label>
              <input
                type="text"
                id="locality"
                value={formData.locality}
                onChange={(e) => setFormData({ ...formData, locality: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                required
              />
            </div>

            <div>
              <label htmlFor="org" className="block text-sm font-medium text-gray-700">
                Organization
              </label>
              <input
                type="text"
                id="org"
                value={formData.org}
                onChange={(e) => setFormData({ ...formData, org: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                required
              />
            </div>

            <div>
              <label htmlFor="orgUnit" className="block text-sm font-medium text-gray-700">
                Organizational Unit
              </label>
              <input
                type="text"
                id="orgUnit"
                value={formData.orgUnit}
                onChange={(e) => setFormData({ ...formData, orgUnit: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                required
              />
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email
              </label>
              <input
                type="email"
                id="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <label htmlFor="challengePassword" className="block text-sm font-medium text-gray-700">
              Challenge Password
            </label>
            <input
              type="password"
              id="challengePassword"
              value={formData.challengePassword}
              onChange={(e) => setFormData({ ...formData, challengePassword: e.target.value })}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            />
          </div>

          <div>
            <button
              type="submit"
              disabled={loading}
              className="inline-flex justify-center rounded-md border border-transparent bg-primary-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50"
            >
              {loading ? 'Generating...' : 'Generate CSR'}
            </button>
          </div>
        </form>

        {result && (
          <div className="mt-6 space-y-4">
            {result.error ? (
              <div className="rounded-md bg-red-50 p-4">
                <div className="text-sm text-red-700">{result.error}</div>
              </div>
            ) : (
              <>
                <div>
                  <h3 className="text-lg font-medium text-gray-900">Private Key</h3>
                  <pre className="mt-2 whitespace-pre-wrap break-all rounded-md bg-gray-50 p-4 text-sm text-gray-900">
                    {result.privateKey}
                  </pre>
                </div>
                <div>
                  <h3 className="text-lg font-medium text-gray-900">CSR</h3>
                  <pre className="mt-2 whitespace-pre-wrap break-all rounded-md bg-gray-50 p-4 text-sm text-gray-900">
                    {result.csr}
                  </pre>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
