export interface ContentBlock {
  type: "Text" | "Image" | "CodeSnippet" | "Quiz";
  order: number;
  config: TextConfig | ImageConfig | CodeSnippetConfig | QuizConfig;
}

export interface TextConfig {
  content: string;
}

export interface ImageConfig {
  src: string;
  alt: string;
  caption?: string;
}

export interface CodeSnippetConfig {
  code: string;
  language: string;
  filename?: string;
}

export interface QuizConfig {
  questions: QuizQuestion[];
}

export interface QuizQuestion {
  id: number;
  type: "multiple_choice" | "true_false";
  options?: string[];
  correct: string;
}

export interface QuizAnswerRequest {
  answers: Record<number, string>;
}

export interface QuizResult {
  score: number;
  totalQuestions: number;
  correctAnswers: number;
  passed: boolean;
}