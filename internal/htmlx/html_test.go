package htmlx

import "testing"

func TestParseToleratesHTMLAndNormalizesText(t *testing.T) {
	root, err := Parse([]byte(`<html lang=ru><body><input disabled><div class=card>A&nbsp; B<img src=x><p>C &amp; D</div></body></html>`))
	if root == nil {
		t.Fatal("nil root")
	}
	// A partial parser error is acceptable for malformed HTML; the tree must still be useful.
	_ = err
	card := FindFirst(root, func(n *Node) bool { return n.HasClassPart("card") })
	if card == nil {
		t.Fatal("card not found")
	}
	if got := Text(card); got != "A B C & D" {
		t.Fatalf("Text() = %q", got)
	}
	if got := ResolveURL("https://somon.tj/a/", "/adv/1"); got != "https://somon.tj/adv/1" {
		t.Fatalf("ResolveURL() = %q", got)
	}
}
