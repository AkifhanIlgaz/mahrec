// Shared animation-wait helpers used by select.js.
// Extracted from workspace/heroui-go's static/js/alert_dialog.js so we don't
// have to pull in the rest of its (unused) modal/dialog logic.

function afterAnimations(el, cb) {
  var done = false;
  var finish = function () {
    if (done) return;
    done = true;
    cb();
  };
  requestAnimationFrame(function () {
    var anims = el.getAnimations ? el.getAnimations({ subtree: true }) : [];
    if (anims.length === 0) return finish();
    Promise.allSettled(anims.map(function (a) { return a.finished; })).then(finish);
  });
  setTimeout(finish, 600);
}

function afterOwnAnimations(elements, cb) {
  var done = false;
  var finish = function () {
    if (done) return;
    done = true;
    cb();
  };
  requestAnimationFrame(function () {
    var anims = [];
    elements.forEach(function (el) {
      if (el && el.getAnimations) anims = anims.concat(el.getAnimations());
    });
    if (anims.length === 0) return finish();
    Promise.allSettled(anims.map(function (a) { return a.finished; })).then(finish);
  });
  setTimeout(finish, 600);
}

// Bumps a goal card's count with a small scale animation whenever the
// /increment endpoint succeeds. The progress bar's own fill transition
// (see goal_item.templ's Fill() calls / .progress-bar__fill in
// @heroui/styles) already animates its width smoothly across the
// morph:outerHTML swap, since idiomorph patches the existing DOM node in
// place rather than replacing it.
(function () {
  document.body.addEventListener("htmx:afterSwap", function (evt) {
    var cfg = evt.detail && evt.detail.requestConfig;
    if (!cfg || !/\/goals\/\d+\/increment$/.test(cfg.path)) return;
    if (window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;

    var card = evt.detail.target;
    var countEl = card && card.querySelector('[data-slot="goal-count"]');
    if (!countEl || !countEl.animate) return;

    countEl.animate(
      [
        { transform: "scale(1)" },
        { transform: "scale(1.18)" },
        { transform: "scale(1)" },
      ],
      { duration: 380, easing: "cubic-bezier(0.34, 1.56, 0.64, 1)" }
    );
  });
})();
