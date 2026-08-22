import React from 'react'
import { Routes, Route } from 'react-router-dom'
import Header from './header/Header'
import Home from './components/home/Home'
import Login from './components/login/Login'
import Register from './components/register/Register'


const App = () => {
  return (
    <div>
      <Header/>
      <Routes>
        <Route path="/" element={<Home/>} />
        <Route path="/movies" element={<Home/>} />
        <Route path="/login" element={<Login />} />
        <Route path='/registeruser' element={<Register/>}/>
      </Routes>
    </div>
  )
}

export default App
