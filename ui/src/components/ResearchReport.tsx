import { useEffect, useRef, useState, useMemo, useCallback } from "react";
import type { TopLevelSpec } from "vega-lite";

/* ─── Vega-Lite Chart Component ─── */
function VegaChart({ spec }: { spec: TopLevelSpec }) {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      if (!containerRef.current) return;
      const vegaEmbed = (await import("vega-embed")).default;
      if (cancelled) return;
      containerRef.current.innerHTML = "";
      await vegaEmbed(containerRef.current, spec as never, {
        actions: false,
        theme: "dark",
        config: {
          background: "transparent",
          axis: { labelColor: "#8888aa", titleColor: "#e8e8f0", gridColor: "#1e1e3a" },
          legend: { labelColor: "#8888aa", titleColor: "#e8e8f0" },
          title: { color: "#e8e8f0" },
        },
      });
    })();
    return () => { cancelled = true; };
  }, [spec]);

  return <div ref={containerRef} className="vega-chart-wrapper" />;
}

/* ─── Simple MDX-like parser ─── */
/*
 * Since we're in a Vite SPA (not Next.js), we parse the MDX string ourselves
 * instead of pulling in a full MDX compiler. We handle:
 *   - headings (#, ##, ###)
 *   - paragraphs
 *   - bold / italic / inline code
 *   - unordered & ordered lists
 *   - fenced code blocks (``` ... ```)
 *   - <VegaLite> JSON blocks → rendered as interactive charts
 *   - blockquotes (>)
 *   - horizontal rules (---)
 *   - tables (| ... |)
 */

interface ParsedBlock {
  type: "h1" | "h2" | "h3" | "p" | "ul" | "ol" | "code" | "vega" | "hr" | "blockquote" | "table";
  content: string;
  items?: string[];
  rows?: string[][];
  lang?: string;
}

function parseMdx(mdx: string): ParsedBlock[] {
  const lines = mdx.split("\n");
  const blocks: ParsedBlock[] = [];
  let i = 0;

  while (i < lines.length) {
    const line = lines[i];

    // Blank line
    if (line.trim() === "") { i++; continue; }

    // Horizontal rule
    if (/^---+$/.test(line.trim())) {
      blocks.push({ type: "hr", content: "" });
      i++; continue;
    }

    // Fenced code / VegaLite
    if (line.trim().startsWith("```")) {
      const lang = line.trim().slice(3).trim();
      const codeLines: string[] = [];
      i++;
      while (i < lines.length && !lines[i].trim().startsWith("```")) {
        codeLines.push(lines[i]);
        i++;
      }
      i++; // skip closing ```
      const code = codeLines.join("\n");

      if (lang.toLowerCase() === "vegalite" || lang.toLowerCase() === "vega-lite") {
        blocks.push({ type: "vega", content: code });
      } else {
        blocks.push({ type: "code", content: code, lang });
      }
      continue;
    }

    // <VegaLite> component blocks
    if (line.trim().startsWith("<VegaLite")) {
      const specLines: string[] = [];
      i++;
      while (i < lines.length && !lines[i].trim().startsWith("</VegaLite")) {
        specLines.push(lines[i]);
        i++;
      }
      i++; // skip closing tag
      blocks.push({ type: "vega", content: specLines.join("\n") });
      continue;
    }

    // Headings
    if (line.startsWith("#### ")) { blocks.push({ type: "h3", content: line.slice(5) }); i++; continue; }
    if (line.startsWith("### ")) { blocks.push({ type: "h3", content: line.slice(4) }); i++; continue; }
    if (line.startsWith("## "))  { blocks.push({ type: "h2", content: line.slice(3) }); i++; continue; }
    if (line.startsWith("# "))   { blocks.push({ type: "h1", content: line.slice(2) }); i++; continue; }

    // Blockquote
    if (line.startsWith("> ")) {
      const quoteLines: string[] = [];
      while (i < lines.length && lines[i].startsWith("> ")) {
        quoteLines.push(lines[i].slice(2));
        i++;
      }
      blocks.push({ type: "blockquote", content: quoteLines.join("\n") });
      continue;
    }

    // Table
    if (line.includes("|") && line.trim().startsWith("|")) {
      const tableRows: string[][] = [];
      while (i < lines.length && lines[i].includes("|") && lines[i].trim().startsWith("|")) {
        const row = lines[i].trim().split("|").filter(Boolean).map(c => c.trim());
        // Skip separator rows
        if (!row.every(c => /^[-:]+$/.test(c))) {
          tableRows.push(row);
        }
        i++;
      }
      blocks.push({ type: "table", content: "", rows: tableRows });
      continue;
    }

    // Unordered list
    if (/^[-*]\s/.test(line.trim())) {
      const items: string[] = [];
      while (i < lines.length && /^[-*]\s/.test(lines[i].trim())) {
        items.push(lines[i].trim().slice(2));
        i++;
      }
      blocks.push({ type: "ul", content: "", items });
      continue;
    }

    // Ordered list
    if (/^\d+\.\s/.test(line.trim())) {
      const items: string[] = [];
      while (i < lines.length && /^\d+\.\s/.test(lines[i].trim())) {
        items.push(lines[i].trim().replace(/^\d+\.\s/, ""));
        i++;
      }
      blocks.push({ type: "ol", content: "", items });
      continue;
    }

    // Paragraph (collect consecutive non-empty lines)
    const paraLines: string[] = [];
    while (i < lines.length && lines[i].trim() !== "" && !lines[i].trim().startsWith("#") && !lines[i].trim().startsWith("```") && !lines[i].trim().startsWith("> ") && !/^[-*]\s/.test(lines[i].trim()) && !/^\d+\.\s/.test(lines[i].trim()) && !lines[i].trim().startsWith("<VegaLite") && !(lines[i].includes("|") && lines[i].trim().startsWith("|")) && !/^---+$/.test(lines[i].trim())) {
      paraLines.push(lines[i]);
      i++;
    }
    if (paraLines.length) {
      blocks.push({ type: "p", content: paraLines.join(" ") });
    } else {
      // Fallback: If no condition matched and we didn't collect any paragraph lines, 
      // forcefully increment `i` to avoid an infinite loop (e.g. unhandled '#### ' heading).
      // We'll just treat the unhandled line as a paragraph block.
      blocks.push({ type: "p", content: lines[i] });
      i++;
    }
  }

  return blocks;
}

