// AlertDialog open/close + trigger/close delegation, vendored from
// workspace/heroui-go's static/js/alert_dialog.js. afterAnimations /
// afterOwnAnimations live in animations.js (load that first).

function syncVisualViewportHeight() {
  if (!window.visualViewport) return;
  document.documentElement.style.setProperty(
    "--visual-viewport-height",
    window.visualViewport.height + "px"
  );
}
if (window.visualViewport) {
  window.visualViewport.addEventListener("resize", syncVisualViewportHeight);
  syncVisualViewportHeight();
}

function lockScroll() {
  var sw = window.innerWidth - document.documentElement.clientWidth;
  document.documentElement.style.overflow = "hidden";
  if (sw > 0) document.body.style.paddingRight = sw + "px";
}

function unlockScroll() {
  document.documentElement.style.overflow = "";
  document.body.style.paddingRight = "";
}

function alertDialogOpen(backdrop) {
  var container = backdrop.querySelector('[data-slot="alert-dialog-container"]');
  backdrop.removeAttribute("hidden");
  lockScroll();
  backdrop.setAttribute("data-entering", "true");
  if (container) container.setAttribute("data-entering", "true");
  afterAnimations(backdrop, function () {
    backdrop.removeAttribute("data-entering");
    if (container) container.removeAttribute("data-entering");
  });
}

function alertDialogClose(backdrop) {
  var container = backdrop.querySelector('[data-slot="alert-dialog-container"]');
  backdrop.setAttribute("data-exiting", "true");
  if (container) container.setAttribute("data-exiting", "true");
  afterOwnAnimations([backdrop, container], function () {
    backdrop.setAttribute("hidden", "");
    backdrop.removeAttribute("data-exiting");
    if (container) container.removeAttribute("data-exiting");
    unlockScroll();
  });
}

document.addEventListener("click", function (e) {
  // Trigger: root'un backdrop dışındaki ilk buton çocuğu
  var trigger = e.target.closest('[data-slot="alert-dialog-root"] > button');
  if (trigger && !trigger.closest('[data-slot="alert-dialog-backdrop"]')) {
    var backdrop = trigger
      .closest('[data-slot="alert-dialog-root"]')
      .querySelector('[data-slot="alert-dialog-backdrop"]');
    if (backdrop) alertDialogOpen(backdrop);
    return;
  }
  // slot="close" olan her element dialogu kapatır
  var closer = e.target.closest('[slot="close"]');
  if (closer) {
    var bd = closer.closest('[data-slot="alert-dialog-backdrop"]');
    if (bd) alertDialogClose(bd);
    return;
  }
  // Backdrop'a tıklama (yalnızca isDismissable)
  if (
    e.target.getAttribute &&
    e.target.getAttribute("data-slot") === "alert-dialog-backdrop" &&
    e.target.getAttribute("data-dismissable") === "true"
  ) {
    alertDialogClose(e.target);
  }
});

// htmx confirm butonları (Sil/Sıfırla), tıklanınca isteği hemen yollar; kapanış
// animasyonu bitmeden hx-swap kartı/li'yi (backdrop dahil) DOM'dan kaldırabilir,
// bu durumda alertDialogClose'un animationend/transitionend dinleyicisi hiç
// tetiklenmez ve unlockScroll() çağrılmadığı için sayfa "overflow: hidden"da
// kilitli kalır. İstek başlar başlamaz kilidi burada açıp backdrop'u
// gizleyerek bu senaryoyu garantiye alıyoruz.
document.body.addEventListener("htmx:beforeRequest", function (e) {
  var backdrop = e.target.closest('[data-slot="alert-dialog-backdrop"]');
  if (backdrop) {
    backdrop.setAttribute("hidden", "");
    backdrop.removeAttribute("data-exiting");
    var container = backdrop.querySelector('[data-slot="alert-dialog-container"]');
    if (container) container.removeAttribute("data-exiting");
    unlockScroll();
  }
});

document.addEventListener("keydown", function (e) {
  if (e.key !== "Escape") return;
  var backdrop = document.querySelector('[data-slot="alert-dialog-backdrop"]:not([hidden])');
  if (backdrop && backdrop.getAttribute("data-keyboard-dismiss-disabled") !== "true") {
    alertDialogClose(backdrop);
  }
});
