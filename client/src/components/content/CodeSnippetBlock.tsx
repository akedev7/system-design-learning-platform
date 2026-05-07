import { CodeSnippetConfig } from "@/lib/content-types";

export function CodeSnippetBlock({ config }: { config: CodeSnippetConfig }) {
  return (
    <div data-testid="code-block" className="my-6 rounded-lg overflow-hidden border border-zinc-200 dark:border-zinc-800">
      {(config.filename || config.language) && (
        <div className="bg-zinc-100 dark:bg-zinc-900 px-4 py-2 text-sm text-zinc-600 dark:text-zinc-400 border-b border-zinc-200 dark:border-zinc-800">
          {config.filename && <span className="mr-4">{config.filename}</span>}
          {config.language && <span className="text-xs uppercase">{config.language}</span>}
        </div>
      )}
      <pre className="bg-zinc-950 dark:bg-black p-4 overflow-x-auto">
        <code className="text-sm font-mono text-zinc-50">{config.code}</code>
      </pre>
    </div>
  );
}