import { useState } from 'react'
import axiosClient from '../../api/axiosConfig'
import Button from '../Button'

const Register = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setMessage('')
    try {
      await axiosClient.post('/registeruser', { email, password })
      setMessage('Registration successful')
    } catch (error) {
      setMessage(error.response?.data?.error || 'Registration failed')
    }
  }

  return (
    <div className="flex justify-center py-16">
      <form
        onSubmit={handleSubmit}
        className="flex w-full max-w-sm flex-col gap-4 rounded-md bg-neutral-900 p-8 text-white"
      >
        <h1 className="text-2xl font-bold">Register</h1>

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="rounded-md bg-neutral-800 px-4 py-2 outline-none focus:ring-2 focus:ring-red-600"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="rounded-md bg-neutral-800 px-4 py-2 outline-none focus:ring-2 focus:ring-red-600"
        />

        {message && <p className="text-sm text-gray-300">{message}</p>}

        <Button type="submit" className="w-full">Register</Button>
      </form>
    </div>
  )
}

export default Register
