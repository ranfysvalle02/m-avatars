.PHONY: check gallery

# Every instrument must play the same melody.
check:
	python3 identicon.py --check golden.json
	go run identicon.go -check golden.json
	node identicon.mjs --check golden.json

# Each image is rendered by a different language, on purpose.
gallery:
	python3 identicon.py --format svg captain-pixel > gallery/captain-pixel.svg
	go run identicon.go -format svg -size 7 -density 50 mary-had-a-little-lamb > gallery/mary-had-a-little-lamb.svg
	node identicon.mjs --format svg --size 9 --density 48 --shape circle ranfysvalle02 > gallery/ranfysvalle02.svg
	python3 identicon.py --format svg --size 11 --density 45 --shape square mongodb > gallery/mongodb.svg
	go run identicon.go -format svg İlkay > gallery/ilkay.svg
	node identicon.mjs --format svg ΟΔΥΣΣΕΑΣ > gallery/odysseas.svg
	python3 identicon.py --format svg --shape square captain-pixel > gallery/captain-pixel-square.svg
	python3 identicon.py --format svg --shape circle captain-pixel > gallery/captain-pixel-circle.svg
	# What the original v0 ports drew for "İlkay". Go's ToLower gave "ilkay"; Python's lower() gave "i" + U+0307.
	python3 identicon.py --format svg ilkay > gallery/ilkay-before-go.svg
	python3 identicon.py --format svg "$$(printf 'i\314\207lkay')" > gallery/ilkay-before-python.svg
