package zones

type CreationServices struct {
	ZoneFactory   *ZoneFactory
	CastleFactory *CastleFactory
	RoadFactory   *RoadFactory
}

func NewCreationServices(castleFactory *CastleFactory, roadFactory *RoadFactory) *CreationServices {
	if castleFactory == nil {
		castleFactory = NewCastleFactory()
	}
	if roadFactory == nil {
		roadFactory = NewRoadFactory()
	}
	return &CreationServices{
		ZoneFactory:   NewZoneFactory(castleFactory, roadFactory),
		CastleFactory: castleFactory,
		RoadFactory:   roadFactory,
	}
}
