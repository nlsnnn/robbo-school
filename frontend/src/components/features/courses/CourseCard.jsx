import { Link } from "react-router-dom";

export const CourseCard = ({ course, onDelete }) => {
  const handleDelete = (e) => {
    e.preventDefault();
    if (window.confirm("Вы уверены, что хотите удалить этот курс?")) {
      onDelete(course.id);
    }
  };

  return (
    <div className="bg-white shadow-md rounded-lg p-6 hover:shadow-lg transition-shadow duration-300 border border-gray-200">
      <h3 className="text-xl font-semibold mb-2 text-gray-800">
        {course.title}
      </h3>
      <p className="text-gray-600 mb-4 line-clamp-3">{course.description}</p>
      <div className="flex justify-between items-center mt-4">
        <Link
          to={`/courses/${course.id}`}
          className="text-blue-600 hover:text-blue-800 font-medium"
        >
          Подробнее
        </Link>
        <div className="flex gap-2 flex-col  xl:flex-row">
          <Link
            to={`/courses/${course.id}/edit`}
            className="px-3 py-1 text-sm text-gray-600 hover:text-gray-800 border border-gray-300 rounded hover:bg-gray-50 transition-colors"
          >
            Редактировать
          </Link>
          <button
            onClick={handleDelete}
            className="px-3 py-1 text-sm text-red-600 hover:text-red-800 border border-red-200 rounded hover:bg-red-50 transition-colors"
          >
            Удалить
          </button>
        </div>
      </div>
    </div>
  );
};
