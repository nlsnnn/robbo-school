import { apiClient } from "../api/client";

const handleResponse = (response) => {
  if (response.data.status === "Error") {
    throw new Error(response.data.error || "Произошла ошибка");
  }
  return response.data;
};

export class CourseService {
  static async getAll() {
    const response = await apiClient.get("/courses");
    return handleResponse(response);
  }

  static async getById(id) {
    const response = await apiClient.get(`/courses/${id}`);
    return handleResponse(response);
  }

  static async create(data) {
    const response = await apiClient.post("/courses", data);
    return handleResponse(response);
  }

  static async update(id, data) {
    const response = await apiClient.put(`/courses/${id}`, data);
    return handleResponse(response);
  }

  static async delete(id) {
    const response = await apiClient.delete(`/courses/${id}`);
    return handleResponse(response);
  }
}
