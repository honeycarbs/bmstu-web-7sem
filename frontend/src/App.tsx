import { useState } from 'react'
import "bootstrap/dist/css/bootstrap.min.css"
import { Navigate, Route, Routes } from 'react-router-dom'
import RegisterForm from './components/RegisterForm'
import LoginForm from './components/LoginForm'
import NoteList from './components/NoteList'
import ShowNote from './components/ShowNote'
import NewNote from './components/NewNote'
import EditNote from './components/EditNote'
import RequireAuth, { RequireAuthProps } from './components/RequireAuth'
import { QueryClient, QueryClientProvider } from 'react-query'
import NoteLayout from './components/NoteLayout'

const defaultRequireAuthProps: Omit<RequireAuthProps, 'outlet'> = {
  authenticationPath: '/login',
};

// const queryClient = new QueryClient()

function App() {

  return (
    <Routes>
      <Route path="/login" element={<LoginForm />} />
      <Route path="/register" element={<RegisterForm />} />
      <Route path="/" element={<RequireAuth {...defaultRequireAuthProps} outlet={<NoteList />} />} />
      <Route path="/new" element={<RequireAuth {...defaultRequireAuthProps} outlet={<NewNote />} />} />
      <Route path="/:id" element={<RequireAuth {...defaultRequireAuthProps} outlet={<NoteLayout />} />}>
        <Route index element={<ShowNote />} />
        <Route path="edit" element={<EditNote />} />
      </Route>
    </Routes >
  )
}

export default App
