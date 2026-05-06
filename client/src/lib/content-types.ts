export interface ContentBlock {
  type: "Text" | "Image" | "CodeSnippet";
  order: number;
  config: TextConfig | ImageConfig | CodeSnippetConfig;
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