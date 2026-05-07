"use client";

import { useState, useEffect, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import axios from "axios";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface ContentBlock {
  type: "Text" | "Image" | "CodeSnippet" | "Quiz" | "ReactFlowDiagram";
  order: number;
  config: Record<string, unknown>;
}

interface Lesson {
  id: number;
  title: string;
  contentJsonb: ContentBlock[];
}

interface SortableBlockProps {
  block: ContentBlock;
  index: number;
  onRemove: (index: number) => void;
  onUpdate: (index: number, config: Record<string, unknown>) => void;
}

function SortableBlock({ block, index, onRemove, onUpdate }: SortableBlockProps) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id: `block-${index}` });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style} className="bg-white border rounded-lg p-4 mb-3">
      <div className="flex items-center justify-between mb-3">
        <span className="font-medium text-zinc-700">{block.type} Block</span>
        <div className="flex gap-2">
          <button
            {...attributes}
            {...listeners}
            className="cursor-grab px-2 py-1 text-zinc-500 hover:text-zinc-700"
          >
            ⬆⬇
          </button>
          <button
            onClick={() => onRemove(index)}
            className="px-2 py-1 text-red-500 hover:text-red-700"
          >
            ×
          </button>
        </div>
      </div>
      <BlockEditor block={block} index={index} onUpdate={onUpdate} />
    </div>
  );
}

function BlockEditor({ block, index, onUpdate }: { block: ContentBlock; index: number; onUpdate: (index: number, config: Record<string, unknown>) => void }) {
  switch (block.type) {
    case "Text":
      return (
        <TextBlockEditor
          config={block.config as { content: string }}
          onChange={(config) => onUpdate(index, config)}
        />
      );
    case "Image":
      return (
        <ImageBlockEditor
          config={block.config as { src: string; alt: string; caption?: string }}
          onChange={(config) => onUpdate(index, config)}
        />
      );
    case "CodeSnippet":
      return (
        <CodeSnippetBlockEditor
          config={block.config as { code: string; language: string; filename?: string }}
          onChange={(config) => onUpdate(index, config)}
        />
      );
    case "Quiz":
      return (
        <QuizBlockEditor
          config={block.config as { questions: unknown[] }}
          onChange={(config) => onUpdate(index, config)}
        />
      );
    case "ReactFlowDiagram":
      return (
        <DiagramBlockEditor
          config={block.config as { nodeTypes: Record<string, number>; edges: unknown[] }}
          onChange={(config) => onUpdate(index, config)}
        />
      );
    default:
      return <div>Unknown block type</div>;
  }
}

function TextBlockEditor({ config, onChange }: { config: { content: string }; onChange: (config: Record<string, unknown>) => void }) {
  return (
    <div>
      <label className="block text-sm font-medium mb-1">Content</label>
      <textarea
        value={config.content || ""}
        onChange={(e) => onChange({ content: e.target.value })}
        className="w-full px-3 py-2 border rounded"
        rows={4}
        placeholder="Enter text content..."
      />
    </div>
  );
}

function ImageBlockEditor({ config, onChange }: { config: { src: string; alt: string; caption?: string }; onChange: (config: Record<string, unknown>) => void }) {
  const [uploading, setUploading] = useState(false);

  const handleUrlChange = (field: string, value: string) => {
    onChange({ ...config, [field]: value });
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    try {
      const res = await fetch("/api/v1/uploads/generate-url", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          filename: file.name,
          contentType: file.type,
        }),
      });
      const data = await res.json();
      if (data.data?.uploadURL) {
        await fetch(data.data.uploadURL, {
          method: "PUT",
          body: file,
          headers: { "Content-Type": file.type },
        });
        onChange({ ...config, src: data.data.publicURL });
      }
    } catch (err) {
      console.error("Upload failed:", err);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-3">
      <div>
        <label className="block text-sm font-medium mb-1">Upload Image</label>
        <input
          type="file"
          accept="image/*"
          onChange={handleFileUpload}
          disabled={uploading}
          className="w-full px-3 py-2 border rounded"
        />
        {uploading && <p className="text-sm text-blue-600 mt-1">Uploading...</p>}
      </div>
      <div className="flex items-center gap-2">
        <span className="text-sm text-gray-500">or</span>
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">Image URL</label>
        <input
          type="text"
          value={config.src || ""}
          onChange={(e) => handleUrlChange("src", e.target.value)}
          className="w-full px-3 py-2 border rounded"
          placeholder="https://..."
        />
      </div>
      {config.src && (
        <div className="mt-2">
          <img src={config.src} alt="Preview" className="max-h-40 rounded" />
        </div>
      )}
      <div>
        <label className="block text-sm font-medium mb-1">Alt Text</label>
        <input
          type="text"
          value={config.alt || ""}
          onChange={(e) => handleUrlChange("alt", e.target.value)}
          className="w-full px-3 py-2 border rounded"
          placeholder="Image description"
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">Caption (optional)</label>
        <input
          type="text"
          value={config.caption || ""}
          onChange={(e) => handleUrlChange("caption", e.target.value)}
          className="w-full px-3 py-2 border rounded"
          placeholder="Caption"
        />
      </div>
    </div>
  );
}

