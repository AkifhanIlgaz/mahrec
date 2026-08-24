// Bridges a goal card's kebab dropdown items (Sil/Sıfırla) to the hidden
// AlertDialog trigger buttons rendered as their siblings, outside the
// dropdown's popover. They can't live inside the popover: dropdown.js hides
// it via the `hidden` attribute, and a `hidden` ancestor forces
// display:none on every descendant — including the confirmation dialog's
// fixed-position backdrop — so the dialog would flash and vanish with it.
document.addEventListener("click", function (e) {
	var item = e.target.closest("[data-open-dialog]");
	if (!item) return;
	var trigger = document.getElementById(item.getAttribute("data-open-dialog"));
	if (trigger) trigger.click();
});
