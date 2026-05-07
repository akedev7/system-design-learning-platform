"use client";

import { useCallback, useState, useMemo } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  Node,
  Edge,
  Connection,
  addEdge,
  useNodesState,
  useEdgesState,
  MarkerType,
  NodeTypes,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import axios from "axios";
import { ReactFlowDiagramConfig, DiagramValidationResult, UserDiagram } from "@/lib/content-types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

interface ReactFlowDiagramBlockProps {
  config: ReactFlowDiagramConfig;
  lessonId: number;
}

const nodeTypeDefinitions = [
  { type: "Client", label: "Client", color: "#fcd34d" },
  { type: "APIGateway", label: "API Gateway", color: "#f472b6" },
  { type: "LoadBalancer", label: "Load Balancer", color: "#fb923c" },
  { type: "AppServer", label: "App Server", color: "#60a5fa" },
  { type: "Cache", label: "Cache (Redis)", color: "#a78bfa" },
  { type: "Database", label: "Database", color: "#34d399" },
  { type: "CDN", label: "CDN", color: "#f87171" },
];

const initialNodes: Node[] = [];
const initialEdges: Edge[] = [];

export function ReactFlowDiagramBlock({ config, lessonId }: ReactFlowDiagramBlockProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [validationResult, setValidationResult] = useState<DiagramValidationResult | null>(null);
  const [isValidating, setIsValidating] = useState(false);
  const [clientValidation, setClientValidation] = useState<{ valid: boolean; errors: string[] } | null>(null);

  const nodeTypeMap = useMemo(() => {
    const map: Record<string, number> = {};
    nodes.forEach((node) => {
      map[node.type || "default"] = (map[node.type || "default"] || 0) + 1;
    });
    return map;
  }, [nodes]);

  const customNodeTypes = useMemo(() => ({}), []);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge({ ...params, markerEnd: { type: MarkerType.ArrowClosed } }, eds)),
    [setEdges]
  );

  const addNode = (nodeType: string, label: string, color: string) => {
    const id = `${nodeType}-${Date.now()}`;
    const newNode: Node = {
      id,
      type: "default",
      position: { x: Math.random() * 400 + 100, y: Math.random() * 300 + 100 },
      data: { label: `${label} (${Object.values(nodeTypeMap).filter((_, i) => Object.keys(nodeTypeMap)[i] === nodeType).length + 1})` },
      style: { background: color, border: "2px solid #333", borderRadius: "8px", padding: "10px", width: 120 },
    };
    setNodes((nds) => [...nds, newNode]);
  };

  const runClientValidation = () => {
    const errors: string[] = [];
    const expectedNodeTypes = config.nodeTypes;
    const userNodeTypes: Record<string, number> = {};

    nodes.forEach((node) => {
      const nodeType = node.type || "default";
      userNodeTypes[nodeType] = (userNodeTypes[nodeType] || 0) + 1;
    });

    for (const [type, count] of Object.entries(expectedNodeTypes)) {
      const userCount = userNodeTypes[type] || 0;
      if (userCount < count) {
        errors.push(`Missing ${count - userCount} ${type} node(s)`);
      } else if (userCount > count) {
        errors.push(`Too many ${type} nodes (have ${userCount}, expected ${count})`);
      }
    }

    for (const type of Object.keys(userNodeTypes)) {
      if (!expectedNodeTypes[type]) {
        errors.push(`Unexpected node type: ${type}`);
      }
    }

    const expectedEdgeSet = new Set(config.edges.map((e) => `${e.from}->${e.to}`));
    const userEdgeTypes = new Set<string>();

    edges.forEach((edge) => {
      const sourceType = nodes.find((n) => n.id === edge.source)?.type || "default";
      const targetType = nodes.find((n) => n.id === edge.target)?.type || "default";
      userEdgeTypes.add(`${sourceType}->${targetType}`);
    });

    config.edges.forEach((expectedEdge) => {
      if (!userEdgeTypes.has(`${expectedEdge.from}->${expectedEdge.to}`)) {
        errors.push(`Missing connection: ${expectedEdge.from} -> ${expectedEdge.to}`);
      }
    });

    const valid = errors.length === 0;
    setClientValidation({ valid, errors });
    return valid;
  };

  const handleValidate = async () => {
    const clientValid = runClientValidation();

    setIsValidating(true);

    try {
      const userDiagram: UserDiagram = {
        nodes: nodes.map((node) => ({
          id: node.id,
          type: node.type || "default",
        })),
        edges: edges.map((edge) => ({
          from: edge.source,
          to: edge.target,
        })),
      };

      const response = await axios.post<DiagramValidationResult>(
        `${API_BASE_URL}/api/v1/lessons/${lessonId}/validate-diagram`,
        { diagram: userDiagram }
      );
      setValidationResult(response.data);
    } catch (err) {
      console.error("Validation error:", err);
    } finally {
      setIsValidating(false);
    }
  };

  const handleClear = () => {
    setNodes([]);
    setEdges([]);
    setValidationResult(null);
    setClientValidation(null);
  };

  const showValidation = clientValidation || validationResult;
  const isPassed = validationResult?.valid || clientValidation?.valid;

  return (
    <div className="space-y-4 rounded-lg border border-zinc-200 dark:border-zinc-800 p-6">
      <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-100">
        System Design Exercise
      </h2>

      <div className="flex flex-wrap gap-2">
        {nodeTypeDefinitions.map((def) => (
          <button
            key={def.type}
            onClick={() => addNode(def.type, def.label, def.color)}
            className="rounded-md px-3 py-1.5 text-sm font-medium text-zinc-900 transition-colors hover:opacity-80"
            style={{ background: def.color }}
          >
            + {def.label}
          </button>
        ))}
      </div>

      <div className="h-96 w-full overflow-hidden rounded-lg border border-zinc-300 dark:border-zinc-700">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={customNodeTypes as NodeTypes}
          fitView
        >
          <Background />
          <Controls />
          <MiniMap />
        </ReactFlow>
      </div>

      <div className="flex gap-3">
        <button
          onClick={handleValidate}
          disabled={isValidating || nodes.length === 0}
          className="rounded-md bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-zinc-400"
        >
          {isValidating ? "Validating..." : "Validate Diagram"}
        </button>
        <button
          onClick={handleClear}
          className="rounded-md bg-zinc-200 px-4 py-2 font-medium text-zinc-700 transition-colors hover:bg-zinc-300 dark:bg-zinc-700 dark:text-zinc-200 dark:hover:bg-zinc-600"
        >
          Clear
        </button>
        <button
          onClick={runClientValidation}
          disabled={nodes.length === 0}
          className="rounded-md bg-purple-600 px-4 py-2 font-medium text-white transition-colors hover:bg-purple-700 disabled:cursor-not-allowed disabled:bg-zinc-400"
        >
          Quick Check
        </button>
      </div>

      {showValidation && (
        <div
          className={`rounded-md border p-4 ${
            isPassed
              ? "border-green-500 bg-green-50 dark:bg-green-900/30"
              : "border-red-500 bg-red-50 dark:bg-red-900/30"
          }`}
        >
          <p className="font-medium">
            {isPassed ? "Correct! Great job!" : "Not quite right"}
          </p>
          {showValidation.errors.length > 0 && (
            <ul className="mt-2 list-inside list-disc text-sm text-zinc-600 dark:text-zinc-300">
              {showValidation.errors.map((err, i) => (
                <li key={i}>{err}</li>
              ))}
            </ul>
          )}
          {validationResult && (
            <p className="mt-2 font-semibold">Score: {validationResult.score}%</p>
          )}
        </div>
      )}

      <div className="text-sm text-zinc-500 dark:text-zinc-400">
        <p className="font-medium">Expected:</p>
        <ul className="mt-1 list-inside list-disc">
          {Object.entries(config.nodeTypes).map(([type, count]) => (
            <li key={type}>
              {count} {type}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}