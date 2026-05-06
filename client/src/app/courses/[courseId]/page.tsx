"use client";

import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface Module {
  id: number;
  courseId: number;
  title: string;
  description: string;
  orderIndex: number;
}

const fetchModules = async (courseId: string): Promise<Module[]> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/courses/${courseId}/modules`);
  return response.data;
};

interface Course {
  id: number;
  title: string;
  description: string;
}

const fetchCourse = async (courseId: string): Promise<Course> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/courses/${courseId}`);
  return response.data;
};

export default function CoursePage() {
  const params = useParams();
  const courseId = params.courseId as string;

  const { data: course } = useQuery({
    queryKey: ["course", courseId],
    queryFn: () => fetchCourse(courseId),
    enabled: !!courseId,
  });

  const { data: modules, isLoading } = useQuery({
    queryKey: ["modules", courseId],
    queryFn: () => fetchModules(courseId),
    enabled: !!courseId,
  });

  return (
    <div className="flex flex-col flex-1 bg-zinc-50 font-sans dark:bg-black">
      <main className="flex flex-1 w-full max-w-3xl flex-col py-16 px-4 mx-auto bg-white dark:bg-black sm:py-32 sm:px-16">
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50 mb-4">
          {course?.title || "Loading..."}
        </h1>
        {course?.description && (
          <p className="text-zinc-600 dark:text-zinc-400 mb-8">{course.description}</p>
        )}
        <h2 className="text-2xl font-medium text-black dark:text-zinc-50 mb-4">Modules</h2>
        {isLoading && <p className="text-zinc-600 dark:text-zinc-400">Loading modules...</p>}
        {!isLoading && (
          <div className="space-y-4">
            {modules?.length === 0 ? (
              <p className="text-zinc-600 dark:text-zinc-400">No modules available yet.</p>
            ) : (
              <ul className="space-y-4">
                {modules?.map((module) => (
                  <li key={module.id}>
                    <a
                      href={`/courses/${courseId}/modules/${module.id}`}
                      className="block p-6 border border-zinc-200 dark:border-zinc-800 rounded-lg hover:border-zinc-400 dark:hover:border-zinc-600 transition-colors"
                    >
                      <h3 className="text-xl font-medium text-black dark:text-zinc-50">
                        {module.orderIndex + 1}. {module.title}
                      </h3>
                      {module.description && (
                        <p className="mt-2 text-zinc-600 dark:text-zinc-400">
                          {module.description}
                        </p>
                      )}
                    </a>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </main>
    </div>
  );
}