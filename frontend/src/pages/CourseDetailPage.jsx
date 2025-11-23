import { useEffect, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { CourseService } from "../services/course-service";

export const CourseDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [course, setCourse] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCourse = async () => {
      try {
        setLoading(true);
        const data = await CourseService.getById(id);
        setCourse(data.course);
      } catch (err) {
        console.error("Failed to fetch course:", err);
        if (err.response && err.response.status === 404) {
          setError("Курс не найден");
        } else {
          setError("Не удалось загрузить детали курса.");
        }
      } finally {
        setLoading(false);
      }
    };

    fetchCourse();
  }, [id]);

  const handleDelete = async () => {
    if (window.confirm("Вы уверены, что хотите удалить этот курс?")) {
      try {
        await CourseService.delete(id);
        navigate("/");
      } catch (err) {
        console.error("Failed to delete course:", err);
        alert("Не удалось удалить курс.");
      }
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error || !course) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        <div className="text-red-600 mb-4">{error || "Курс не найден"}</div>
        <Link to="/" className="text-blue-600 hover:underline">
          Назад к списку
        </Link>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <div className="mb-6">
        <Link
          to="/"
          className="text-gray-500 hover:text-gray-700 flex items-center"
        >
          Назад к списку
        </Link>
      </div>

      <div className="bg-white shadow-lg rounded-lg overflow-hidden border border-gray-200">
        <div className="p-8">
          <div className="flex justify-between items-start mb-6">
            <h1 className="text-3xl font-bold text-gray-900">{course.title}</h1>
            <span className="text-sm text-gray-500 bg-gray-100 px-3 py-1 rounded-full">
              ID: {course.id}
            </span>
          </div>

          <div className="prose max-w-none mb-8">
            <h3 className="text-lg font-semibold text-gray-700 mb-2">
              Описание
            </h3>
            <p className="text-gray-600 whitespace-pre-wrap">
              {course.description}
            </p>
          </div>

          <div className="border-t border-gray-100 pt-6 mt-6">
            <div className="flex justify-between items-center">
              <div className="text-sm text-gray-500">
                Создан: {new Date(course.createdAt).toLocaleDateString()}
              </div>
              <div className="space-x-3">
                <Link
                  to={`/courses/${course.id}/edit`}
                  className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
                >
                  Редактировать
                </Link>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-white text-red-600 border border-red-200 rounded hover:bg-red-50 transition-colors"
                >
                  Удалить
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
