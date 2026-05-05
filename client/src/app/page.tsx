"use client";

import { useQuery } from "@tanstack/react-query";
import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface Course {
  id: number;
  title: string;
  description: string;
}

const fetchCourses = async (): Promise<Course[]> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/courses`);
  return response.data;
};

export default function Home() {
  const { data: courses, isLoading, error } = useQuery({
    queryKey: ["courses"],
    queryFn: fetchCourses,
  });

  return (
    <div className="flex flex-col flex-1 bg-zinc-50 font-sans dark:bg-black">
      <main className="flex flex-1 w-full max-w-3xl flex-col py-16 px-4 mx-auto bg-white dark:bg-black sm:py-32 sm:px-16">
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50 mb-8">
          Course Catalog
        </h1>
        {isLoading && <p className="text-zinc-600 dark:text-zinc-400">Loading courses...</p>}
        {error && <p className="text-red-500">Failed to load courses.</p>}
        {!isLoading && !error && (
          <div className="space-y-4">
            {courses?.length === 0 ? (
              <p className="text-zinc-600 dark:text-zinc-400">No courses available yet.</p>
            ) : (
              <ul className="space-y-4">
                {courses?.map((course) => (
                  <li
                    key={course.id}
                    className="p-6 border border-zinc-200 dark:border-zinc-800 rounded-lg"
                  >
                    <h2 className="text-xl font-medium text-black dark:text-zinc-50">
                      {course.title}
                    </h2>
                    <p className="mt-2 text-zinc-600 dark:text-zinc-400">
                      {course.description}
                    </p>
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
