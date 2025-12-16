package embedfiles

import "embed"

//go:embed assets/carrier.db
var CarrierDB embed.FS
