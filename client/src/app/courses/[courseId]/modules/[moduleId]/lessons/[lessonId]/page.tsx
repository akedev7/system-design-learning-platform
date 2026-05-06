"use client";

import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import axios from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface Lesson {
  id: number;
  moduleId: number;
  title: string;
  description: string | null;
  contentJsonb: ContentBlock[];
}

const fetchLesson = async (lessonId: string): Promise<Lesson> => {
  const response = await axios.get(`${API_BASE_URL}/api/v1/lessons/${lessonId}`);
  return response.data;
};

import { ContentBlock, TextConfig, ImageConfig, CodeSnippetConfig } from "@/lib/content-types";
import { TextBlock } from "@/components/content/TextBlock";
import { ImageBlock } from "@/components/content/ImageBlock";
import { CodeSnippetBlock } from "@/components/content/CodeSnippetBlock";

function ContentBlockRenderer({ block }: { block: ContentBlock }) {
  switch (block.type) {
    case "Text":
      return <TextBlock config={block.config as TextConfig} />;
    case "Image":
      return <ImageBlock config={block.config as ImageConfig} />;
    case "CodeSnippet":
      return <CodeSnippetBlock config={block.config as CodeSnippetConfig} />;
    default:
      return null;
  }
}

export default function LessonPage() {
  const params = useParams();
  const courseId = params.courseId as string;
  const moduleId = params.moduleId as string;
  const lessonId = params.lessonId as string;

  const { data: lesson, isLoading } = useQuery({
    queryKey: ["lesson", lessonId],
    queryFn: () => fetchLesson(lessonId),
    enabled: !!lessonId,
  });

  return (
    <div className="flex flex-col flex-1 bg-zinc-50 font-sans dark:bg-black">
      <main className="flex flex-1 w-full max-w-3xl flex-col py-16 px-4 mx-auto bg-white dark:bg-black sm:py-32 sm:px-16">
        <a
          href={`/courses/${courseId}/modules/${moduleId}`}
          className="text-zinc-600 dark:text-zinc-400 hover:text-black dark:hover:text-zinc-50 mb-4"
        >
          &larr; Back to module
        </a>
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50 mb-4">
          {lesson?.title || "Loading..."}
        </h1>
        {lesson?.description && (
          <p className="text-zinc-600 dark:text-zinc-400 mb-8">{lesson.description}</p>
        )}
        {isLoading && <p className="text-zinc-600 dark:text-zinc-400">Loading content...</p>}
        {!isLoading && lesson?.contentJsonb && (
          <div className="space-y-6">
            {lesson.contentJsonb.map((block, index) => (
              <ContentBlockRenderer key={index} block={block} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}