function CodeSnippetBlockEditor({ config, onChange }: { config: { code: string; language: string; filename?: string }; onChange: (config: Record<string, unknown>) => void }) {
  const handleChange = (field: string, value: string) => {
    onChange({ ...config, [field]: value });
  };

  return (
    <div className="space-y-3">
      <div>
        <label className="block text-sm font-medium mb-1">Language</label>
        <select
          value={config.language || "javascript"}
          onChange={(e) => handleChange("language", e.target.value)}
          className="w-full px-3 py-2 border rounded"
        >
          <option value="javascript">JavaScript</option>
          <option value="typescript">TypeScript</option>
          <option value="python">Python</option>
          <option value="go">Go</option>
          <option value="sql">SQL</option>
          <option value="bash">Bash</option>
          <option value="json">JSON</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">Filename (optional)</label>
        <input
          type="text"
          value={config.filename || ""}
          onChange={(e) => handleChange("filename", e.target.value)}
          className="w-full px-3 py-2 border rounded"
          placeholder="example.ts"
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">Code</label>
        <textarea
          value={config.code || ""}
          onChange={(e) => handleChange("code", e.target.value)}
          className="w-full px-3 py-2 border rounded font-mono text-sm"
          rows={6}
          placeholder="// Your code here..."
        />
      </div>
    </div>
  );
}

interface QuizQuestion {
  id: number;
  type: "multiple_choice" | "true_false";
  options?: string[];
  correct: string;
}

function QuizBlockEditor({ config, onChange }: { config: { questions: QuizQuestion[] }; onChange: (config: Record<string, unknown>) => void }) {
  const questions = config.questions || [];

  const addQuestion = () => {
    const newQuestion: QuizQuestion = {
      id: Date.now(),
      type: "multiple_choice",
      options: ["A", "B", "C", "D"],
      correct: "A",
    };
    onChange({ questions: [...questions, newQuestion] });
  };

  const updateQuestion = (index: number, updates: Partial<QuizQuestion>) => {
    const updated = [...questions];
    updated[index] = { ...updated[index], ...updates };
    onChange({ questions: updated });
  };

  const removeQuestion = (index: number) => {
    onChange({ questions: questions.filter((_, i) => i !== index) });
  };

  return (
    <div className="space-y-4">
      {questions.map((q, i) => (
        <div key={q.id} className="border p-3 rounded">
          <div className="flex justify-between items-center mb-2">
            <span className="font-medium">Question {i + 1}</span>
            <button onClick={() => removeQuestion(i)} className="text-red-500">×</button>
          </div>
          <div className="mb-2">
            <label className="text-sm">Type:</label>
            <select
              value={q.type}
              onChange={(e) => updateQuestion(i, { type: e.target.value as "multiple_choice" | "true_false" })}
              className="ml-2 px-2 py-1 border rounded"
            >
              <option value="multiple_choice">Multiple Choice</option>
              <option value="true_false">True/False</option>
            </select>
          </div>
          {q.type === "multiple_choice" && q.options && (
            <div className="space-y-2">
              {q.options.map((opt, j) => (
                <div key={j} className="flex items-center gap-2">
                  <input
                    type="radio"
                    name={`q${i}-correct`}
                    checked={q.correct === opt}
                    onChange={() => updateQuestion(i, { correct: opt })}
                  />
                  <input
                    type="text"
                    value={opt}
                    onChange={(e) => {
                      const newOptions = [...(q.options || [])];
                      newOptions[j] = e.target.value;
                      updateQuestion(i, { options: newOptions });
                    }}
                    className="flex-1 px-2 py-1 border rounded text-sm"
                  />
                </div>
              ))}
            </div>
          )}
          {q.type === "true_false" && (
            <div className="flex gap-4">
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  name={`q${i}-correct`}
                  checked={q.correct === "true"}
                  onChange={() => updateQuestion(i, { correct: "true" })}
                />
                True
              </label>
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  name={`q${i}-correct`}
                  checked={q.correct === "false"}
                  onChange={() => updateQuestion(i, { correct: "false" })}
                />
                False
              </label>
            </div>
          )}
        </div>
      ))}
      <button onClick={addQuestion} className="px-3 py-1 bg-blue-100 text-blue-700 rounded text-sm">
        + Add Question
      </button>
    </div>
  );
}

