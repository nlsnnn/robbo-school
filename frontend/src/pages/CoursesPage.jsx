import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { CourseService } from "../services/course-service";
import { CourseCard } from "../components/features/courses/CourseCard";

export const CoursesPage = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchCourses = async () => {
    try {
      setLoading(true);
      const data = await CourseService.getAll();
      setCourses(data.courses || []);
      setError(null);
    } catch (err) {
      console.error("Failed to fetch courses:", err);
      setError("Не удалось загрузить курсы. Пожалуйста, попробуйте позже.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  const handleDelete = async (id) => {
    try {
      await CourseService.delete(id);
      // Refresh list after delete
      fetchCourses();
    } catch (err) {
      console.error("Failed to delete course:", err);
      alert("Не удалось удалить курс. Пожалуйста, попробуйте снова.");
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-10">
        <div className="text-red-600 mb-4">{error}</div>
        <button
          onClick={fetchCourses}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Повторить
        </button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Курсы</h1>
        <Link
          to="/courses/new"
          className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors shadow-sm"
        >
          Добавить курс
        </Link>
      </div>

      {courses.length === 0 ? (
        <div className="text-center py-12 bg-gray-50 rounded-lg border border-gray-200">
          <p className="text-gray-500 text-lg mb-4">Курсы не найдены.</p>
          <Link
            to="/courses/new"
            className="text-blue-600 hover:text-blue-800 font-medium"
          >
            Создать первый курс
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {courses.map((course) => (
            <CourseCard
              key={course.id}
              course={course}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
};
