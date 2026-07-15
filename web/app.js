// Vanilla-JS hash router for the interactive lesson tour. No dependencies.
(function () {
  "use strict";

  const STORAGE_PREFIX = "go-from-rust:edit:";

  const els = {
    title: document.getElementById("lesson-title"),
    value: document.getElementById("lesson-value"),
    question: document.getElementById("lesson-question"),
    editor: document.getElementById("editor"),
    output: document.getElementById("output-code"),
    status: document.getElementById("status-msg"),
    select: document.getElementById("lesson-select"),
    prevBtn: document.getElementById("prev-btn"),
    nextBtn: document.getElementById("next-btn"),
    runBtn: document.getElementById("run-btn"),
    formatBtn: document.getElementById("format-btn"),
    resetBtn: document.getElementById("reset-btn"),
    connStatus: document.getElementById("conn-status"),
  };

  let lessons = [];
  let current = null;
  let activeRequest = null;

  function lessonKey(id) {
    return String(id).padStart(2, "0");
  }

  function storageKey(lesson) {
    return `${STORAGE_PREFIX}${lesson.Filename}:${lesson.Hash}`;
  }

  function setStatus(text) {
    els.status.textContent = text;
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
    els.title.textContent = `${lessonKey(lesson.ID)} — ${lesson.Filename}`;
    els.value.textContent = lesson.ExpectedFailure
      ? `Value ${lesson.Value}/5 · Expected compiler failure`
      : `Value ${lesson.Value}/5`;
    els.question.textContent = lesson.Question;
    els.editor.value = loadEditedSource(lesson);
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
    if (els.editor.value === current.Source) {
      localStorage.removeItem(storageKey(current));
    } else {
      localStorage.setItem(storageKey(current), els.editor.value);
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
    return { ok: res.ok, status: res.status, data };
  }

  async function runCode() {
    if (!current) return;
    const lesson = current;
    const controller = beginRequest();
    setStatus("Running…");
    try {
      const { ok, data } = await postJSON("/api/run", { body: els.editor.value }, controller.signal);
      if (current !== lesson) return;
      if (!ok || !data) {
        els.output.textContent = data && data.error ? data.error.message : "Run failed.";
        els.output.parentElement.dataset.state = "error";
        setStatus("Run failed");
        return;
      }
      els.output.textContent = extractOutput(data);
      if (data.Errors) {
        els.output.parentElement.dataset.state = lesson.ExpectedFailure ? "expected" : "error";
        setStatus(lesson.ExpectedFailure ? "Expected compiler error" : "Compile failed");
      } else if (lesson.ExpectedFailure) {
        els.output.parentElement.dataset.state = "error";
        setStatus("Expected compiler error did not occur");
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
      const { ok, data } = await postJSON("/api/format", { body: els.editor.value }, controller.signal);
      if (current !== lesson) return;
      if (!ok || !data) {
        setStatus(data && data.error ? data.error.message : "Format failed");
        return;
      }
      els.editor.value = data.body;
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
    els.editor.value = current.Source;
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
    els.editor.addEventListener("input", persistEdit);
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
      if (e.key === "Tab") {
        e.preventDefault();
        const start = els.editor.selectionStart;
        const end = els.editor.selectionEnd;
        els.editor.value = els.editor.value.slice(0, start) + "\t" + els.editor.value.slice(end);
        els.editor.selectionStart = els.editor.selectionEnd = start + 1;
        persistEdit();
      }
    });
  }

  async function init() {
    try {
      const res = await fetch("/api/lessons");
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
