import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { CourseService } from "../services/course-service";
import { CourseForm } from "../components/features/courses/CourseForm";

export const AddCoursePage = () => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (values) => {
    try {
      setIsLoading(true);
      setError(null);
      await CourseService.create(values);
      navigate("/");
    } catch (err) {
      console.error("Failed to create course:", err);
      const message =
        err.message || "Не удалось создать курс. Пожалуйста, попробуйте снова.";
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <div className="mb-6">
        <Link to="/" className="text-gray-500 hover:text-gray-700">
          Назад к списку
        </Link>
      </div>

      <div className="bg-white shadow-md rounded-lg p-8 border border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">
          Создать новый курс
        </h1>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 text-red-700">
            <p>{error}</p>
          </div>
        )}

        <CourseForm
          onSubmit={handleSubmit}
          isLoading={isLoading}
          buttonText="Создать курс"
        />
      </div>
    </div>
  );
};
