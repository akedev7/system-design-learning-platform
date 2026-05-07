"use client";

import { useState } from "react";
import axios from "axios";
import { QuizConfig, QuizAnswerRequest, QuizResult } from "@/lib/content-types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface QuizBlockProps {
  config: QuizConfig;
  lessonId: number;
}

export function QuizBlock({ config, lessonId }: QuizBlockProps) {
  const [answers, setAnswers] = useState<Record<number, string>>({});
  const [result, setResult] = useState<QuizResult | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submitted, setSubmitted] = useState(false);

  const handleAnswerChange = (questionId: number, answer: string) => {
    if (submitted) return;
    setAnswers((prev) => ({ ...prev, [questionId]: answer }));
  };

  const handleSubmit = async () => {
    setIsSubmitting(true);
    setError(null);

    try {
      const response = await axios.post<{ result: QuizResult }>(
        `${API_BASE_URL}/api/v1/lessons/${lessonId}/submit-quiz`,
        { answers } as QuizAnswerRequest
      );
      setResult(response.data.result);
      setSubmitted(true);
    } catch (err) {
      setError("Failed to submit quiz. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const allAnswered = config.questions.every((q) => answers[q.id]);

  const getOptionClass = (questionId: number, option: string) => {
    if (!submitted) {
      return answers[questionId] === option
        ? "border-blue-500 bg-blue-50 dark:bg-blue-900/30"
        : "border-zinc-300 dark:border-zinc-700 hover:border-zinc-400 dark:hover:border-zinc-600";
    }

    const question = config.questions.find((q) => q.id === questionId);
    if (!question) return "";

    if (option === question.correct) {
      return "border-green-500 bg-green-50 dark:bg-green-900/30";
    }
    if (answers[questionId] === option && option !== question.correct) {
      return "border-red-500 bg-red-50 dark:bg-red-900/30";
    }
    return "border-zinc-300 dark:border-zinc-700 opacity-50";
  };

  return (
    <div data-testid="quiz-block" className="space-y-6 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
      <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-100">
        Quiz
      </h2>

      {config.questions.map((question) => (
        <div key={question.id} className="space-y-3">
          <p className="font-medium text-zinc-800 dark:text-zinc-200">
            {question.id}. {question.type === "multiple_choice" ? "Multiple Choice" : "True/False"}
          </p>

          {question.type === "multiple_choice" && question.options && (
            <div className="space-y-2">
              {question.options.map((option, idx) => (
                <button
                  key={idx}
                  onClick={() => handleAnswerChange(question.id, option)}
                  disabled={submitted}
                  className={`w-full rounded-md border p-3 text-left transition-colors ${getOptionClass(
                    question.id,
                    option
                  )}`}
                >
                  <span className="mr-2 font-mono text-sm text-zinc-500">
                    {String.fromCharCode(65 + idx)}.
                  </span>
                  {option}
                </button>
              ))}
            </div>
          )}

          {question.type === "true_false" && (
            <div className="flex gap-4">
              {["true", "false"].map((option) => (
                <button
                  key={option}
                  onClick={() => handleAnswerChange(question.id, option)}
                  disabled={submitted}
                  className={`flex-1 rounded-md border p-3 capitalize transition-colors ${getOptionClass(
                    question.id,
                    option
                  )}`}
                >
                  {option}
                </button>
              ))}
            </div>
          )}
        </div>
      ))}

      {!submitted && (
        <button
          onClick={handleSubmit}
          disabled={!allAnswered || isSubmitting}
          className="rounded-md bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-zinc-400 dark:bg-blue-600 dark:hover:bg-blue-700"
        >
          {isSubmitting ? "Submitting..." : "Submit Answers"}
        </button>
      )}

      {error && (
        <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
      )}

      {result && (
        <div
          className={`rounded-md border p-4 ${
            result.passed
              ? "border-green-500 bg-green-50 dark:bg-green-900/30"
              : "border-yellow-500 bg-yellow-50 dark:bg-yellow-900/30"
          }`}
        >
          <p className="font-medium">
            Score: {result.score}% ({result.correctAnswers}/{result.totalQuestions}{" "}
            correct)
          </p>
          <p
            className={`font-semibold ${
              result.passed ? "text-green-700 dark:text-green-400" : "text-yellow-700 dark:text-yellow-400"
            }`}
          >
            {result.passed ? "Passed!" : "Not quite - try again!"}
          </p>
        </div>
      )}
    </div>
  );
}