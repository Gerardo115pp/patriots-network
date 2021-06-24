module gw-node

replace (
	github.com/Gerardo115pp/PatriotLib/PatriotBlockchain => ../../patriotslibs.com/PatriotBlockchain
	github.com/Gerardo115pp/PatriotLib/PatriotRouter => ../../patriotslibs.com/Router
)

go 1.16

require (
	github.com/Gerardo115pp/PatriotLib/PatriotBlockchain v0.0.0-00010101000000-000000000000
	github.com/Gerardo115pp/PatriotLib/PatriotRouter v0.0.0-00010101000000-000000000000
)
