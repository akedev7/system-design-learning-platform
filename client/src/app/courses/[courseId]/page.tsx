"use client";

import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import Link from "next/link";
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

interface CourseProgress {
  courseId: number;
  totalModules: number;
  completedModules: number;
  totalLessons: number;
  completedLessons: number;
  completionPercentage: number;
}

const fetchCourseProgress = async (courseId: string): Promise<CourseProgress> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/courses/${courseId}/progress`);
  return response.data;
};

interface ResumeLesson {
  moduleId: number;
  lessonId: number;
  passed: boolean;
}

const fetchResumeLesson = async (courseId: string): Promise<ResumeLesson> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/courses/${courseId}/resume`);
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

  const { data: progress } = useQuery({
    queryKey: ["courseProgress", courseId],
    queryFn: () => fetchCourseProgress(courseId),
    enabled: !!courseId,
  });

  const { data: resumeLesson } = useQuery({
    queryKey: ["resumeLesson", courseId],
    queryFn: () => fetchResumeLesson(courseId),
    enabled: !!courseId,
    retry: false,
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

        {progress && progress.completionPercentage > 0 && (
          <div className="mb-8 p-4 bg-zinc-100 dark:bg-zinc-800 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">Your Progress</span>
              <span className="text-sm font-medium text-zinc-700 dark:text-zinc-300">{progress.completionPercentage}%</span>
            </div>
            <div className="w-full bg-zinc-300 dark:bg-zinc-600 rounded-full h-2">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all"
                style={{ width: `${progress.completionPercentage}%` }}
              />
            </div>
            <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
              {progress.completedLessons} of {progress.totalLessons} lessons completed
            </p>
            {resumeLesson && (
              <Link
                href={`/courses/${courseId}/modules/${resumeLesson.moduleId}/lessons/${resumeLesson.lessonId}`}
                className="mt-3 inline-block px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                {resumeLesson.passed ? "Review" : "Continue"} Learning
              </Link>
            )}
          </div>
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