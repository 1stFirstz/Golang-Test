package constants

const (
	WipingClothProductId = "WIPING-CLOTH"
	CleanerSuffix        = "-CLEANNER"
	BundleSeparator      = "/"
	ProductIdSeparator   = "-"
)

const (
	ProductIdPattern     = `(FG[0-9A-Z]+-[A-Z]+-[A-Z0-9\-]+)(\*\d+)?`
	QtyMultiplierPattern = `\*(\d+)`
)

const (
	DefaultQtyMultiplier = 1
	DefaultUnitPrice     = 0.0
	DefaultTotalPrice    = 0.0
	MinProductParts      = 3
)

const (
	MaterialIdPartIndex = 1
	MinRegexMatches     = 2
	QtyRegexMatches     = 2
)

var (
	TextureOrder = []string{"CLEAR", "MATTE", "PRIVACY"}
)
