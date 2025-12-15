import { useState, useEffect } from 'react'

type Status = 'pending' | 'downloading' | 'analyzing' | 'processing' | 'completed' | 'failed'

interface JobStatusResponse {
  id: string
  status: Status
  message: string
}

function App() {
  const [url, setUrl] = useState('')
  const [jobId, setJobId] = useState<string | null>(null)
  const [status, setStatus] = useState<Status | null>(null)
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)

  const handleGenerate = async () => {
    if (!url) return
    setLoading(true)
    setStatus('pending')
    setMessage('Starting...')
    
    try {
      const res = await fetch('http://localhost:8080/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url })
      })
      if (!res.ok) throw new Error('Failed to start')
      
      const data = await res.json()
      setJobId(data.job_id)
    } catch (e) {
      console.error(e)
      setStatus('failed')
      setMessage('Failed to start job')
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!jobId || status === 'completed' || status === 'failed') return

    const interval = setInterval(async () => {
      try {
        const res = await fetch(`http://localhost:8080/status/${jobId}`)
        if (res.ok) {
          const data: JobStatusResponse = await res.json()
          setStatus(data.status)
          setMessage(data.message)
          
          if (data.status === 'completed' || data.status === 'failed') {
            setLoading(false)
          }
        }
      } catch (e) {
        console.error("Poll error", e)
      }
    }, 2000)

    return () => clearInterval(interval)
  }, [jobId, status])

  const steps = [
    { key: 'downloading', label: 'Downloading Video' },
    { key: 'analyzing', label: 'Analyzing Content' },
    { key: 'processing', label: 'Processing Clip' },
  ]

  const getStepClass = (stepKey: string) => {
    // simple linear progression logic for UI
    const order = ['pending', 'downloading', 'analyzing', 'processing', 'completed']
    const currentIndex = order.indexOf(status || 'pending')
    const stepIndex = order.indexOf(stepKey)
    
    if (status === 'failed') return 'status-check failed'
    if (currentIndex > stepIndex) return 'status-check done'
    if (currentIndex === stepIndex) return 'status-check active'
    return 'status-check'
  }

  return (
    <div className="container">
      <div>
        <h1>Auto-Clip 🎬</h1>
        <p>Turn YouTube Videos into TikToks instantly.</p>
      </div>

      <div className="input-group">
        <input 
          type="text" 
          placeholder="Paste YouTube URL here..." 
          value={url}
          onChange={e => setUrl(e.target.value)}
          disabled={loading}
        />
        <button onClick={handleGenerate} disabled={loading || !url}>
          {loading ? 'Processing...' : 'Generate Clips'}
        </button>
      </div>

      {status && (
        <div className="status-card">
          <div className="status-title">Status: {status.toUpperCase()}</div>
          <p>{message}</p>
          
          <div style={{marginTop: '1rem'}}>
             {steps.map(s => (
               <div key={s.key} className={getStepClass(s.key)}>
                 {getStepClass(s.key).includes('done') ? '✅' : 
                  getStepClass(s.key).includes('active') ? '🔄' : '⚪'} {s.label}
               </div>
             ))}
          </div>

          {status === 'completed' && (
            <div style={{marginTop: '1.5rem', textAlign: 'center'}}>
              <a 
                href={`http://localhost:8080/download/${jobId}`} 
                target="_blank" 
                className="download-link"
              >
                Download Clip ⬇️
              </a>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default App
