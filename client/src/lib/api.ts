import axios, { AxiosInstance, AxiosResponse } from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export interface ApiResponse<T> {
  status: "success" | "error";
  data: T;
  message?: string;
}

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.response.use(
  (response: AxiosResponse<ApiResponse<unknown>>) => {
    if (response.data.status === "error") {
      return Promise.reject(new Error(response.data.message || "API error"));
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default apiClient;
export { API_BASE_URL };