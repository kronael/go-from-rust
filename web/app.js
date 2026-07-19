// Vanilla-JS hash router for the interactive lesson tour. No dependencies.
(function () {
  "use strict";

  function wordSet(text) {
    return new Set(text.split(" "));
  }

  const STORAGE_PREFIX = "go-from-rust:edit:";
  const GO_KEYWORDS = wordSet(
    "break case chan const continue default defer else fallthrough for func go goto " +
      "if import interface map package range return select struct switch type var",
  );
  const GO_TYPES = wordSet(
    "any bool byte comparable complex64 complex128 error float32 float64 int int8 " +
      "int16 int32 int64 rune string uint uint8 uint16 uint32 uint64 uintptr",
  );
  const GO_BUILTINS = wordSet(
    "append cap clear close complex copy delete imag len make max min new panic " +
      "print println real recover",
  );
  const GO_CONSTANTS = wordSet("false iota nil true");
  const NUMBER_PATTERN = new RegExp(
    "^(?:" +
      "0[xX][0-9a-fA-F_]+(?:\\.[0-9a-fA-F_]*)?(?:[pP][+-]?[0-9_]+)?" +
      "|0[bB][01_]+|0[oO][0-7_]+" +
      "|(?:[0-9][0-9_]*)(?:\\.[0-9_]*)?(?:[eE][+-]?[0-9_]+)?" +
      "|\\.[0-9_]+(?:[eE][+-]?[0-9_]+)?" +
      ")(?:i)?",
  );

  const els = {
    title: document.getElementById("lesson-title"),
    value: document.getElementById("lesson-value"),
    summary: document.getElementById("lesson-summary"),
    editor: document.getElementById("editor"),
    highlight: document.getElementById("editor-highlight"),
    highlightCode: document.getElementById("editor-highlight-code"),
    output: document.getElementById("output-code"),
    status: document.getElementById("status-msg"),
    select: document.getElementById("lesson-select"),
    prevBtn: document.getElementById("prev-btn"),
    nextBtn: document.getElementById("next-btn"),
    importsBtn: document.getElementById("imports-btn"),
    runBtn: document.getElementById("run-btn"),
    formatBtn: document.getElementById("format-btn"),
    resetBtn: document.getElementById("reset-btn"),
    lineNumbers: document.getElementById("line-numbers"),
    connStatus: document.getElementById("conn-status"),
  };

  let lessons = [];
  let current = null;
  let activeRequest = null;
  let importsVisible = false;
  let hiddenImports = "";

  function lessonKey(id) {
    return String(id).padStart(2, "0");
  }

  function storageKey(lesson) {
    return `${STORAGE_PREFIX}${lesson.Filename}:${lesson.Hash}`;
  }

  function setStatus(text) {
    els.status.textContent = text;
  }

  function escapeHTML(text) {
    return text.replace(/[&<>]/g, (character) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;" })[character]);
  }

  function highlightedToken(className, text) {
    return `<span class="${className}">${escapeHTML(text)}</span>`;
  }

  function isIdentifierStart(character) {
    return character === "_" || /[A-Za-z]/.test(character);
  }

  function isIdentifierPart(character) {
    return character === "_" || /[A-Za-z0-9]/.test(character);
  }

  function quotedEnd(source, start, quote) {
    let index = start + 1;
    while (index < source.length) {
      if (source[index] === quote) return index + 1;
      if (quote !== "`" && source[index] === "\\") {
        index += 2;
      } else {
        index += 1;
      }
    }
    return source.length;
  }

  function highlightGo(source) {
    let html = "";
    let index = 0;
    while (index < source.length) {
      if (source.startsWith("//", index)) {
        const newline = source.indexOf("\n", index);
        const end = newline === -1 ? source.length : newline;
        html += highlightedToken("tok-comment", source.slice(index, end));
        index = end;
        continue;
      }
      if (source.startsWith("/*", index)) {
        const close = source.indexOf("*/", index + 2);
        const end = close === -1 ? source.length : close + 2;
        html += highlightedToken("tok-comment", source.slice(index, end));
        index = end;
        continue;
      }

      const character = source[index];
      if (character === '"' || character === "'" || character === "`") {
        const end = quotedEnd(source, index, character);
        html += highlightedToken("tok-string", source.slice(index, end));
        index = end;
        continue;
      }

      const number = source.slice(index).match(NUMBER_PATTERN);
      if (number) {
        html += highlightedToken("tok-number", number[0]);
        index += number[0].length;
        continue;
      }

      if (isIdentifierStart(character)) {
        let end = index + 1;
        while (end < source.length && isIdentifierPart(source[end])) end += 1;
        const word = source.slice(index, end);
        if (GO_KEYWORDS.has(word)) {
          html += highlightedToken("tok-keyword", word);
        } else if (GO_TYPES.has(word)) {
          html += highlightedToken("tok-type", word);
        } else if (GO_BUILTINS.has(word)) {
          html += highlightedToken("tok-builtin", word);
        } else if (GO_CONSTANTS.has(word)) {
          html += highlightedToken("tok-constant", word);
        } else {
          html += escapeHTML(word);
        }
        index = end;
        continue;
      }

      html += escapeHTML(character);
      index += 1;
    }
    return html;
  }

  function syncHighlightScroll() {
    els.highlight.scrollTop = els.editor.scrollTop;
    els.highlight.scrollLeft = els.editor.scrollLeft;
    els.lineNumbers.scrollTop = els.editor.scrollTop;
  }

  function renderHighlight() {
    const source = els.editor.value;
    els.highlightCode.innerHTML = highlightGo(source) + (source.endsWith("\n") ? " " : "");
    els.lineNumbers.textContent = source
      .split("\n")
      .map((_, index) => index + 1)
      .join("\n");
    syncHighlightScroll();
  }

  function splitImports(source) {
    const match = /^(package\s+[A-Za-z_]\w*\n)\n?(import\s+(?:\([\s\S]*?\n\)|[^\n]+)\n)\n?/.exec(source);
    if (!match) return { visible: source, imports: "" };
    return {
      visible: `${match[1]}\n${source.slice(match[0].length)}`,
      imports: match[2].trimEnd(),
    };
  }

  function joinImports(source, imports) {
    if (!imports) return source;
    return source.replace(/^(package\s+[A-Za-z_]\w*\n)\n?/, `$1\n${imports}\n\n`);
  }

  function fullEditorSource() {
    return importsVisible ? els.editor.value : joinImports(els.editor.value, hiddenImports);
  }

  function updateImportsButton() {
    els.importsBtn.disabled = hiddenImports === "";
    els.importsBtn.textContent = importsVisible ? "Imports on" : "Imports off";
    els.importsBtn.setAttribute("aria-pressed", String(importsVisible));
  }

  function setEditorSource(source) {
    const split = splitImports(source);
    hiddenImports = split.imports;
    els.editor.value = importsVisible ? source : split.visible;
    updateImportsButton();
    renderHighlight();
  }

  function toggleImports() {
    const source = fullEditorSource();
    importsVisible = !importsVisible;
    setEditorSource(source);
  }

  function findLesson(id) {
    return lessons.find((l) => l.ID === id) || null;
  }

  function populateSelect() {
    els.select.innerHTML = "";
    for (const lesson of lessons) {
      const opt = document.createElement("option");
      opt.value = lessonKey(lesson.ID);
      opt.textContent = `${lessonKey(lesson.ID)} — ${lesson.Filename}`;
      els.select.appendChild(opt);
    }
  }

  function loadEditedSource(lesson) {
    const saved = localStorage.getItem(storageKey(lesson));
    return saved !== null ? saved : lesson.Source;
  }

  function renderLesson(lesson) {
    cancelRequest();
    current = lesson;
    importsVisible = false;
    els.title.textContent = `${lessonKey(lesson.ID)} — ${lesson.Filename}`;
    els.value.textContent = `Value ${lesson.Value}/5`;
    els.summary.textContent = lesson.Summary;
    setEditorSource(loadEditedSource(lesson));
    els.output.textContent = "Run the lesson to see output here.";
    els.output.parentElement.dataset.state = "idle";
    els.select.value = lessonKey(lesson.ID);
    els.prevBtn.disabled = lesson.ID <= lessons[0].ID;
    els.nextBtn.disabled = lesson.ID >= lessons[lessons.length - 1].ID;
    setStatus("");
    document.title = `${lessonKey(lesson.ID)} ${lesson.Filename} — Go for Rust Engineers`;
  }

  function navigateTo(id, replace) {
    const hash = `#${lessonKey(id)}`;
    if (replace) {
      history.replaceState(null, "", hash);
    } else {
      location.hash = hash;
    }
  }

  function idFromHash() {
    const match = /^#(\d+)$/.exec(location.hash);
    return match ? parseInt(match[1], 10) : null;
  }

  function onHashChange() {
    const id = idFromHash();
    const lesson = id !== null ? findLesson(id) : null;
    if (lesson) {
      renderLesson(lesson);
    } else {
      navigateTo(lessons[0].ID, true);
      renderLesson(lessons[0]);
    }
  }

  function persistEdit() {
    if (!current) return;
    const source = fullEditorSource();
    if (source === current.Source) {
      localStorage.removeItem(storageKey(current));
    } else {
      localStorage.setItem(storageKey(current), source);
    }
  }

  function extractOutput(payload) {
    if (payload.Errors) {
      return payload.Errors;
    }
    const events = payload.Events || [];
    if (events.length === 0) {
      return "(no output)";
    }
    return events.map((e) => e.Message).join("");
  }

  function setBusy(busy) {
    els.runBtn.disabled = busy;
    els.formatBtn.disabled = busy;
  }

  function beginRequest() {
    cancelRequest();
    activeRequest = new AbortController();
    setBusy(true);
    return activeRequest;
  }

  function cancelRequest() {
    if (activeRequest) {
      activeRequest.abort();
      activeRequest = null;
    }
    setBusy(false);
  }

  function finishRequest(controller) {
    if (activeRequest === controller) {
      activeRequest = null;
      setBusy(false);
    }
  }

  async function postJSON(url, body, signal) {
    const res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
      signal,
    });
    const data = await res.json().catch(() => null);
    return { ok: res.ok, data };
  }

  async function runCode() {
    if (!current) return;
    const lesson = current;
    const controller = beginRequest();
    setStatus("Running…");
    try {
      const { ok, data } = await postJSON("api/run", { body: fullEditorSource() }, controller.signal);
      if (current !== lesson) return;
      if (!ok || !data) {
        els.output.textContent = data && data.error ? data.error.message : "Run failed.";
        els.output.parentElement.dataset.state = "error";
        setStatus("Run failed");
        return;
      }
      els.output.textContent = extractOutput(data);
      if (data.Errors) {
        els.output.parentElement.dataset.state = "error";
        setStatus("Compile failed");
      } else {
        els.output.parentElement.dataset.state = "success";
        setStatus("Run complete");
      }
    } catch (err) {
      if (err.name === "AbortError") return;
      els.output.textContent = "Network error while running.";
      els.output.parentElement.dataset.state = "error";
      setStatus("Run failed");
    } finally {
      finishRequest(controller);
    }
  }

  async function formatCode() {
    if (!current) return;
    const lesson = current;
    const controller = beginRequest();
    setStatus("Formatting…");
    try {
      const { ok, data } = await postJSON("api/format", { body: fullEditorSource() }, controller.signal);
      if (current !== lesson) return;
      if (!ok || !data) {
        setStatus(data && data.error ? data.error.message : "Format failed");
        return;
      }
      setEditorSource(data.body);
      persistEdit();
      setStatus("Formatted");
    } catch (err) {
      if (err.name === "AbortError") return;
      setStatus("Format failed: network error");
    } finally {
      finishRequest(controller);
    }
  }

  function resetCode() {
    if (!current) return;
    setEditorSource(current.Source);
    localStorage.removeItem(storageKey(current));
    setStatus("Reset to original source");
  }

  function stepLesson(delta) {
    if (!current) return;
    const idx = lessons.findIndex((l) => l.ID === current.ID);
    const next = lessons[idx + delta];
    if (next) navigateTo(next.ID, false);
  }

  function wireEvents() {
    window.addEventListener("hashchange", onHashChange);
    els.editor.addEventListener("input", () => {
      if (!importsVisible) {
        const split = splitImports(els.editor.value);
        if (split.imports) {
          hiddenImports = split.imports;
          els.editor.value = split.visible;
        }
        updateImportsButton();
      }
      renderHighlight();
      persistEdit();
    });
    els.editor.addEventListener("scroll", syncHighlightScroll);
    els.importsBtn.addEventListener("click", toggleImports);
    els.runBtn.addEventListener("click", runCode);
    els.formatBtn.addEventListener("click", formatCode);
    els.resetBtn.addEventListener("click", resetCode);
    els.prevBtn.addEventListener("click", () => stepLesson(-1));
    els.nextBtn.addEventListener("click", () => stepLesson(1));
    els.select.addEventListener("change", () => {
      navigateTo(parseInt(els.select.value, 10), false);
    });
    els.editor.addEventListener("keydown", (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
        e.preventDefault();
        runCode();
      }
    });
  }

  async function init() {
    try {
      const res = await fetch("api/lessons");
      if (!res.ok) throw new Error(`lessons request failed: ${res.status}`);
      const data = await res.json();
      lessons = (data.lessons || []).slice().sort((a, b) => a.ID - b.ID);
    } catch (err) {
      els.title.textContent = "Failed to load lessons";
      setStatus("Could not reach the server. Reload to retry.");
      return;
    }
    if (lessons.length === 0) {
      els.title.textContent = "No lessons available";
      return;
    }
    populateSelect();
    wireEvents();
    els.connStatus.textContent = `${lessons.length} lessons · edits stay in this browser`;
    onHashChange();
  }

  init();
})();
