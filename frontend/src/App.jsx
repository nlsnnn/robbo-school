import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { CoursesPage } from './pages/CoursesPage'
import { CourseDetailPage } from './pages/CourseDetailPage'
import { AddCoursePage } from './pages/AddCoursePage'
import { EditCoursePage } from './pages/EditCoursePage'

function App() {

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <nav className="bg-white shadow-sm border-b border-gray-200">
          <div className="container mx-auto px-4 py-4 text-center">
            <a href="/" className="text-xl font-bold">Robbo School</a>
          </div>
        </nav>
        <Routes>
          <Route path="/" element={<CoursesPage />} />
          <Route path="/courses/new" element={<AddCoursePage />} />
          <Route path="/courses/:id" element={<CourseDetailPage />} />
          <Route path="/courses/:id/edit" element={<EditCoursePage />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default App
