module JD_central

replace (
	github.com/Gerardo115pp/PatriotLib/PatriotBlockchain => ../../patriotslibs.com/PatriotBlockchain
	github.com/Gerardo115pp/PatriotLib/PatriotRouter => ../../patriotslibs.com/Router
	github.com/Gerardo115pp/PatriotLib/PatriotUtils => ../../patriotslibs.com/patriot-utils
)

go 1.16

require (
	github.com/Gerardo115pp/PatriotLib/PatriotBlockchain v0.0.0-00010101000000-000000000000
	github.com/Gerardo115pp/PatriotLib/PatriotRouter v0.0.0-00010101000000-000000000000
	github.com/Gerardo115pp/PatriotLib/PatriotUtils v0.0.0-00010101000000-000000000000
)
