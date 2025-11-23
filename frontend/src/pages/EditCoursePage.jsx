import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { CourseService } from "../services/course-service";
import { CourseForm } from "../components/features/courses/CourseForm";

export const EditCoursePage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [course, setCourse] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
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

  const handleSubmit = async (values) => {
    try {
      setSaving(true);
      setError(null);
      await CourseService.update(id, values);
      navigate("/");
    } catch (err) {
      console.error("Failed to update course:", err);
      const message =
        err.message ||
        "Не удалось обновить курс. Пожалуйста, попробуйте снова.";
      setError(message);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!course && !loading) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        <div className="text-red-600 mb-4">Курс не найден</div>
        <Link to="/" className="text-blue-600 hover:underline">
          Назад к списку
        </Link>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <div className="mb-6">
        <Link to="/" className="text-gray-500 hover:text-gray-700">
          Назад к списку
        </Link>
      </div>

      <div className="bg-white shadow-md rounded-lg p-8 border border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">
          Редактировать курс
        </h1>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700">
            <p>{error}</p>
          </div>
        )}

        <CourseForm
          initialValues={course}
          onSubmit={handleSubmit}
          isLoading={saving}
          buttonText="Обновить курс"
        />
      </div>
    </div>
  );
};