/* Inline markdown formatting */
function renderInline(text: string): string {
  return text
    .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.+?)\*/g, "<em>$1</em>")
    .replace(/`(.+?)`/g, "<code>$1</code>");
}

/* ─── Main Report Component ─── */
interface ResearchReportProps {
  mdxContent: string;
  onBack: () => void;
}

export default function ResearchReport({ mdxContent, onBack }: ResearchReportProps) {
  const blocks = useMemo(() => parseMdx(mdxContent), [mdxContent]);
  const [revealed, setRevealed] = useState(false);

  useEffect(() => {
    // Trigger animation on mount
    requestAnimationFrame(() => setRevealed(true));
  }, []);

  const renderBlock = useCallback((block: ParsedBlock, idx: number) => {
    switch (block.type) {
      case "h1":
        return <h1 key={idx} dangerouslySetInnerHTML={{ __html: renderInline(block.content) }} />;
      case "h2":
        return <h2 key={idx} dangerouslySetInnerHTML={{ __html: renderInline(block.content) }} />;
      case "h3":
        return <h3 key={idx} dangerouslySetInnerHTML={{ __html: renderInline(block.content) }} />;
      case "p":
        return <p key={idx} dangerouslySetInnerHTML={{ __html: renderInline(block.content) }} />;
      case "hr":
        return <hr key={idx} />;
      case "blockquote":
        return <blockquote key={idx} dangerouslySetInnerHTML={{ __html: renderInline(block.content) }} />;
      case "ul":
        return (
          <ul key={idx}>
            {block.items!.map((item, j) => (
              <li key={j} dangerouslySetInnerHTML={{ __html: renderInline(item) }} />
            ))}
          </ul>
        );
      case "ol":
        return (
          <ol key={idx}>
            {block.items!.map((item, j) => (
              <li key={j} dangerouslySetInnerHTML={{ __html: renderInline(item) }} />
            ))}
          </ol>
        );
      case "code":
        return (
          <pre key={idx}>
            <code>{block.content}</code>
          </pre>
        );
      case "table":
        if (!block.rows || block.rows.length === 0) return null;
        return (
          <table key={idx}>
            <thead>
              <tr>
                {block.rows[0].map((cell, j) => (
                  <th key={j} dangerouslySetInnerHTML={{ __html: renderInline(cell) }} />
                ))}
              </tr>
            </thead>
            <tbody>
              {block.rows.slice(1).map((row, ri) => (
                <tr key={ri}>
                  {row.map((cell, j) => (
                    <td key={j} dangerouslySetInnerHTML={{ __html: renderInline(cell) }} />
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        );
      case "vega":
        try {
          const spec = JSON.parse(block.content) as TopLevelSpec;
          return <VegaChart key={idx} spec={spec} />;
        } catch {
          return (
            <pre key={idx} style={{ color: "var(--red)" }}>
              <code>Invalid Vega-Lite spec: {block.content}</code>
            </pre>
          );
        }
      default:
        return null;
    }
  }, []);

  return (
    <div className="report-container" id="research-report">
      <div className="report-header">
        <button className="btn-secondary report-back" onClick={onBack} id="report-back-btn">
          ← New Research
        </button>
      </div>
      <div className="report-content">
        {blocks.map((block, idx) => (
          <div
            key={idx}
            className={`report-block-reveal ${revealed ? "visible" : ""}`}
            style={{ animationDelay: `${idx * 80}ms` }}
          >
            {renderBlock(block, idx)}
          </div>
        ))}
      </div>
    </div>
  );
}
