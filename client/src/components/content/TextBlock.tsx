import { TextConfig } from "@/lib/content-types";

export function TextBlock({ config }: { config: TextConfig }) {
  return (
    <div className="prose dark:prose-invert max-w-none">
      <p>{config.content}</p>
    </div>
  );
}