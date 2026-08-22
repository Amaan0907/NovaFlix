import React from 'react'
import { Routes, Route } from 'react-router-dom'
import Header from './header/Header'
import Home from './components/home/Home'


const App = () => {
  return (
    <div>
      <Header/>
      <Routes>
        <Route path="/" element={<Home/>} />
        <Route path="/movies" element={<Home/>} />
        {/* <Route path="/tv-shows" element={<ComingSoon title="TV Shows" />} />
        <Route path="/my-list" element={<ComingSoon title="My List" />} />
        <Route path="/login" element={<ComingSoon title="Sign In" />} /> */}
      </Routes>
    </div>
  )
}

export default App
