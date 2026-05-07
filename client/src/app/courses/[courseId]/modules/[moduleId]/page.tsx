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
  // API returns {status: "success", data: [...]}
  return response.data.data || [];
};

interface Module {
  id: number;
  courseId: number;
  title: string;
  description: string;
}

const fetchModule = async (moduleId: string): Promise<Module> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/modules/${moduleId}`);
  // API returns {status: "success", data: {...}}
  return response.data.data || null;
};

interface ModuleProgress {
  moduleId: number;
  totalLessons: number;
  completedLessons: number;
  completionPercentage: number;
}

const fetchModuleProgress = async (moduleId: string): Promise<ModuleProgress> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/modules/${moduleId}/progress`);
  // API returns {status: "success", data: {...}}
  return response.data.data || null;
};

interface LessonProgress {
  lessonId: number;
  score: number;
  passed: boolean;
  completedAt?: string;
}

const fetchLessonProgress = async (lessonId: string): Promise<LessonProgress> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/lessons/${lessonId}/progress`);
  // API returns {status: "success", data: {...}}
  return response.data.data || null;
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

  const { data: moduleProgress } = useQuery({
    queryKey: ["moduleProgress", moduleId],
    queryFn: () => fetchModuleProgress(moduleId),
    enabled: !!moduleId,
    retry: false,
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
          <p className="text-zinc-600 dark:text-zinc-400 mb-4">{module.description}</p>
        )}

        {moduleProgress && moduleProgress.completionPercentage > 0 && (
          <div className="mb-6 p-4 bg-zinc-100 dark:bg-zinc-800 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">Module Progress</span>
              <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">{moduleProgress.completionPercentage}%</span>
            </div>
            <div className="w-full bg-zinc-300 dark:bg-zinc-600 rounded-full h-2">
              <div
                className="bg-green-600 h-2 rounded-full transition-all"
                style={{ width: `${moduleProgress.completionPercentage}%` }}
              />
            </div>
            <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
              {moduleProgress.completedLessons} of {moduleProgress.totalLessons} lessons completed
            </p>
          </div>
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
                  <LessonItem
                    key={lesson.id}
                    lesson={lesson}
                    courseId={courseId}
                    moduleId={moduleId}
                  />
                ))}
              </ul>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

function LessonItem({ lesson, courseId, moduleId }: { lesson: Lesson; courseId: string; moduleId: string }) {
  const { data: progress } = useQuery({
    queryKey: ["lessonProgress", lesson.id],
    queryFn: () => fetchLessonProgress(String(lesson.id)),
    retry: false,
  });

  const isCompleted = progress?.passed;

  return (
    <li data-testid="lesson-card">
      <a
        href={`/courses/${courseId}/modules/${moduleId}/lessons/${lesson.id}`}
        className="block p-6 border border-zinc-200 dark:border-zinc-800 rounded-lg hover:border-zinc-400 dark:hover:border-zinc-600 transition-colors"
      >
        <div className="flex items-center gap-3">
          {isCompleted ? (
            <span className="text-green-600 text-xl">✓</span>
          ) : (
            <span className="text-zinc-400 text-xl">○</span>
          )}
          <div className="flex-1">
            <h3 className="text-xl font-medium text-black dark:text-zinc-50">
              {lesson.orderIndex + 1}. {lesson.title}
            </h3>
            {lesson.description && (
              <p className="mt-2 text-zinc-600 dark:text-zinc-400">
                {lesson.description}
              </p>
            )}
          </div>
          {progress && progress.score > 0 && (
            <span className={`text-sm font-medium ${isCompleted ? 'text-green-600' : 'text-zinc-500'}`}>
              {progress.score}%
            </span>
          )}
        </div>
      </a>
    </li>
  );
}