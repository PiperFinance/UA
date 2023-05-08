package conf

const connectionCheckKey = "IsCacheManagerOK"

var (
// CacheManager *cache.Cache[string]
)

func init() {
	//ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
	//	NumCounters: 1000,
	//	MaxCost:     100,
	//	BufferItems: 64,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	//
	//CacheManager = cache.New[string](ristrettoStore)
	//ctx := context.Background()
	//
	//err = CacheManager.Set(ctx, connectionCheckKey, "YES", store.WithCost(2))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//value, err := CacheManager.Get(ctx, connectionCheckKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if value != "YES" {
	//	log.Fatalf("Bad Response Value @%s, suspected 'YES'", value)
	//}
	//err = CacheManager.Delete(ctx, connectionCheckKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
