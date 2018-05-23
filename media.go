package ticketswitch

// Media describes some event media asset
type Media struct {
	// caption in plain text describing the asset.
	Caption string
	// caption as html describing the asset.
	Caption_html string
	// name of the asset.
	Name string
	// url for the asset.
	URL string
	// indicates if the assert url is secure or not.
	Secure bool
	// width of the asset in pixels. Only present on the video
	Width int
	// height of the asset in pixels. Only present on the video asset.
	Height int
}