interface DiagramEdge {
  from: string;
  to: string;
}

interface DiagramConfig {
  nodeTypes: Record<string, number>;
  edges: DiagramEdge[];
}

function DiagramBlockEditor({ config, onChange }: { config: DiagramConfig; onChange: (config: Record<string, unknown>) => void }) {
  const nodeTypes = config.nodeTypes || {};
  const edges = config.edges || [];
  const [newNodeType, setNewNodeType] = useState("AppServer");
  const [newNodeCount, setNewNodeCount] = useState(1);

  const addNodeType = () => {
    onChange({
      nodeTypes: { ...nodeTypes, [newNodeType]: newNodeCount },
      edges,
    });
  };

  const removeNodeType = (type: string) => {
    const updated = { ...nodeTypes };
    delete updated[type];
    onChange({ nodeTypes: updated, edges });
  };

  const addEdge = (from: string, to: string) => {
    onChange({ nodeTypes, edges: [...edges, { from, to }] });
  };

  const removeEdge = (index: number) => {
    onChange({ nodeTypes, edges: edges.filter((_, i) => i !== index) });
  };

  const availableTypes = ["Client", "APIGateway", "LoadBalancer", "AppServer", "Cache", "Database", "CDN"];

  return (
    <div className="space-y-4">
      <div className="flex gap-2 items-end">
        <div>
          <label className="block text-sm font-medium mb-1">Node Type</label>
          <select
            value={newNodeType}
            onChange={(e) => setNewNodeType(e.target.value)}
            className="px-2 py-1 border rounded"
          >
            {availableTypes.map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Count</label>
          <input
            type="number"
            value={newNodeCount}
            onChange={(e) => setNewNodeCount(parseInt(e.target.value) || 1)}
            className="w-16 px-2 py-1 border rounded"
            min={1}
          />
        </div>
        <button onClick={addNodeType} className="px-3 py-1 bg-blue-100 text-blue-700 rounded">
          Add
        </button>
      </div>

      <div className="flex flex-wrap gap-2">
        {Object.entries(nodeTypes).map(([type, count]) => (
          <div key={type} className="flex items-center gap-1 bg-zinc-100 px-2 py-1 rounded">
            <span className="text-sm">{type}: {count}</span>
            <button onClick={() => removeNodeType(type)} className="text-red-500">×</button>
          </div>
        ))}
      </div>

      <div className="border-t pt-3">
        <h4 className="font-medium mb-2">Expected Connections</h4>
        <div className="flex gap-2 mb-2">
          <select id="edge-from" className="px-2 py-1 border rounded">
            {Object.keys(nodeTypes).map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
          <span>→</span>
          <select id="edge-to" className="px-2 py-1 border rounded">
            {Object.keys(nodeTypes).map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
          <button
            onClick={() => {
              const from = (document.getElementById("edge-from") as HTMLSelectElement).value;
              const to = (document.getElementById("edge-to") as HTMLSelectElement).value;
              if (from && to) addEdge(from, to);
            }}
            className="px-2 py-1 bg-blue-100 text-blue-700 rounded"
          >
            Add
          </button>
        </div>
        <div className="flex flex-wrap gap-2">
          {edges.map((e, i) => (
            <div key={i} className="flex items-center gap-1 bg-zinc-100 px-2 py-1 rounded">
              <span className="text-sm">{e.from} → {e.to}</span>
              <button onClick={() => removeEdge(i)} className="text-red-500">×</button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default function EditLessonPage() {
  const params = useParams();
  const router = useRouter();
  const lessonId = params.lessonId as string;

  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [blocks, setBlocks] = useState<ContentBlock[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  );

  useEffect(() => {
    if (lessonId) {
      fetchLesson();
    }
  }, [lessonId]);

  const fetchLesson = async () => {
    try {
      const res = await axios.get<Lesson>(`${API_BASE_URL}/api/v1/lessons/${lessonId}`);
      setLesson(res.data);
      setBlocks(res.data.contentJsonb || []);
    } catch (error) {
      console.error("Failed to fetch lesson:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await axios.put(`${API_BASE_URL}/api/v1/lessons/${lessonId}/content`, {
        contentJsonb: blocks,
      });
      alert("Content saved!");
    } catch (error) {
      console.error("Failed to save:", error);
      alert("Failed to save");
    } finally {
      setSaving(false);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (active.id !== over?.id) {
      setBlocks((items) => {
        const oldIndex = items.findIndex((_, i) => `block-${i}` === active.id);
        const newIndex = items.findIndex((_, i) => `block-${i}` === over?.id);
        return arrayMove(items, oldIndex, newIndex);
      });
    }
  };

  const addBlock = (type: ContentBlock["type"]) => {
    const newBlock: ContentBlock = {
      type,
      order: blocks.length,
      config: getDefaultConfig(type),
    };
    setBlocks([...blocks, newBlock]);
  };

  const removeBlock = (index: number) => {
    setBlocks(blocks.filter((_, i) => i !== index));
  };

  const updateBlock = (index: number, config: Record<string, unknown>) => {
    const updated = [...blocks];
    updated[index] = { ...updated[index], config };
    setBlocks(updated);
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div className="min-h-screen bg-zinc-50 p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <button onClick={() => router.push("/admin")} className="text-zinc-600 hover:text-zinc-800 mb-2">
              ← Back to Admin
            </button>
            <h1 className="text-2xl font-semibold">Edit: {lesson?.title}</h1>
          </div>
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-zinc-400"
          >
            {saving ? "Saving..." : "Save Content"}
          </button>
        </div>

        <div className="mb-6">
          <h3 className="font-medium mb-2">Add Block</h3>
          <div className="flex flex-wrap gap-2">
            <button onClick={() => addBlock("Text")} className="px-3 py-2 bg-zinc-200 rounded hover:bg-zinc-300">
              + Text
            </button>
            <button onClick={() => addBlock("Image")} className="px-3 py-2 bg-zinc-200 rounded hover:bg-zinc-300">
              + Image
            </button>
            <button onClick={() => addBlock("CodeSnippet")} className="px-3 py-2 bg-zinc-200 rounded hover:bg-zinc-300">
              + Code
            </button>
            <button onClick={() => addBlock("Quiz")} className="px-3 py-2 bg-zinc-200 rounded hover:bg-zinc-300">
              + Quiz
            </button>
            <button onClick={() => addBlock("ReactFlowDiagram")} className="px-3 py-2 bg-zinc-200 rounded hover:bg-zinc-300">
              + Diagram
            </button>
          </div>
        </div>

        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={blocks.map((_, i) => `block-${i}`)} strategy={verticalListSortingStrategy}>
            {blocks.map((block, index) => (
              <SortableBlock key={`block-${index}`} block={block} index={index} onRemove={removeBlock} onUpdate={updateBlock} />
            ))}
          </SortableContext>
        </DndContext>

        {blocks.length === 0 && (
          <div className="text-center py-12 text-zinc-400 border-2 border-dashed rounded-lg">
            No content blocks yet. Add a block above to get started.
          </div>
        )}
      </div>
    </div>
  );
}

function getDefaultConfig(type: string): Record<string, unknown> {
  switch (type) {
    case "Text":
      return { content: "" };
    case "Image":
      return { src: "", alt: "", caption: "" };
    case "CodeSnippet":
      return { code: "", language: "javascript", filename: "" };
    case "Quiz":
      return { questions: [] };
    case "ReactFlowDiagram":
      return { nodeTypes: {}, edges: [] };
    default:
      return {};
  }
}