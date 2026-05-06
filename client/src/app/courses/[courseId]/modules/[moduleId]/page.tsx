"use client";

import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface Lesson {
  id: number;
  moduleId: number;
  title: string;
  description: string;
  orderIndex: number;
}

const fetchLessons = async (moduleId: string): Promise<Lesson[]> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/modules/${moduleId}/lessons`);
  return response.data;
};

interface Module {
  id: number;
  courseId: number;
  title: string;
  description: string;
}

const fetchModule = async (moduleId: string): Promise<Module> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/modules/${moduleId}`);
  return response.data;
};

export default function ModulePage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const moduleId = params.moduleId as string;

  const { data: module } = useQuery({
    queryKey: ["module", moduleId],
    queryFn: () => fetchModule(moduleId),
    enabled: !!moduleId,
  });

  const { data: lessons, isLoading } = useQuery({
    queryKey: ["lessons", moduleId],
    queryFn: () => fetchLessons(moduleId),
    enabled: !!moduleId,
  });

  return (
    <div className="flex flex-col flex-1 bg-zinc-50 font-sans dark:bg-black">
      <main className="flex flex-1 w-full max-w-3xl flex-col py-16 px-4 mx-auto bg-white dark:bg-black sm:py-32 sm:px-16">
        <a
          href={`/courses/${courseId}`}
          className="text-zinc-600 dark:text-zinc-400 hover:text-black dark:hover:text-zinc-50 mb-4"
        >
          &larr; Back to course
        </a>
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50 mb-4">
          {module?.title || "Loading..."}
        </h1>
        {module?.description && (
          <p className="text-zinc-600 dark:text-zinc-400 mb-8">{module.description}</p>
        )}
        <h2 className="text-2xl font-medium text-black dark:text-zinc-50 mb-4">Lessons</h2>
        {isLoading && <p className="text-zinc-600 dark:text-zinc-400">Loading lessons...</p>}
        {!isLoading && (
          <div className="space-y-4">
            {lessons?.length === 0 ? (
              <p className="text-zinc-600 dark:text-zinc-400">No lessons available yet.</p>
            ) : (
              <ul className="space-y-4">
                {lessons?.map((lesson) => (
                  <li key={lesson.id}>
                    <a
                      href={`/courses/${courseId}/modules/${moduleId}/lessons/${lesson.id}`}
                      className="block p-6 border border-zinc-200 dark:border-zinc-800 rounded-lg hover:border-zinc-400 dark:hover:border-zinc-600 transition-colors"
                    >
                      <h3 className="text-xl font-medium text-black dark:text-zinc-50">
                        {lesson.orderIndex + 1}. {lesson.title}
                      </h3>
                      {lesson.description && (
                        <p className="mt-2 text-zinc-600 dark:text-zinc-400">
                          {lesson.description}
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