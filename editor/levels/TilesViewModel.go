package levels

type tilesViewModel struct {
	selectedTiles     tileCoordinates
	textureDisplay    TextureDisplay
	shadowDisplay     ColorDisplay
	cyberColorDisplay ColorDisplay

	restoreFocus bool
	windowOpen   bool
}

func freshTilesViewModel() tilesViewModel {
	return tilesViewModel{
		textureDisplay:    TextureDisplayFloor,
		shadowDisplay:     ColorDisplayNone,
		cyberColorDisplay: ColorDisplayNone,
	}
}
