// Pinch-to-zoom + pan + double-tap for the "Okunuş" image inside an alert
// dialog. The page itself disables native pinch-zoom (viewport meta:
// user-scalable=no) because position:fixed dialogs don't track the visual
// viewport under browser-level pinch-zoom on mobile Safari, which let the
// page underneath show through. This reimplements zoom scoped to just the
// <img>, driven by CSS transforms, so the backdrop is never affected.

(function () {
  var MAX_SCALE = 4;
  var DOUBLE_TAP_SCALE = 2.5;
  var DOUBLE_TAP_MS = 300;
  var DOUBLE_TAP_SLOP = 24;

  function clamp(v, min, max) {
    return Math.max(min, Math.min(max, v));
  }

  function setupZoomable(container) {
    if (container.__zoomSetup) return;
    container.__zoomSetup = true;

    var img = container.querySelector("img");
    if (!img) return;

    var state = { scale: 1, x: 0, y: 0 };
    var pointers = new Map();
    var pinchStartDist = 0;
    var pinchStartScale = 1;
    var panStart = null;
    var lastTapTime = 0;
    var lastTapPos = null;

    function apply() {
      img.style.transform =
        "translate(" + state.x + "px," + state.y + "px) scale(" + state.scale + ")";
    }

    container.__resetZoom = function () {
      state.scale = 1;
      state.x = 0;
      state.y = 0;
      apply();
    };

    function clampPan() {
      var rect = container.getBoundingClientRect();
      var maxX = (rect.width * (state.scale - 1)) / 2;
      var maxY = (rect.height * (state.scale - 1)) / 2;
      state.x = clamp(state.x, -maxX, maxX);
      state.y = clamp(state.y, -maxY, maxY);
    }

    function midpoint(a, b) {
      return { x: (a.x + b.x) / 2, y: (a.y + b.y) / 2 };
    }

    function dist(a, b) {
      return Math.hypot(a.x - b.x, a.y - b.y);
    }

    function setScaleAtPoint(newScale, point) {
      var rect = container.getBoundingClientRect();
      var cx = point.x - rect.left - rect.width / 2;
      var cy = point.y - rect.top - rect.height / 2;
      var ratio = newScale / state.scale;
      state.x = cx - (cx - state.x) * ratio;
      state.y = cy - (cy - state.y) * ratio;
      state.scale = newScale;
      clampPan();
      apply();
    }

    container.addEventListener("pointerdown", function (e) {
      pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });
      container.setPointerCapture(e.pointerId);

      if (pointers.size === 2) {
        var pts = Array.from(pointers.values());
        pinchStartDist = dist(pts[0], pts[1]);
        pinchStartScale = state.scale;
        panStart = null;
      } else if (pointers.size === 1) {
        if (state.scale > 1) {
          panStart = { x: e.clientX, y: e.clientY, ox: state.x, oy: state.y };
        }

        var now = Date.now();
        if (
          lastTapPos &&
          now - lastTapTime < DOUBLE_TAP_MS &&
          Math.hypot(e.clientX - lastTapPos.x, e.clientY - lastTapPos.y) < DOUBLE_TAP_SLOP
        ) {
          if (state.scale > 1) {
            container.__resetZoom();
          } else {
            setScaleAtPoint(DOUBLE_TAP_SCALE, { x: e.clientX, y: e.clientY });
          }
          lastTapTime = 0;
          lastTapPos = null;
        } else {
          lastTapTime = now;
          lastTapPos = { x: e.clientX, y: e.clientY };
        }
      }
    });

    container.addEventListener("pointermove", function (e) {
      if (!pointers.has(e.pointerId)) return;
      pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });

      if (pointers.size === 2) {
        var pts = Array.from(pointers.values());
        var d = dist(pts[0], pts[1]);
        if (pinchStartDist > 0) {
          var newScale = clamp((d / pinchStartDist) * pinchStartScale, 1, MAX_SCALE);
          setScaleAtPoint(newScale, midpoint(pts[0], pts[1]));
        }
      } else if (pointers.size === 1 && panStart) {
        state.x = panStart.ox + (e.clientX - panStart.x);
        state.y = panStart.oy + (e.clientY - panStart.y);
        clampPan();
        apply();
      }
    });

    function endPointer(e) {
      pointers.delete(e.pointerId);
      if (pointers.size < 2) pinchStartDist = 0;
      if (pointers.size === 0) panStart = null;
      if (state.scale <= 1) container.__resetZoom();
    }

    container.addEventListener("pointerup", endPointer);
    container.addEventListener("pointercancel", endPointer);
  }

  document.body.addEventListener("htmx:afterSettle", scan);
  document.addEventListener("DOMContentLoaded", scan);
  function scan() {
    document.querySelectorAll('[data-slot="zoomable-image"]').forEach(setupZoomable);
  }
  scan();

  // Reset zoom whenever a dialog closes so reopening starts fresh.
  document.addEventListener("click", function (e) {
    var closer = e.target.closest('[slot="close"]');
    var backdrop = closer && closer.closest('[data-slot="alert-dialog-backdrop"]');
    if (!backdrop) return;
    backdrop.querySelectorAll('[data-slot="zoomable-image"]').forEach(function (el) {
      if (el.__resetZoom) el.__resetZoom();
    });
  });
  document.body.addEventListener("htmx:beforeRequest", function (e) {
    var backdrop = e.target.closest('[data-slot="alert-dialog-backdrop"]');
    if (!backdrop) return;
    backdrop.querySelectorAll('[data-slot="zoomable-image"]').forEach(function (el) {
      if (el.__resetZoom) el.__resetZoom();
    });
  });
})();
